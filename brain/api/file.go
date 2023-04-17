package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type FileParams struct {
	FileName   string
	TargetPath string
	FileBuf    string
}

func FileUpload(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	targetPath := ctx.PostForm("targetPath")
	rmsg := message.Result{
		Rmsg: "OK",
	}

	var sb strings.Builder
	sb.WriteString(config.GlobalConfig.Workspace.Root)
	sb.WriteString(targetPath)

	os.Mkdir(sb.String(), os.ModePerm)

	sb.WriteString(file.Filename)

	dst, err := os.Create(sb.String())
	if err != nil {
		logger.Brain.Println("FileCreate")
		rmsg.Rmsg = "FileCreate:" + err.Error()
		ctx.JSON(403, rmsg)
		return
	}
	defer dst.Close()

	src, _ := file.Open()
	io.Copy(dst, src)
	ctx.JSON(200, rmsg)
}

type FileSpreadParams struct {
	SourcePath  string
	TargetPath  string
	FileName    string
	TargetNodes []string
}

type BasicNodeResults struct {
	Name   string
	Result string
}

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
	var localFile strings.Builder
	localFile.WriteString(config.GlobalConfig.Workspace.Root)
	localFile.WriteString(fsParams.SourcePath)
	localFile.WriteString(fsParams.FileName)
	f, err := os.OpenFile(localFile.String(), os.O_RDONLY, os.ModePerm)
	if err != nil {
		logger.Brain.Println("OpenFile")
		rmsg.Rmsg = "OpenFile:" + err.Error()
		ctx.JSON(403, rmsg)
		return
	}
	defer f.Close()

	raw, _ := io.ReadAll(f)
	content := base64.RawStdEncoding.EncodeToString(raw)
	finfo := FileParams{
		FileName:   fsParams.FileName,
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
	ctx.JSON(200, results)
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
		pathsb.WriteString(config.GlobalConfig.Workspace.Root)
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
	file, _ := ctx.FormFile("file")
	fileName := file.Filename
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

	multipart, err := file.Open()
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

	content := base64.RawStdEncoding.EncodeToString(raw)
	finfo := FileParams{
		FileName:   fileName,
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
	ctx.JSON(200, results)
}
