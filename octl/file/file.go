package file

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"octl/config"
	"octl/output"
	"octl/task"
	"os"
	"strings"
	"time"

	"github.com/mholt/archiver"
)

func UpLoadFile(localFileOrDir string, targetPath string) {
	if targetPath == "." {
		targetPath = ""
	} else if targetPath[len(targetPath)-1] != '/' {
		targetPath = targetPath + "/"
	}

	if localFileOrDir[len(localFileOrDir)-1] == '/' {
		localFileOrDir = localFileOrDir[:len(localFileOrDir)-1]
	}

	packName := fmt.Sprintf("%d.zip", time.Now().Nanosecond())
	// err := exec.Command("tar", "-cf", packName, "-C", filepath.Dir(localFileOrDir), filepath.Base(localFileOrDir)).Run()

	// if err != nil {
	// 	output.PrintFatal("cmd.Run")
	// }
	archiver.DefaultZip.OverwriteExisting = true
	err := archiver.DefaultZip.Archive([]string{localFileOrDir}, packName)
	if err != nil {
		output.PrintFatal("Archive")
	}
	defer os.Remove(packName)

	f, err := os.OpenFile(packName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		output.PrintFatal("err")
	}
	defer f.Close()

	bodyBuffer := bytes.Buffer{}
	bodyWriter := multipart.NewWriter(&bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("tarfile", packName)
	io.Copy(fileWriter, f)

	bodyWriter.WriteField("targetPath", targetPath)

	contentType := bodyWriter.FormDataContentType()

	bodyWriter.Close()

	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.FileUpload,
	)
	res, err := http.Post(url, contentType, &bodyBuffer)
	if err != nil {
		output.PrintFatal("post")
	}

	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		output.PrintFatal("ReadAll")
	}

	if res.StatusCode != 202 {
		output.PrintFatal("Request submit error: " + string(msg))
		return
	}
	results, err := task.WaitTask("UPLOADING...", string(msg))
	if err != nil {
		output.PrintFatal("Task processing error: " + err.Error())
		return
	}
	output.PrintJSON(results)
}

type FileSpreadParams struct {
	TargetPath  string
	FileOrDir   string
	TargetNodes []string
}

func SpreadFile(FileOrDir string, targetPath string, nodes []string) {
	fsParams := &FileSpreadParams{
		TargetPath:  targetPath,
		FileOrDir:   FileOrDir,
		TargetNodes: nodes,
	}
	buf, _ := config.Jsoner.Marshal(fsParams)

	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.FileSpread,
	)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		output.PrintFatal("Post")
	}
	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		output.PrintFatal("ReadAll")
	}

	if res.StatusCode != 202 {
		output.PrintFatal("Request submit error: " + string(msg))
		return
	}
	results, err := task.WaitTask("SPREADING...", string(msg))
	if err != nil {
		output.PrintFatal("Task processing error: " + err.Error())
		return
	}
	output.PrintJSON(results)
}

type FileDistribParams struct {
	LocalFile   string
	TargetPath  string
	TargetNodes []string
}

func DistribFile(localFileOrDir string, targetPath string, nodes []string) {
	if targetPath == "." {
		targetPath = ""
	} else if targetPath[len(targetPath)-1] != '/' {
		targetPath = targetPath + "/"
	}

	if localFileOrDir[len(localFileOrDir)-1] == '/' {
		localFileOrDir = localFileOrDir[:len(localFileOrDir)-1]
	}

	packName := fmt.Sprintf("%d.zip", time.Now().Nanosecond())
	// err := exec.Command("tar", "-cf", tarName, "-C", filepath.Dir(localFileOrDir), filepath.Base(localFileOrDir)).Run()
	// if err != nil {
	// 	output.PrintFatal("cmd.Run")
	// }
	archiver.DefaultZip.OverwriteExisting = true
	err := archiver.DefaultZip.Archive([]string{localFileOrDir}, packName)
	if err != nil {
		output.PrintFatal("Archive")
	}
	defer os.Remove(packName)

	f, err := os.OpenFile(packName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		output.PrintFatal("err")
	}
	defer f.Close()

	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)

	bodyBuffer := bytes.Buffer{}
	bodyWriter := multipart.NewWriter(&bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("packfiles", packName)
	io.Copy(fileWriter, f)
	bodyWriter.WriteField("targetPath", targetPath)
	bodyWriter.WriteField("targetNodes", string(nodes_serialized))

	contentType := bodyWriter.FormDataContentType()

	bodyWriter.Close()

	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.FileDistrib,
	)

	res, err := http.Post(url, contentType, &bodyBuffer)
	if err != nil {
		output.PrintFatal("post")
	}
	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		output.PrintFatal("ReadAll")
	}

	if res.StatusCode != 202 {
		output.PrintFatal("Request submit error: " + string(msg))
		return
	}
	results, err := task.WaitTask("DISTRIBUTING...", string(msg))
	if err != nil {
		output.PrintFatal("Task processing error: " + err.Error())
		return
	}
	output.PrintJSON(results)
}

func ListAllFile(pathtype string, node string, subdir string) {
	url := fmt.Sprintf("http://%s:%d/%s%s?pathtype=%s&name=%s&subdir=%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.FileTree,
		pathtype,
		node,
		subdir,
	)
	res, err := http.Get(url)
	if err != nil {
		output.PrintFatal("Get")
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
	if err != nil {
		output.PrintFatal("ReadAll")
	}
	output.PrintJSON(msg)
}

type Result struct {
	// Result code. Reserved
	Rcode int

	// Result message. OK, or Error Reason
	Rmsg string

	// For script or command execution. Script output.
	Output string

	// For version control. Version hash code.
	Version string

	// For version control. Modified flag.
	Modified bool
}

func PullFile(pathtype string, node string, fileOrDir string, targetdir string) {
	url := fmt.Sprintf("http://%s:%d/%s%s?pathtype=%s&name=%s&fileOrDir=%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.FilePull,
		pathtype,
		node,
		fileOrDir,
	)
	res, err := http.Get(url)
	if err != nil {
		output.PrintFatal("Get")
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
	if err != nil {
		output.PrintFatal("ReadAll")
	}

	result := Result{}
	if res.StatusCode == 200 {
		// directly get the file from master
		config.Jsoner.Unmarshal(msg, &result)

		// unpack result.Output
		err = unpackFiles(result.Output, targetdir)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Success")
	} else if res.StatusCode == 202 {
		// have to wait
		resultmsg, err := task.WaitTask("PULLING...", string(msg))
		if err != nil {
			output.PrintFatal("Task processing error: " + err.Error())
			return
		}
		config.Jsoner.Unmarshal(resultmsg, &result)
		if len(result.Output) == 0 {
			output.PrintJSON(resultmsg)
		} else {
			// unpack result.Output
			err = unpackFiles(result.Output, targetdir)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("Success")
		}
	} else {
		// some error
		output.PrintJSON(msg)
	}
}

type ErrUnpack struct{}

func (ErrUnpack) Error() string { return "ErrUnpack" }

func unpackFiles(packb64 string, dir string) error {
	var file strings.Builder
	// tmp tar file name
	fname := fmt.Sprintf("%d.zip", time.Now().Nanosecond())

	// make complete path and filename
	// file.WriteString(config.GlobalConfig.Workspace.Root)
	file.WriteString(dir)
	path := file.String()
	os.Mkdir(path, os.ModePerm)
	if path[len(path)-1] != '/' {
		file.WriteByte('/')
	}
	file.WriteString(fname)

	err := saveFile(packb64, file.String())
	if err != nil {
		return err
	}

	archiver.DefaultZip.OverwriteExisting = true
	err = archiver.DefaultZip.Unarchive(file.String(), path)
	if err != nil {
		return ErrUnpack{}
	}

	os.Remove(file.String())
	return nil
}

func saveFile(filebufb64, filename string) error {
	Offset := 0
	Len := len(filebufb64)
	ChunkSize := 4096 * 4
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	for Offset < Len {
		end := Offset + ChunkSize
		if Offset+ChunkSize > Len {
			end = Len
		}
		content, err := base64.RawStdEncoding.DecodeString(filebufb64[Offset:end])
		if err != nil {
			return err
		}
		_, err = f.Write(content)
		if err != nil {
			return err
		}
		Offset += ChunkSize
	}
	return nil
}
