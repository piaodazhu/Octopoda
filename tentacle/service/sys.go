package service

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
	"github.com/piaodazhu/Octopoda/tentacle/task"
)

type ScriptParams struct {
	FileName   string
	TargetPath string
	FileBuf    string
	DelayTime  int
}

var shellPath string

func init() {
	shellPath = "/bin/bash"
	_, err := os.Stat(shellPath)
	if err == nil {
		return
	}
	shellPath = "/bin/sh"
	_, err = os.Stat(shellPath)
	if err == nil {
		return
	}
	shellPath = "sh"
}

func execScript(sparams *ScriptParams, dir string, cmdChan chan *exec.Cmd) ([]byte, error) {
	var content []byte
	var err error
	var scriptFile strings.Builder
	// var f *os.File
	content, err = base64.RawStdEncoding.DecodeString(sparams.FileBuf)
	if err != nil {
		logger.Exceptions.Println("FileDecode")
		return nil, err
	}

	scriptFile.WriteString(config.GlobalConfig.Workspace.Root)
	scriptFile.WriteString(sparams.TargetPath)

	scriptDir := scriptFile.String()
	os.Mkdir(scriptDir, os.ModePerm)
	if scriptFile.String()[scriptFile.Len()-1] != '/' {
		scriptFile.WriteByte('/')
	}

	scriptFile.WriteString(sparams.FileName)
	err = os.WriteFile(scriptFile.String(), content, os.ModePerm)
	if err != nil {
		logger.Exceptions.Println("WriteFile")
		return nil, err
	}
	outputFile := scriptFile.String() + ".output"
	output, _ := os.Create(outputFile)

	// fbuf, _ := os.ReadFile(scriptFile.String())
	// logger.Client.Println(string(fbuf))

	cmd := exec.Command(shellPath, scriptFile.String())
	cmd.Dir = dir
	cmd.Env = append(syscall.Environ(), config.OctopodaEnv(scriptDir, sparams.FileName, outputFile)...)

	if cmdChan != nil {
		cmdChan <- cmd
		close(cmdChan)
	}

	scriptErr := cmd.Run()
	if scriptErr != nil {
		logger.Exceptions.Println("Run cmd", err)
	}

	// read output
	// result := []byte{}
	result, err := io.ReadAll(output)
	if err != nil {
		logger.Exceptions.Println(err)
	}
	output.Close()

	// logger.Client.Println("exit code:", cmd.ProcessState.ExitCode())

	os.Remove(scriptFile.String())
	os.Remove(outputFile)

	return result, scriptErr
}

type CommandParams struct {
	Command    string
	Background bool
	DelayTime  int
}

func RunCmd(conn net.Conn, serialNum uint32, raw []byte) {
	cparams := CommandParams{}

	if err := config.Jsoner.Unmarshal(raw, &cparams); err != nil {
		logger.Exceptions.Println("invalid arguments: ", err)
		// SNED BACK
		err = protocols.SendMessageUnique(conn, protocols.TypeRunCommandResponse, serialNum, []byte{})
		if err != nil {
			logger.Comm.Println("TypeRunCommandResponse send error")
		}
		return
	}

	var utaskFunc func() *protocols.Result
	var ucancelFunc func()
	var cmd *exec.Cmd
	var delayTimer *time.Timer = nil

	if cparams.Background {
		cmd = exec.Command(shellPath, "-c", cparams.Command)
		cmd.Dir = config.GlobalConfig.Workspace.Root
		cmd.Env = append(syscall.Environ(), config.OctopodaEnv(config.GlobalConfig.Workspace.Root, "NONE", "NONE")...)
	} else {
		cmd = exec.Command(shellPath, "-c", cparams.Command)
		cmd.Dir = config.GlobalConfig.Workspace.Root
		cmd.Env = append(syscall.Environ(), config.OctopodaEnv(config.GlobalConfig.Workspace.Root, "NONE", "NONE")...)
	}

	utaskFunc = func() *protocols.Result {
		rmsg := protocols.Result{
			Rmsg: "OK",
		}
		runFunc := func() {
			if cparams.Background {
				startErr := cmd.Start()
				if startErr != nil {
					emsg := fmt.Sprintf("run cmd background error when start the process: %s. command is: %s", startErr.Error(), cmd.String())
					logger.Exceptions.Println(emsg)
					rmsg.Rmsg = emsg
					rmsg.Output = "0"
				}
				rmsg.Output = fmt.Sprint(cmd.Process.Pid)
				go func() {
					execErr := cmd.Wait()
					if execErr != nil {
						emsg := fmt.Sprintf("run cmd background error when wait the process: %s. command is: %s", execErr.Error(), cmd.String())
						logger.Exceptions.Println(emsg)
						rmsg.Rmsg = emsg
						rmsg.Output = "0"
					}
				} ()
			} else {
				stdoutPipe, err := cmd.StdoutPipe()
				if err != nil {
					emsg := fmt.Sprintf("run cmd foreground error when open stdout pipe: %s. command is: %s", err.Error(), cmd.String())
					logger.Exceptions.Println(emsg)
					rmsg.Rmsg = emsg
					return
				}
				defer stdoutPipe.Close()

				stderrPipe, err := cmd.StderrPipe()
				if err != nil {
					emsg := fmt.Sprintf("run cmd foreground error when open stderr pipe: %s. command is: %s", err.Error(), cmd.String())
					logger.Exceptions.Println(emsg)
					rmsg.Rmsg = emsg
					return
				}
				defer stderrPipe.Close()

				var stdoutSb, stderrSb strings.Builder

				wg := sync.WaitGroup{}
				wg.Add(2)
				go func() {
					defer wg.Done()
					scanner := bufio.NewScanner(stdoutPipe)
					for scanner.Scan() {
						stdoutSb.WriteString(scanner.Text())
						stdoutSb.WriteByte('\n')
					}
				}()
				go func() {
					defer wg.Done()
					scanner := bufio.NewScanner(stderrPipe)
					for scanner.Scan() {
						stderrSb.WriteString(scanner.Text())
						stderrSb.WriteByte('\n')
					}
				}()

				if err := cmd.Start(); err != nil {
					emsg := fmt.Sprintf("run cmd foreground error when start the process: %s. command is: %s", err.Error(), cmd.String())
					logger.Exceptions.Println(emsg)
					rmsg.Rmsg = emsg
					return
				}

				wg.Wait()
				if err := cmd.Wait(); err != nil {
					emsg := fmt.Sprintf("run cmd foreground error when wait the process: %s. command is: %s", err.Error(), cmd.String())
					logger.Exceptions.Println(emsg)
					rmsg.Rmsg = emsg
					rmsg.Output = stderrSb.String()
					return
				}
				rmsg.Output = stdoutSb.String()
			}
		}

		if cparams.DelayTime < 0 {
			runFunc()
		} else {
			delayTimer = time.AfterFunc(time.Duration(cparams.DelayTime)*time.Second, runFunc)
			rmsg.Output = "loaded"
		}
		return &rmsg
	}

	ucancelFunc = func() {
		cmd.Process.Kill()
		if delayTimer != nil {
			delayTimer.Stop()
		}
	}

	taskId, err := task.TaskManager.CreateTask(cmd.String(), utaskFunc, ucancelFunc)
	if err != nil {
		// ERROR
		logger.Exceptions.Println("cannot create task: ", err)
		err = protocols.SendMessageUnique(conn, protocols.TypeRunCommandResponse, serialNum, []byte{})
		if err != nil {
			logger.Comm.Println("TypeRunCommandResponse send error")
		}
		return
	}
	err = protocols.SendMessageUnique(conn, protocols.TypeRunCommandResponse, serialNum, []byte(taskId))
	if err != nil {
		logger.Comm.Println("TypeRunCommandResponse send error")
	}
}

func RunScript(conn net.Conn, serialNum uint32, raw []byte) {
	sparams := ScriptParams{}
	if err := config.Jsoner.Unmarshal(raw, &sparams); err != nil {
		logger.Exceptions.Println(err)
		err = protocols.SendMessageUnique(conn, protocols.TypeRunScriptResponse, serialNum, []byte{})
		if err != nil {
			logger.Comm.Println("TypeRunScriptResponse send error")
		}
		return
	}

	var utaskFunc func() *protocols.Result
	var ucancelFunc func()
	cmdChan := make(chan *exec.Cmd, 1)
	var delayTimer *time.Timer = nil

	utaskFunc = func() *protocols.Result {
		rmsg := protocols.Result{
			Rmsg: "OK",
		}
		runFunc := func() {
			output, err := execScript(&sparams, config.GlobalConfig.Workspace.Root, cmdChan)
			if err != nil {
				rmsg.Rmsg = err.Error()
			} else {
				rmsg.Output = string(output)
			}
		}
		if sparams.DelayTime < 0 {
			runFunc()
		} else {
			delayTimer = time.AfterFunc(time.Duration(sparams.DelayTime)*time.Second, runFunc)
			rmsg.Output = "loaded"
		}
		return &rmsg
	}

	ucancelFunc = func() {
		cmd := <-cmdChan
		cmd.Process.Kill()
		if delayTimer != nil {
			delayTimer.Stop()
		}
	}

	taskId, err := task.TaskManager.CreateTask(sparams.FileName, utaskFunc, ucancelFunc)
	if err != nil {
		// ERROR
		logger.Exceptions.Println("cannot create task: ", err)
		err = protocols.SendMessageUnique(conn, protocols.TypeRunScriptResponse, serialNum, []byte{})
		if err != nil {
			logger.Comm.Println("TypeRunScriptResponse send error")
		}
		return
	}
	err = protocols.SendMessageUnique(conn, protocols.TypeRunScriptResponse, serialNum, []byte(taskId))
	if err != nil {
		logger.Comm.Println("TypeRunCommandResponse send error")
	}
}
