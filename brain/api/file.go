package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"brain/rdb"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type FileParams struct {
	TarName    string
	TargetPath string
	FileBuf    string
}

func FileUpload(ctx *gin.Context) {
	tarfile, _ := ctx.FormFile("tarfile")
	targetPath := ctx.PostForm("targetPath")
	rmsg := message.Result{
		Rmsg: "OK",
	}

	var sb strings.Builder
	sb.WriteString(config.GlobalConfig.Workspace.Store)
	sb.WriteString(targetPath)

	path := sb.String()
	os.Mkdir(path, os.ModePerm)

	sb.WriteString(tarfile.Filename)
	dst, err := os.Create(sb.String())
	if err != nil {
		logger.Brain.Println("FileCreate")
		rmsg.Rmsg = "FileCreate:" + err.Error()
		ctx.JSON(403, rmsg)
		return
	}

	src, _ := tarfile.Open()

	io.Copy(dst, src)
	dst.Close()
	src.Close()

	// fast return
	taskid := rdb.TaskIdGen()
	if !rdb.TaskNew(taskid, 3600) {
		logger.Brain.Print("TaskNew")
	}
	ctx.String(202, taskid)

	go func() {
		err = exec.Command("tar", "-xf", sb.String(), "-C", path).Run()
		os.Remove(sb.String())
		if err != nil {
			logger.Brain.Println("UnpackFile")
			rmsg.Rmsg = "UnpackFile:" + err.Error()
			// ctx.JSON(403, rmsg)
			if !rdb.TaskMarkFailed(taskid, rmsg, 3600) {
				logger.Brain.Print("TaskMarkFailed Error")
			}
			return
		}

		// ctx.JSON(200, rmsg)
		if !rdb.TaskMarkDone(taskid, rmsg, 3600) {
			logger.Brain.Print("TaskMarkDone Error")
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
		logger.Brain.Println("FileCreate")
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
	} else if fsParams.TargetPath[len(fsParams.TargetPath)-1] != '/' {
		logger.Brain.Println("Invalid targetPath")
		rmsg.Rmsg = "Invalid targetPath:" + fsParams.TargetPath
		ctx.JSON(400, rmsg)
		return
	}

	var sb strings.Builder
	sb.WriteString(config.GlobalConfig.Workspace.Store)
	sb.WriteString(fsParams.FileOrDir)
	fname := sb.String()
	
	_, err = os.Stat(fname)
	if err != nil {
		logger.Brain.Println("FileOrDir Not Found")
		rmsg.Rmsg = "FileOrDir Not Found:" + err.Error()
		ctx.JSON(403, rmsg)
		return
	}

	// fast return
	taskid := rdb.TaskIdGen()
	if !rdb.TaskNew(taskid, 3600) {
		logger.Brain.Print("TaskNew")
	}
	ctx.String(202, taskid)

	go func() {
		tarName := fmt.Sprintf("%d.tar", time.Now().Nanosecond())
		err = exec.Command("tar", "-cf", tarName, "-C", filepath.Dir(fname), filepath.Base(fname)).Run()
		if err != nil {
			logger.Brain.Println("PackFile")
			rmsg.Rmsg = "PackFile:" + err.Error()
			if !rdb.TaskMarkFailed(taskid, rmsg, 3600) {
				logger.Brain.Print("TaskMarkFailed Error")
			}
			return
		}

		raw, _ := os.ReadFile(tarName)
		os.Remove(tarName)

		content := base64.RawStdEncoding.EncodeToString(raw)
		finfo := FileParams{
			TarName:    tarName,
			TargetPath: fsParams.TargetPath,
			FileBuf:    content,
		}
		payload, _ := json.Marshal(&finfo)

		// check target nodes
		// spread file
		results := make([]BasicNodeResults, len(fsParams.TargetNodes))
		var wg sync.WaitGroup

		for i := range fsParams.TargetNodes {
			name := fsParams.TargetNodes[i]
			results[i].Name = name
			if addr, exists := model.GetNodeAddress(name); exists {
				wg.Add(1)
				go pushFile(addr, payload, &wg, &results[i].Result)
			} else {
				results[i].Result = "NodeNotExists"
			}
		}
		wg.Wait()
		if !rdb.TaskMarkDone(taskid, results, 3600) {
			logger.Brain.Print("TaskMarkDone Error")
		}
	}()
}

func pushFile(addr string, payload []byte, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "UnknownError"

	if conn, err := net.Dial("tcp", addr); err != nil {
		return
	} else {
		defer conn.Close()

		message.SendMessage(conn, message.TypeFilePush, payload)
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeFilePushResponse {
			logger.Tentacle.Println("TypeFilePushResponse", err)
			*result = "NetError"
			return
		}

		var rmsg message.Result
		err = json.Unmarshal(raw, &rmsg)
		if err != nil {
			logger.Tentacle.Println("UnmarshalNodeState", err)
			*result = "MasterError"
			return
		}
		logger.Tentacle.Print(rmsg.Rmsg)
		if rmsg.Rmsg != "OK" {
			*result = "NodeError"
		} else {
			*result = "OK"
		}
	}
}

func FileTree(ctx *gin.Context) {
	name, ok := ctx.GetQuery("name")
	if !ok {
		ctx.Status(400)
		return
	}
	subdir, _ := ctx.GetQuery("subdir")
	raw := getFileTree(name, subdir)
	ctx.Data(200, "application/json", raw)
}

func getFileTree(name string, subdir string) []byte {
	var pathsb strings.Builder
	if name == "master" {
		pathsb.WriteString(config.GlobalConfig.Workspace.Store)
		pathsb.WriteString(subdir)
		return allFiles(pathsb.String())
	}
	addr, ok := model.GetNodeAddress(name)
	if !ok {
		return []byte{}
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		logger.Tentacle.Print("FileTree Dial")
		return []byte{}
	}
	defer conn.Close()

	err = message.SendMessage(conn, message.TypeFileTree, []byte(subdir))
	if err != nil {
		logger.Tentacle.Print("FileTree")
		return []byte{}
	}
	mtype, raw, err := message.RecvMessage(conn)
	if err != nil || mtype != message.TypeFileTreeResponse {
		logger.Tentacle.Print("FileTreeResponse")
		return []byte{}
	}
	return raw
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
	serialized, _ := json.Marshal(&finfos)
	return serialized
}

func walkDir(path string, files *[]FileInfo) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return
	}

	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() {
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
	tarfile, _ := ctx.FormFile("tarfile")
	tarName := tarfile.Filename
	targetPath := ctx.PostForm("targetPath")
	targetNodes := ctx.PostForm("targetNodes")
	rmsg := message.Result{
		Rmsg: "OK",
	}

	nodes := []string{}
	err := json.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Rmsg = "targetNodes:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}

	multipart, err := tarfile.Open()
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
		logger.Brain.Print("TaskNew")
	}
	ctx.String(202, taskid)

	go func() {
		content := b64Encode(raw)
		// content := base64.RawStdEncoding.EncodeToString(raw)
		finfo := FileParams{
			TarName:    tarName,
			TargetPath: targetPath,
			FileBuf:    content,
		}
		payload, _ := json.Marshal(&finfo)

		// check target nodes
		// spread file
		results := make([]BasicNodeResults, len(nodes))
		var wg sync.WaitGroup

		for i := range nodes {
			name := nodes[i]
			results[i].Name = name
			if addr, exists := model.GetNodeAddress(name); exists {
				wg.Add(1)
				go pushFile(addr, payload, &wg, &results[i].Result)
			} else {
				results[i].Result = "NodeNotExists"
			}
		}
		wg.Wait()
		if !rdb.TaskMarkDone(taskid, results, 3600) {
			logger.Brain.Print("TaskMarkDone Error")
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
		if Offset + ChunkSize > Len {
			end = Len 
		}
		buffer.WriteString(base64.RawStdEncoding.EncodeToString(raw[Offset:end]))
		Offset += ChunkSize
	}
	return buffer.String()
}