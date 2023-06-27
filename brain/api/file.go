package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"brain/rdb"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver/v3"
)

type FileParams struct {
	PackName   string
	PathType   string
	TargetPath string
	FileBuf    string
}

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
	rmsg := message.Result{
		Rmsg: "OK",
	}

	// tmpPath: the zip file will be unpacked under this path
	var tmpSb strings.Builder
	tmpSb.WriteString(config.GlobalConfig.Workspace.Store)
	tmpSb.WriteString(".octopoda_tmp/")

	tmpPath := tmpSb.String()

	// dstPath: the unpacked files will be moved under this path
	fmt.Println("targetPath:", targetPath)
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
			ctx.JSON(403, rmsg)
			return
		}
	}
	tmpExtDir.WriteString("/*")

	tmpDst, err := os.Create(tmpSb.String())
	if err != nil {
		logger.Exceptions.Println("FileCreate")
		rmsg.Rmsg = "FileCreate:" + err.Error()
		ctx.JSON(403, rmsg)

		os.RemoveAll(tmpPath)
		return
	}

	src, _ := tarfile.Open()

	io.Copy(tmpDst, src)
	tmpDst.Close()
	src.Close()

	// fast return
	taskid := rdb.TaskIdGen()
	if !rdb.TaskNew(taskid, 3600) {
		logger.Exceptions.Print("TaskNew")
	}
	ctx.String(202, taskid)

	go func() {
		archiver.DefaultZip.OverwriteExisting = true
		err = archiver.DefaultZip.Unarchive(tmpSb.String(), tmpPath)
		defer os.RemoveAll(tmpPath)

		if err != nil {
			logger.Exceptions.Println("Unarchive")
			rmsg.Rmsg = "Unarchive:" + err.Error()
			// ctx.JSON(403, rmsg)
			if !rdb.TaskMarkFailed(taskid, rmsg, 3600) {
				logger.Exceptions.Print("TaskMarkFailed Error")
			}
			return
		}

		cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("mkdir -p %s && cp -r %s %s", dstPath, tmpExtDir.String(), dstPath))
		err = cmd.Run()
		if err != nil {
			logger.Exceptions.Println("cp -r")
			rmsg.Rmsg = "cp -r error:" + fmt.Sprintf("cp -r %s %s", tmpExtDir.String(), dstPath)
			if !rdb.TaskMarkFailed(taskid, rmsg, 3600) {
				logger.Exceptions.Print("TaskMarkFailed Error")
			}
			return
		}

		// ctx.JSON(200, rmsg)
		if !rdb.TaskMarkDone(taskid, rmsg, 3600) {
			logger.Exceptions.Print("TaskMarkDone Error")
		}
	}()
}

type FileSpreadParams struct {
	TargetPath  string
	FileOrDir   string
	TargetNodes []string
}

type BasicNodeResults struct {
	Name   string
	Result string
}

// need fix
func FileSpread(ctx *gin.Context) {
	var fsParams FileSpreadParams
	err := ctx.ShouldBind(&fsParams)
	rmsg := message.Result{
		Rmsg: "OK",
	}

	if err != nil {
		logger.Exceptions.Println("FileCreate")
		rmsg.Rmsg = "FileCreate:" + err.Error()
		ctx.JSON(403, rmsg)
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
	// else if fsParams.TargetPath[len(fsParams.TargetPath)-1] != '/' {
	// 	logger.Exceptions.Println("Invalid targetPath")
	// 	rmsg.Rmsg = "Invalid targetPath:" + fsParams.TargetPath
	// 	ctx.JSON(400, rmsg)
	// 	return
	// }

	// var sb strings.Builder
	// sb.WriteString(config.GlobalConfig.Workspace.Store)
	// sb.WriteString(fsParams.FileOrDir)
	// fname := sb.String()

	fname := pathFixing(fsParams.FileOrDir, config.GlobalConfig.Workspace.Store)
	// fmt.Println(fname, "----")

	// we dont check because fname may be a pattern
	// _, err = os.Stat(fname)
	// if err != nil {
	// 	logger.Exceptions.Println("FileOrDir Not Found")
	// 	rmsg.Rmsg = "FileOrDir Not Found:" + err.Error()
	// 	ctx.JSON(403, rmsg)
	// 	return
	// }

	// fast return
	taskid := rdb.TaskIdGen()
	if !rdb.TaskNew(taskid, 3600) {
		logger.Exceptions.Print("TaskNew")
	}
	ctx.String(202, taskid)

	go func() {
		// if fname[len(fname)-1] == '/' {
		// 	fname = fname + "."
		// }

		// packName := fmt.Sprintf("%d.zip", time.Now().Nanosecond())
		// // err = exec.Command("tar", "-cf", tarName, "-C", filepath.Dir(fname), filepath.Base(fname)).Run()
		// archiver.DefaultZip.OverwriteExisting = true
		// err = archiver.DefaultZip.Archive([]string{fname}, packName)

		wrapName := fmt.Sprintf("%d.wrap", time.Now().Nanosecond())
		os.Mkdir(wrapName, os.ModePerm)
		cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -r %s %s", fname, wrapName))
		err := cmd.Run()
		if err != nil {
			logger.Exceptions.Print("Wrap files: " + fname + "-->" + wrapName + " | " + cmd.String())
		}
		defer os.RemoveAll(wrapName)

		packName := fmt.Sprintf("%s.zip", wrapName)
		// err := exec.Command("tar", "-cf", tarName, "-C", filepath.Dir(localFileOrDir), filepath.Base(localFileOrDir)).Run()
		// if err != nil {
		// 	output.PrintFatal("cmd.Run")
		// }
		archiver.DefaultZip.OverwriteExisting = true
		err = archiver.DefaultZip.Archive([]string{wrapName}, packName)
		if err != nil {
			logger.Exceptions.Println("Archive")
			rmsg.Rmsg = "Archive:" + err.Error()
			if !rdb.TaskMarkFailed(taskid, rmsg, 3600) {
				logger.Exceptions.Print("TaskMarkFailed Error")
			}
			return
		}

		raw, _ := os.ReadFile(packName)
		os.Remove(packName)

		content := base64.RawStdEncoding.EncodeToString(raw)
		finfo := FileParams{
			PackName:   packName,
			TargetPath: fsParams.TargetPath,
			FileBuf:    content,
		}
		payload, _ := config.Jsoner.Marshal(&finfo)

		// check target nodes
		// spread file
		results := make([]BasicNodeResults, len(fsParams.TargetNodes))
		var wg sync.WaitGroup

		for i := range fsParams.TargetNodes {
			name := fsParams.TargetNodes[i]
			results[i].Name = name
			wg.Add(1)
			go pushFile(name, payload, &wg, &results[i].Result)
		}
		wg.Wait()
		if !rdb.TaskMarkDone(taskid, results, 3600) {
			logger.Exceptions.Print("TaskMarkDone Error")
		}
	}()
}

func pushFile(name string, payload []byte, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "UnknownError"

	raw, err := model.Request(name, message.TypeFilePush, payload)
	if err != nil {
		logger.Comm.Println("TypeFilePushResponse", err)
		*result = "TypeFilePushResponse"
		return
	}
	var rmsg message.Result
	err = config.Jsoner.Unmarshal(raw, &rmsg)
	if err != nil {
		logger.Exceptions.Println("UnmarshalRmsg", err)
		*result = "MasterError"
		return
	}
	if rmsg.Rmsg != "OK" {
		logger.Exceptions.Print(rmsg.Rmsg)
		*result = rmsg.Rmsg
	} else {
		*result = "OK"
	}
}

func FileTree(ctx *gin.Context) {
	rmsg := message.Result{
		Rmsg: "OK",
	}
	name, ok := ctx.GetQuery("name")
	if !ok {
		rmsg.Rmsg = "Lack name"
		ctx.JSON(400, rmsg)
		return
	}
	pathtype, ok := ctx.GetQuery("pathtype")
	if !ok {
		rmsg.Rmsg = "Lack pathtype"
		ctx.JSON(400, rmsg)
		return
	}
	subdir, _ := ctx.GetQuery("subdir")
	raw, err := getFileTree(pathtype, name, subdir)
	if err != nil {
		rmsg.Rmsg = err.Error()
		ctx.JSON(404, rmsg)
		return
	}
	// rmsg.Output = string(raw) //  marshal?
	// ctx.JSON(200, rmsg)
	ctx.Data(200, "application/json", raw)
}

type ErrInvalidPathType struct{ pathtype, node string }

func (e ErrInvalidPathType) Error() string {
	return fmt.Sprintf("Invalid path type: %s on %s\n", e.pathtype, e.node)
}

type ErrInvalidNode struct{ node string }

func (e ErrInvalidNode) Error() string { return fmt.Sprintf("Invalid node: %s\n", e.node) }

type ErrNetworkError struct{ node string }

func (e ErrNetworkError) Error() string { return fmt.Sprintf("Network error: %s\n", e.node) }

func getFileTree(pathtype string, name string, subdir string) ([]byte, error) {
	var pathsb strings.Builder
	if name == "master" {
		switch pathtype {
		case "store":
			pathsb.WriteString(config.GlobalConfig.Workspace.Store)
		case "log":
			pathsb.WriteString(config.GlobalConfig.Logger.Path)
		default:
			return nil, ErrInvalidPathType{pathtype: pathtype, node: name}
		}
		pathsb.WriteString(subdir)
		return allFiles(pathsb.String()), nil
	}

	params, _ := config.Jsoner.Marshal(&FileParams{
		PathType:   pathtype,
		TargetPath: subdir,
	})
	raw, err := model.Request(name, message.TypeFileTree, params)
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
	targetNodes := ctx.PostForm("targetNodes")
	rmsg := message.Result{
		Rmsg: "OK",
	}

	nodes := []string{}
	err := config.Jsoner.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Rmsg = "targetNodes:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}

	multipart, err := packfiles.Open()
	if err != nil {
		rmsg.Rmsg = "Open:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}
	defer multipart.Close()

	raw, err := io.ReadAll(multipart)
	if err != nil {
		rmsg.Rmsg = "ReadAll:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}

	// fast return
	taskid := rdb.TaskIdGen()
	if !rdb.TaskNew(taskid, 3600) {
		logger.Exceptions.Print("TaskNew")
	}
	ctx.String(202, taskid)

	go func() {
		content := b64Encode(raw)
		// content := base64.RawStdEncoding.EncodeToString(raw)
		finfo := FileParams{
			PackName:   packName,
			TargetPath: targetPath,
			FileBuf:    content,
		}
		payload, _ := config.Jsoner.Marshal(&finfo)

		// check target nodes
		// spread file
		results := make([]BasicNodeResults, len(nodes))
		var wg sync.WaitGroup

		for i := range nodes {
			name := nodes[i]
			results[i].Name = name
			// if addr, exists := model.GetNodeAddress(name); exists {

			// } else {
			// 	results[i].Result = "NodeNotExists"
			// }
			wg.Add(1)
			go pushFile(name, payload, &wg, &results[i].Result)
		}
		wg.Wait()
		if !rdb.TaskMarkDone(taskid, results, 3600) {
			logger.Exceptions.Print("TaskMarkDone Error")
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
	rmsg := message.Result{
		Rmsg: "OK",
	}
	name, ok := ctx.GetQuery("name")
	if !ok {
		rmsg.Rmsg = "Lack name"
		ctx.JSON(400, rmsg)
		return
	}
	pathtype, ok := ctx.GetQuery("pathtype")
	if !ok {
		rmsg.Rmsg = "Lack pathtype"
		ctx.JSON(400, rmsg)
		return
	}
	fileOrDir, ok := ctx.GetQuery("fileOrDir")
	if !ok {
		rmsg.Rmsg = "Lack fileOrDir"
		ctx.JSON(400, rmsg)
		return
	}

	// pull master file?
	if name == "master" {
		var pathsb strings.Builder
		switch pathtype {
		case "store":
			pathsb.WriteString(config.GlobalConfig.Workspace.Store)
		case "log":
			pathsb.WriteString(config.GlobalConfig.Logger.Path)
		default:
			rmsg.Rmsg = ErrInvalidPathType{pathtype: pathtype, node: name}.Error()
			ctx.JSON(400, rmsg)
			return
		}
		pathsb.WriteString(fileOrDir)
		_, err := os.Stat(pathsb.String())
		if err != nil {
			rmsg.Rmsg = "file or path not found"
			ctx.JSON(404, rmsg)
			return
		}
		// pack the file or dir
		packName := packFile(pathsb.String())
		if packName == "" {
			rmsg.Rmsg = "Error when packing files"
			ctx.JSON(500, rmsg)
			return
		}
		defer os.Remove(packName)
		rmsg.Output = loadFile(packName)
		ctx.JSON(200, rmsg)
		return
	}

	// fast return
	taskid := rdb.TaskIdGen()
	if !rdb.TaskNew(taskid, 3600) {
		logger.Exceptions.Print("TaskNew")
	}
	ctx.String(202, taskid)

	go func() {
		params, _ := config.Jsoner.Marshal(&FileParams{
			PathType:   pathtype,
			TargetPath: fileOrDir,
		})
		raw, err := model.Request(name, message.TypeFilePull, params)
		if err != nil {
			logger.Comm.Print("TypeFilePullResponse")
			rmsg.Rmsg = ErrNetworkError{node: name}.Error()
			if !rdb.TaskMarkFailed(taskid, rmsg, 3600) {
				logger.Exceptions.Print("TaskMarkFailed Error")
			}
			return
		}

		if len(raw) == 0 {
			rmsg.Rmsg = "file or path not found"
			if !rdb.TaskMarkFailed(taskid, rmsg, 3600) {
				logger.Exceptions.Print("TaskMarkFailed Error")
			}
			return
		}

		// rmsg.Output = encodeBuf(raw)
		rmsg.Output = string(raw)
		if !rdb.TaskMarkDone(taskid, rmsg, 3600) {
			logger.Exceptions.Print("TaskMarkDone Error")
		}
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

// func encodeBuf(buf []byte) string {
// 	var bufb64 strings.Builder

// 	// prepare enough buffer capacity
// 	bufb64.Grow(base64.RawStdEncoding.EncodedLen(len(buf)))

// 	// read and encode to base64
// 	Offset := 0
// 	Len := len(buf)
// 	ChunkSize := 4096 * 4
// 	for Offset < Len {
// 		end := Offset + ChunkSize
// 		if end > Len {
// 			end = Len
// 		}
// 		bufb64.WriteString(base64.RawStdEncoding.EncodeToString(buf[Offset:end]))
// 		Offset += ChunkSize
// 	}
// 	return bufb64.String()
// }
