package file

import (
	"encoding/base64"
	"fmt"
	"github.com/piaodazhu/Octopoda/octl/output"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mholt/archiver/v3"
)

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
		return fmt.Errorf("error when unpacking directory %s", packName)
	}

	os.MkdirAll(targetDir, os.ModePerm)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe", "/C", fmt.Sprintf("cp -Force -r %s/* %s", wpath, targetDir))
	} else {
		cmd = exec.Command("/bin/bash", "-c", fmt.Sprintf("cp -r %s/* %s", wpath, targetDir))
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
		if !filepath.IsAbs(path) {
			result.WriteString(base)
		}
		result.WriteString(path)
	}
	return result.String()
}
