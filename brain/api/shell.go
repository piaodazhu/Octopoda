package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/gin-gonic/gin"
)

type sshInfo struct {
	Addr     string
	Username string
	Password string
}

func SSHInfo(ctx *gin.Context) {
	var name string
	var ok bool
	rmsg := message.Result{
		Rmsg: "OK",
	}
	if name, ok = ctx.GetQuery("name"); !ok {
		rmsg.Rmsg = "Lack node name"
		ctx.JSON(404, rmsg)
		return
	}
	if name == "master" {
		sinfo := sshInfo{
			Addr:     fmt.Sprintf("%s:%d", config.GlobalConfig.Sshinfo.Ip, config.GlobalConfig.Sshinfo.Port),
			Username: config.GlobalConfig.Sshinfo.Username,
			Password: config.GlobalConfig.Sshinfo.Password,
		}
		ctx.JSON(200, sinfo)
		return
	}
	var conn *net.Conn
	if conn, ok = model.GetNodeMsgConn(name); !ok {
		rmsg.Rmsg = "Invalid node"
		ctx.JSON(404, rmsg)
		return
	}

	raw, err := message.Request(conn, message.TypeCommandSSH, []byte{})
	if err != nil {
		rmsg.Rmsg = "Node Error:" + err.Error()
		ctx.JSON(404, rmsg)
		return
	}
	ctx.Data(200, "application/json", raw)
}

type ScriptParams struct {
	FileName   string
	TargetPath string
	FileBuf    string
}

func RunScript(ctx *gin.Context) {
	script, _ := ctx.FormFile("script")
	targetNodes := ctx.PostForm("targetNodes")
	rmsg := message.Result{
		Rmsg: "OK",
	}

	if script.Size == 0 || len(targetNodes) == 0 {
		logger.Request.Println("RunScript Args Error")
		rmsg.Rmsg = "ERORR: arguments"
		ctx.JSON(400, rmsg)
		return
	}

	f, err := script.Open()
	if err != nil {
		rmsg.Rmsg = "Open:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}
	defer f.Close()

	raw, _ := io.ReadAll(f)
	content := base64.RawStdEncoding.EncodeToString(raw)
	sparams := ScriptParams{
		FileName:   script.Filename,
		TargetPath: "scripts/",
		FileBuf:    content,
	}
	payload, _ := config.Jsoner.Marshal(&sparams)

	nodes := []string{}
	err = config.Jsoner.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Rmsg = "ERROR: targetNodes"
		ctx.JSON(400, rmsg)
		return
	}

	results := make([]BasicNodeResults, len(nodes))
	var wg sync.WaitGroup

	for i := range nodes {
		name := nodes[i]
		results[i].Name = name
		// if addr, exists := model.GetNodeAddress(name); exists {
		// 	wg.Add(1)
		// 	go runScript(addr, payload, &wg, &results[i].Result)
		// } else {
		// 	results[i].Result = "NodeNotExists"
		// }
		wg.Add(1)
		go runScript(name, payload, &wg, &results[i].Result)
	}
	wg.Wait()
	ctx.JSON(200, results)
}

func RunCmd(ctx *gin.Context) {
	cmd := ctx.PostForm("command")
	targetNodes := ctx.PostForm("targetNodes")
	rmsg := message.Result{
		Rmsg: "OK",
	}

	if len(cmd) == 0 || len(targetNodes) == 0 {
		logger.Request.Println("RunCmd Args Error")
		rmsg.Rmsg = "ERORR: arguments"
		ctx.JSON(400, rmsg)
		return
	}
	nodes := []string{}
	err := config.Jsoner.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Rmsg = "Unmarshal:" + err.Error()
		ctx.JSON(400, rmsg)
		return
	}

	results := make([]BasicNodeResults, len(nodes))
	var wg sync.WaitGroup

	for i := range nodes {
		name := nodes[i]
		results[i].Name = name
		// if addr, exists := model.GetNodeAddress(name); exists {
		// 	wg.Add(1)
		// 	go runCmd(addr, cmd, &wg, &results[i].Result)
		// } else {
		// 	results[i].Result = "NodeNotExists"
		// }
		wg.Add(1)
		go runCmd(name, cmd, &wg, &results[i].Result)
	}
	wg.Wait()
	ctx.JSON(200, results)
}

func runCmd(name string, cmd string, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "UnknownError"

	var ok bool
	var conn *net.Conn
	if conn, ok = model.GetNodeMsgConn(name); !ok {
		logger.Comm.Println("GetNodeMsgConn")
		*result = "Connection not exists"
		return
	}
	
	raw, err := message.Request(conn, message.TypeCommandRun, []byte(cmd))
	if err != nil {
		logger.Comm.Println("Request", err)
		*result = "Request error"
		return
	}
	var rmsg message.Result
	err = config.Jsoner.Unmarshal(raw, &rmsg)
	if err != nil {
		logger.Exceptions.Println("UnmarshalNodeState", err)
		*result = "BrainError"
		return
	}
	// logger.SysInfo.Print(rmsg.Rmsg)
	// *result = rmsg.Output
	*result = fmt.Sprintf("[%s]: %s", rmsg.Rmsg, rmsg.Output)
}

func runScript(name string, payload []byte, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "UnknownError"

	var ok bool
	var conn *net.Conn
	if conn, ok = model.GetNodeMsgConn(name); !ok {
		logger.Comm.Println("GetNodeMsgConn")
		*result = "Connection not exists"
		return
	}
	raw, err := message.Request(conn, message.TypeCommandRunScript, payload)
	if err != nil {
		logger.Comm.Println("Request", err)
		*result = "Request error"
		return
	}
	var rmsg message.Result
	err = config.Jsoner.Unmarshal(raw, &rmsg)
	if err != nil {
		logger.Exceptions.Println("UnmarshalNodeState", err)
		*result = "Brain Error"
		return
	}
	// logger.SysInfo.Print(rmsg.Rmsg)
	*result = fmt.Sprintf("[%s]: %s", rmsg.Rmsg, rmsg.Output)
}
