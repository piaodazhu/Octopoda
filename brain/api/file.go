package api

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver/v3"
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/brain/model"
	"github.com/piaodazhu/Octopoda/brain/workgroup"
	"github.com/piaodazhu/Octopoda/protocols"
)

func pathFixing(path string, base string) string {
	// dstPath: the unpacked files will be moved under this path
	var result strings.Builder
	// find ~
	homePos := -1
	for i, c := range path {
		if c == '~' {
			homePos = i
		}
	}
	if homePos != -1 {
		result.WriteString(path[homePos:])
	} else {
		if len(path) == 0 || path[0] != '/' {
			result.WriteString(base)
		}
		result.WriteString(path)
	}
	return result.String()
}

func FileUpload(ctx *gin.Context) {
	tarfile, _ := ctx.FormFile("tarfile")
	targetPath := ctx.PostForm("targetPath")
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	// tmpPath: the zip file will be unpacked under this path
	var tmpSb strings.Builder
	tmpSb.WriteString(config.GlobalConfig.Workspace.Store)
	tmpSb.WriteString(".octopoda_tmp/")

	tmpPath := tmpSb.String()

	// dstPath: the unpacked files will be moved under this path
	// fmt.Println("targetPath:", targetPath)
	dstPath := pathFixing(targetPath, config.GlobalConfig.Workspace.Store)

	os.Mkdir(tmpPath, os.ModePerm)

	tmpSb.WriteString(tarfile.Filename)

	// tmpExtDir: the dir name after being extracted
	var tmpExtDir strings.Builder
	for i := tmpSb.Len() - 1; i > 0; i-- {
		if tmpSb.String()[i] == '.' {
			tmpExtDir.WriteString(tmpSb.String()[:i])
			break
		} else if tmpSb.String()[i] == '/' {
			logger.Exceptions.Println("tmpExtDir")
			rmsg.Rmsg = "bad tmpExtDir:" + tmpSb.String()
			ctx.JSON(http.StatusBadRequest, rmsg)
			return
		}
	}
	tmpExtDir.WriteString("/*")

	tmpDst, err := os.Create(tmpSb.String())
	if err != nil {
		logger.Exceptions.Println("FileCreate")
		rmsg.Rmsg = "FileCreate:" + err.Error()
		ctx.JSON(http.StatusBadRequest, rmsg)

		os.RemoveAll(tmpPath)
		return
	}

	src, _ := tarfile.Open()

	io.Copy(tmpDst, src)
	tmpDst.Close()
	src.Close()

	taskid := model.BrainTaskManager.CreateTask(1)
	ctx.String(http.StatusAccepted, taskid)

	go func() {
		result := protocols.ExecutionResult{
			Name:   "brain",
			Code:   protocols.ExecOK,
			Result: "OK",
		}

		archiver.DefaultZip.OverwriteExisting = true
		err = archiver.DefaultZip.Unarchive(tmpSb.String(), tmpPath)
		defer os.RemoveAll(tmpPath)

		if err != nil {
			emsg := fmt.Sprintf("unarchive %s to %s err:%v", tmpSb.String(), tmpPath, err)
			result.Code = protocols.ExecProcessError
			result.CommunicationErrorMsg = emsg
			model.BrainTaskManager.AddFailedSubTask(taskid, model.TaskIdGen(), &result)
			return
		}

		cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("mkdir -p %s && cp -r %s %s", dstPath, tmpExtDir.String(), dstPath))
		err = cmd.Run()
		if err != nil {
			emsg := fmt.Sprintf("copy %s to %s err:%v", tmpExtDir.String(), dstPath, err)
			result.Code = protocols.ExecProcessError
			result.CommunicationErrorMsg = emsg
			model.BrainTaskManager.AddFailedSubTask(taskid, model.TaskIdGen(), &result)
			return
		}

		subtaskid := model.TaskIdGen()
		model.BrainTaskManager.AddSubTask(taskid, subtaskid, &result)
		model.BrainTaskManager.DoneSubTask(taskid, subtaskid, &result)
	}()
}

// need fix
func FileSpread(ctx *gin.Context) {
	var fsParams protocols.FileSpreadParams
	err := ctx.ShouldBind(&fsParams)
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	if err != nil {
		logger.Exceptions.Println("FileCreate")
		rmsg.Rmsg = "FileCreate:" + err.Error()
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	// check file
	if fsParams.FileOrDir == "/" {
		fsParams.FileOrDir = "."
	} else if fsParams.FileOrDir[len(fsParams.FileOrDir)-1] == '/' {
		fsParams.FileOrDir = fsParams.FileOrDir[:len(fsParams.FileOrDir)-1]
	}

	if fsParams.TargetPath == "/" || fsParams.TargetPath == "./" {
		fsParams.TargetPath = ""
	}

	fname := pathFixing(fsParams.FileOrDir, config.GlobalConfig.Workspace.Store)
	wrapName := fmt.Sprintf("%d.wrap", time.Now().Nanosecond())
	os.Mkdir(wrapName, os.ModePerm)
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -r %s %s", fname, wrapName))
	err = cmd.Run()
	if err != nil {
		emsg := fmt.Sprintf("error when file %s to %s: %v", fname, wrapName, err)
		logger.Exceptions.Print(emsg)
		rmsg.Rmsg = emsg
		ctx.JSON(http.StatusInternalServerError, rmsg)
		return
	}
	defer os.RemoveAll(wrapName)

	packName := fmt.Sprintf("%s.zip", wrapName)
	archiver.DefaultZip.OverwriteExisting = true
	err = archiver.DefaultZip.Archive([]string{wrapName}, packName)
	if err != nil {
		emsg := fmt.Sprintf("error archive %s to %s: %v", wrapName, packName, err)
		logger.Exceptions.Print(emsg)
		rmsg.Rmsg = emsg
		ctx.JSON(http.StatusInternalServerError, rmsg)
		return
	}

	raw, _ := os.ReadFile(packName)
	os.Remove(packName)

	content := base64.RawStdEncoding.EncodeToString(raw)
	finfo := protocols.FileParams{
		PackName:   packName,
		TargetPath: fsParams.TargetPath,
		FileBuf:    content,
	}
	payload, _ := config.Jsoner.Marshal(&finfo)

	taskid := model.BrainTaskManager.CreateTask(len(fsParams.TargetNodes))
	ctx.String(http.StatusAccepted, taskid)

	go func() {
		for i := range fsParams.TargetNodes {
			go pushFile(taskid, fsParams.TargetNodes[i], payload)
		}
	}()
}

func pushFile(taskid string, name string, payload []byte) {
	result := protocols.ExecutionResult{
		Name: name,
		Code: protocols.ExecOK,
	}
	subtask_id := model.TaskIdGen() // just random
	model.BrainTaskManager.AddSubTask(taskid, subtask_id, &result)
	defer model.BrainTaskManager.DoneSubTask(taskid, subtask_id, &result)

	raw, err := model.Request(name, protocols.TypeFilePush, payload)
	if err != nil {
		logger.Comm.Println("TypeFilePushResponse", err)
		result.Code = protocols.ExecCommunicationError
		result.CommunicationErrorMsg = "TypeFilePushResponse"
		return
	}
	var rmsg protocols.Result
	err = config.Jsoner.Unmarshal(raw, &rmsg)
	if err != nil {
		logger.Exceptions.Println("UnmarshalRmsg", err)
		result.Code = protocols.ExecCommunicationError
		result.CommunicationErrorMsg = "BrainError"
		return
	}
	if rmsg.Rmsg != "OK" {
		result.Code = protocols.ExecProcessError
		result.ProcessErrorMsg = rmsg.Rmsg
		logger.Exceptions.Print(rmsg.Rmsg)
	} else {
		result.Result = "OK"
	}
}

func FileTree(ctx *gin.Context) {
	rmsg := protocols.Result{
		Rmsg: "OK",
	}
	name, ok := ctx.GetQuery("name")
	if !ok {
		rmsg.Rmsg = "Lack name"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}
	subdir, _ := ctx.GetQuery("subdir")
	raw, err := getFileTree(name, subdir)
	if err != nil {
		rmsg.Rmsg = err.Error()
		ctx.JSON(http.StatusNotFound, rmsg)
		return
	}
	ctx.Data(http.StatusOK, "application/json", raw)
}

type ErrInvalidNode struct{ node string }

func (e ErrInvalidNode) Error() string { return fmt.Sprintf("Invalid node: %s\n", e.node) }

type ErrNetworkError struct{ node string }

func (e ErrNetworkError) Error() string { return fmt.Sprintf("Network error: %s\n", e.node) }

func getFileTree(name string, subdir string) ([]byte, error) {
	if name == "brain" {
		subdir = config.ParsePathWithEnv(subdir)
		return allFiles(subdir), nil
	}

	params, _ := config.Jsoner.Marshal(&protocols.FileParams{
		TargetPath: subdir,
	})
	raw, err := model.Request(name, protocols.TypeFileTree, params)
	if err != nil {
		logger.Comm.Print("FileTree")
		return nil, ErrNetworkError{node: name}
	}
	return raw, nil
}

type FileInfo struct {
	Name       string
	Size       int64
	ModifyTime string
}

func allFiles(path string) []byte {
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	finfos := []FileInfo{}
	walkDir(path, &finfos)
	hideRoot(path, &finfos)
	serialized, _ := config.Jsoner.Marshal(&finfos)
	return serialized
}

func walkDir(path string, files *[]FileInfo) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return
	}

	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.Name()[0] == '.' {
			continue
		} else if fi.IsDir() {
			walkDir(path+PthSep+fi.Name(), files)
		} else {
			detail, _ := fi.Info()
			modtimestr := detail.ModTime().Format("01月02日 15:04")
			finfo := FileInfo{
				Name:       path + PthSep + fi.Name(),
				Size:       detail.Size(),
				ModifyTime: modtimestr,
			}

			*files = append(*files, finfo)
		}
	}
}

func hideRoot(root string, files *[]FileInfo) {
	for i := range *files {
		(*files)[i].Name = (*files)[i].Name[len(root):]
	}
}

func FileDistrib(ctx *gin.Context) {
	packfiles, _ := ctx.FormFile("packfiles")
	packName := packfiles.Filename
	targetPath := ctx.PostForm("targetPath")
	isForceStr := ctx.PostForm("isForce")
	isForce := false
	if isForceStr == "true" {
		isForce = true
	}
	targetNodes := ctx.PostForm("targetNodes")
	rmsg := protocols.Result{
		Rmsg: "OK",
	}

	nodes := []string{}
	err := config.Jsoner.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Rmsg = "targetNodes:" + err.Error()
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	if !workgroup.IsInScope(ctx.GetStringMapString("octopoda_scope"), nodes...) {
		rmsg.Rmsg = "ERROR: some nodes are invalid or out of scope."
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	multipart, err := packfiles.Open()
	if err != nil {
		rmsg.Rmsg = "Open:" + err.Error()
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}
	defer multipart.Close()

	raw, err := io.ReadAll(multipart)
	if err != nil {
		rmsg.Rmsg = "ReadAll:" + err.Error()
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}
	content := b64Encode(raw)
	// content := base64.RawStdEncoding.EncodeToString(raw)
	finfo := protocols.FileParams{
		PackName:    packName,
		TargetPath:  targetPath,
		FileBuf:     content,
		ForceCreate: isForce,
	}
	payload, _ := config.Jsoner.Marshal(&finfo)

	taskid := model.BrainTaskManager.CreateTask(len(nodes))
	ctx.String(http.StatusAccepted, taskid)

	go func() {
		for i := range nodes {
			go pushFile(taskid, nodes[i], payload)
		}
	}()
}

func b64Encode(raw []byte) string {
	var buffer strings.Builder
	Offset := 0
	Len := len(raw)
	ChunkSize := 4096 * 3

	buffer.Grow(base64.RawStdEncoding.EncodedLen(Len))
	for Offset < Len {
		end := Offset + ChunkSize
		if Offset+ChunkSize > Len {
			end = Len
		}
		buffer.WriteString(base64.RawStdEncoding.EncodeToString(raw[Offset:end]))
		Offset += ChunkSize
	}
	return buffer.String()
}

func FilePull(ctx *gin.Context) {
	rmsg := protocols.Result{
		Rmsg: "OK",
	}
	name, ok := ctx.GetQuery("name")
	if !ok {
		rmsg.Rmsg = "Lack name"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}
	if !workgroup.IsInScope(ctx.GetStringMapString("octopoda_scope"), name) {
		rmsg.Rmsg = "ERROR: some nodes are invalid or out of scope."
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	fileOrDir, ok := ctx.GetQuery("fileOrDir")
	if !ok {
		rmsg.Rmsg = "Lack fileOrDir"
		ctx.JSON(http.StatusBadRequest, rmsg)
		return
	}

	// pull brain file?
	if name == "brain" {
		fileOrDir = config.ParsePathWithEnv(fileOrDir)
		_, err := os.Stat(fileOrDir)
		if err != nil {
			rmsg.Rmsg = "file or path not found"
			ctx.JSON(http.StatusNotFound, rmsg)
			return
		}
		// pack the file or dir
		packName := packFile(fileOrDir)
		if packName == "" {
			rmsg.Rmsg = "Error when packing files"
			ctx.JSON(http.StatusInternalServerError, rmsg)
			return
		}
		defer os.Remove(packName)
		rmsg.Output = loadFile(packName)
		ctx.JSON(http.StatusOK, rmsg)
		return
	}

	taskid := model.BrainTaskManager.CreateTask(1)
	ctx.String(http.StatusAccepted, taskid)

	go func() {
		result := protocols.ExecutionResult{
			Name: name,
			Code: protocols.ExecOK,
		}

		params, _ := config.Jsoner.Marshal(&protocols.FileParams{
			TargetPath: fileOrDir,
		})
		raw, err := model.Request(name, protocols.TypeFilePull, params)
		if err != nil {
			emsg := fmt.Sprintf("Send filepull request: %v", err)
			logger.Comm.Println(emsg)
			result.Code = protocols.ExecProcessError
			result.ProcessErrorMsg = emsg
			model.BrainTaskManager.AddFailedSubTask(taskid, model.TaskIdGen(), &result)
			return
		}

		if len(raw) == 0 {
			emsg := "file or path not found"
			logger.Comm.Println(emsg)
			result.Code = protocols.ExecProcessError
			result.ProcessErrorMsg = emsg
			model.BrainTaskManager.AddFailedSubTask(taskid, model.TaskIdGen(), &result)
			return
		}

		result.Result = string(raw)
		subtask_id := model.TaskIdGen()
		model.BrainTaskManager.AddSubTask(taskid, subtask_id, &result)
		defer model.BrainTaskManager.DoneSubTask(taskid, subtask_id, &result)
	}()
}

func packFile(fileOrDir string) string {
	packName := fmt.Sprintf("%d.zip", time.Now().Nanosecond())
	archiver.DefaultZip.OverwriteExisting = true
	err := archiver.DefaultZip.Archive([]string{fileOrDir}, packName)
	if err != nil {
		return ""
	}
	return packName
}

func loadFile(packName string) string {
	var filebufb64 strings.Builder
	f, err := os.Open(packName)
	if err != nil {
		return ""
	}
	defer f.Close()

	// prepare enough buffer capacity
	info, _ := f.Stat()
	filebufb64.Grow(base64.RawStdEncoding.EncodedLen(int(info.Size())))

	// read and encode to base64
	ChunkSize := 4096 * 4
	ChunkBuf := make([]byte, ChunkSize)
	for {
		n, err := f.Read(ChunkBuf)
		if err == io.EOF {
			break
		}
		filebufb64.WriteString(base64.RawStdEncoding.EncodeToString(ChunkBuf[:n]))
	}
	return filebufb64.String()
}
