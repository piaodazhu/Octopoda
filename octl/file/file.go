package file

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"octl/config"
	"octl/nameclient"
	"octl/output"
	"octl/task"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mholt/archiver/v3"
)

func UpLoadFile(localFileOrDir string, targetPath string) {
	if targetPath == "." {
		targetPath = ""
	} else if targetPath[len(targetPath)-1] != '/' {
		targetPath = targetPath + "/"
	}

	// if localFileOrDir[len(localFileOrDir)-1] == '/' {
	// 	localFileOrDir = localFileOrDir[:len(localFileOrDir)-1]
	// }
	pwd, _ := os.Getwd()
	srcPath := pathFixing(localFileOrDir, pwd+string(filepath.Separator))

	// wrap the files first
	wrapName := fmt.Sprintf("%d.wrap", time.Now().Nanosecond())
	os.Mkdir(wrapName, os.ModePerm)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe", "/C", fmt.Sprintf("cp -Force %s %s", srcPath, wrapName))
	} else {
		cmd = exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -r %s %s", srcPath, wrapName))
	}
	err := cmd.Run()
	if err != nil {
		output.PrintFatalln("Wrap files: " + srcPath + "-->" + wrapName + " | " + cmd.String())
	}
	defer os.RemoveAll(wrapName)

	packName := fmt.Sprintf("%s.zip", wrapName)
	// err := exec.Command("tar", "-cf", tarName, "-C", filepath.Dir(localFileOrDir), filepath.Base(localFileOrDir)).Run()
	// if err != nil {
	// 	output.PrintFatalln("cmd.Run")
	// }
	archiver.DefaultZip.OverwriteExisting = true
	err = archiver.DefaultZip.Archive([]string{wrapName}, packName)
	if err != nil {
		output.PrintFatalln("Archive")
	}
	defer os.Remove(packName)

	f, err := os.OpenFile(packName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		output.PrintFatalln("err")
	}
	defer f.Close()

	bodyBuffer := bytes.Buffer{}
	bodyWriter := multipart.NewWriter(&bodyBuffer)
	fileWriter, _ := bodyWriter.CreateFormFile("tarfile", packName)
	io.Copy(fileWriter, f)

	bodyWriter.WriteField("targetPath", targetPath)

	contentType := bodyWriter.FormDataContentType()

	bodyWriter.Close()

	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.FileUpload,
	)
	res, err := http.Post(url, contentType, &bodyBuffer)
	if err != nil {
		output.PrintFatalln("post")
	}

	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		output.PrintFatalln("ReadAll")
	}

	if res.StatusCode != 202 {
		output.PrintFatalln("Request submit error: " + string(msg))
		return
	}
	results, err := task.WaitTask("UPLOADING...", string(msg))
	if err != nil {
		output.PrintFatalln("Task processing error: " + err.Error())
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

	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.FileSpread,
	)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		output.PrintFatalln("Post")
	}
	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		output.PrintFatalln("ReadAll")
	}

	if res.StatusCode != 202 {
		output.PrintFatalln("Request submit error: " + string(msg))
		return
	}
	results, err := task.WaitTask("SPREADING...", string(msg))
	if err != nil {
		output.PrintFatalln("Task processing error: " + err.Error())
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

	// if localFileOrDir[len(localFileOrDir)-1] == '/' {
	// 	localFileOrDir = localFileOrDir[:len(localFileOrDir)-1]
	// }

	pwd, _ := os.Getwd()
	srcPath := pathFixing(localFileOrDir, pwd+string(filepath.Separator))

	// wrap the files first
	wrapName := fmt.Sprintf("%d.wrap", time.Now().Nanosecond())
	os.Mkdir(wrapName, os.ModePerm)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe", "/C", fmt.Sprintf("cp -Force %s %s", srcPath, wrapName))
	} else {
		cmd = exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -r %s %s", srcPath, wrapName))
	}
	err := cmd.Run()
	if err != nil {
		output.PrintFatalln("Wrap files: " + srcPath + "-->" + wrapName + " | " + cmd.String())
	}
	defer os.RemoveAll(wrapName)

	packName := fmt.Sprintf("%s.zip", wrapName)
	// err := exec.Command("tar", "-cf", tarName, "-C", filepath.Dir(localFileOrDir), filepath.Base(localFileOrDir)).Run()
	// if err != nil {
	// 	output.PrintFatalln("cmd.Run")
	// }
	archiver.DefaultZip.OverwriteExisting = true
	err = archiver.DefaultZip.Archive([]string{wrapName}, packName)
	if err != nil {
		output.PrintFatalln("Archive")
	}
	defer os.Remove(packName)

	f, err := os.OpenFile(packName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		output.PrintFatalln("err")
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

	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.FileDistrib,
	)

	res, err := http.Post(url, contentType, &bodyBuffer)
	if err != nil {
		output.PrintFatalln("post")
	}
	msg, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		output.PrintFatalln("ReadAll")
	}

	if res.StatusCode != 202 {
		output.PrintFatalln("Request submit error: " + string(msg))
		return
	}
	results, err := task.WaitTask("DISTRIBUTING...", string(msg))
	if err != nil {
		output.PrintFatalln("Task processing error: " + err.Error())
		return
	}
	output.PrintJSON(results)
}

func ListAllFile(pathtype string, node string, subdir string) {
	url := fmt.Sprintf("http://%s/%s%s?pathtype=%s&name=%s&subdir=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.FileTree,
		pathtype,
		node,
		subdir,
	)
	res, err := http.Get(url)
	if err != nil {
		output.PrintFatalln("Get")
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
	if err != nil {
		output.PrintFatalln("ReadAll")
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

type FilePullParams struct {
	PackName   string
	PathType   string
	TargetPath string
	FileBuf    string
}

func PullFile(pathtype string, node string, fileOrDir string, targetdir string) {
	url := fmt.Sprintf("http://%s/%s%s?pathtype=%s&name=%s&fileOrDir=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.FilePull,
		pathtype,
		node,
		fileOrDir,
	)
	res, err := http.Get(url)
	if err != nil {
		output.PrintFatalln("Get")
	}
	defer res.Body.Close()
	msg, err := io.ReadAll(res.Body)
	if err != nil {
		output.PrintFatalln("ReadAll")
	}

	result := Result{}
	finfo := FilePullParams{}
	if res.StatusCode == 200 {
		// get the file info structure from master
		err = config.Jsoner.Unmarshal(msg, &result)
		if err != nil {
			output.PrintFatalln(err.Error())
		}

		// marshal the file info
		err = config.Jsoner.Unmarshal([]byte(result.Output), &finfo)
		if err != nil {
			output.PrintFatalln(err.Error())
		}

		// unpack result.Output
		err = unpackFiles(finfo.FileBuf, finfo.PackName, targetdir)
		if err != nil {
			output.PrintFatalln(err.Error())
		}
		output.PrintInfoln("Success")
	} else if res.StatusCode == 202 {
		// have to wait
		resultmsg, err := task.WaitTask("PULLING...", string(msg))
		if err != nil {
			output.PrintFatalln("Task processing error: " + err.Error())
			return
		}
		config.Jsoner.Unmarshal(resultmsg, &result)
		if len(result.Output) == 0 {
			output.PrintJSON(resultmsg)
		} else {
			// marshal the file info
			err = config.Jsoner.Unmarshal([]byte(result.Output), &finfo)
			if err != nil {
				output.PrintFatalln(err.Error())
			}
			// unpack result.Output
			err = unpackFiles(finfo.FileBuf, finfo.PackName, targetdir)
			if err != nil {
				output.PrintFatalln(err.Error())
			}
			output.PrintInfoln("Success")
		}
	} else {
		// some error
		output.PrintJSON(msg)
	}
}

type ErrUnpack struct{}

func (ErrUnpack) Error() string { return "ErrUnpack" }

func unpackFiles(packb64 string, packName string, targetDir string) error {
	pname := filepath.Base(packName)
	extPos := -1
	for i := len(pname) - 1; i > 0; i-- {
		if pname[i] == '.' {
			extPos = i
			break
		}
	}
	if extPos == -1 {
		return fmt.Errorf("extPos == -1")
	}
	wname := pname[:extPos]
	tmpDir := ".octopoda_tmp/"
	os.Mkdir(tmpDir, os.ModePerm)
	defer os.RemoveAll(tmpDir)

	ppath := tmpDir + pname
	wpath := tmpDir + wname

	err := saveFile(packb64, ppath)
	if err != nil {
		return err
	}

	archiver.DefaultZip.OverwriteExisting = true
	err = archiver.DefaultZip.Unarchive(ppath, tmpDir)
	if err != nil {
		return ErrUnpack{}
	}

	os.MkdirAll(targetDir, os.ModePerm)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe", "/C", fmt.Sprintf("cp -Force %s %s", wpath, targetDir))
	} else {
		cmd = exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -r %s %s", wpath, targetDir))
	}
	err = cmd.Run()
	if err != nil {
		output.PrintWarningln("Wrap files: " + wpath + "-->" + targetDir + " | " + cmd.String())
		return fmt.Errorf("cp -r")
	}

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
		if path[0] != '/' {
			result.WriteString(base)
		}
		result.WriteString(path)
	}
	return result.String()
}
