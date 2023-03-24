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

type RMSG struct {
	Msg string
}

func FileUpload(ctx *gin.Context) {
	file, _ := ctx.FormFile("file")
	targetPath := ctx.PostForm("targetPath")
	rmsg := RMSG{"OK"}

	var sb strings.Builder
	sb.WriteString(config.GlobalConfig.Workspace.Root)
	sb.WriteString(targetPath)

	os.Mkdir(sb.String(), os.ModePerm)

	sb.WriteString(file.Filename)

	dst, err := os.Create(sb.String())
	if err != nil {
		logger.Brain.Println("FileCreate")
		rmsg.Msg = "FileCreate"
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

func FileSpread(ctx *gin.Context) {
	var fsParams FileSpreadParams
	err := ctx.ShouldBind(&fsParams)
	rmsg := RMSG{"OK"}

	if err != nil {
		logger.Brain.Println("FileCreate")
		rmsg.Msg = "FileCreate"
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
		rmsg.Msg = "OpenFile"
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
	var wg sync.WaitGroup

	for i := range fsParams.TargetNodes {
		name := fsParams.TargetNodes[i]
		if addr, exists := model.GetNodeAddress(name); exists {
			wg.Add(1)
			go pushFile(addr, payload, &wg)
		}
	}
	wg.Wait()
	ctx.JSON(200, rmsg)
}

func pushFile(addr string, payload []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	if conn, err := net.Dial("tcp", addr); err != nil {
		return
	} else {
		defer conn.Close()

		message.SendMessage(conn, message.TypeFilePush, payload)
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeFilePushResponse {
			logger.Tentacle.Println("TypeFilePushResponse", err)
			return
		}

		var rmsg RMSG
		err = json.Unmarshal(raw, &rmsg)
		if err != nil {
			logger.Tentacle.Println("UnmarshalNodeState", err)
			return
		}
		logger.Tentacle.Print(rmsg.Msg)
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
	if path[len(path) - 1] == '/' {
		path = path[:len(path) - 1]
	}
	finfos := []FileInfo{}
	walkDir(path, &finfos)
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
