package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"encoding/base64"
	"encoding/json"
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
	var name, addr string
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
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
	if addr, ok = model.GetNodeAddress(name); !ok {
		ctx.JSON(404, struct{}{})
		return
	}

	if conn, err := net.Dial("tcp", addr); err != nil {
		ctx.JSON(404, struct{}{})
	} else {
		defer conn.Close()
		message.SendMessage(conn, message.TypeCommandSSH, []byte{})
		mtype, payload, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeCommandResponse {
			logger.Tentacle.Println("NodeReboot", err)
			ctx.JSON(404, struct{}{})
			return
		}
		ctx.Data(200, "application/json", payload)
	}
}

type ScriptParams struct {
	FileName   string
	TargetPath string
	FileBuf    string
}

func RunScript(ctx *gin.Context) {
	script, _ := ctx.FormFile("script")
	targetNodes := ctx.PostForm("targetNodes")
	rmsg := RMSG{}

	if script.Size == 0 || len(targetNodes) == 0 {
		logger.Brain.Println("RunScript Args Error")
		rmsg.Msg = "ERORR: arguments"
		ctx.JSON(400, rmsg)
		return
	}

	f, err := script.Open()
	if err != nil {
		rmsg.Msg = "ERORR: Open Script"
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
	payload, _ := json.Marshal(&sparams)

	nodes := []string{}
	err = json.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Msg = "ERROR: targetNodes"
		ctx.JSON(400, rmsg)
		return
	}

	results := make([]UploadResults, len(nodes))
	var wg sync.WaitGroup

	for i := range nodes {
		name := nodes[i]
		results[i].Name = name
		if addr, exists := model.GetNodeAddress(name); exists {
			wg.Add(1)
			go runScript(addr, payload, &wg, &results[i].Result)
		} else {
			results[i].Result = "NodeNotExists"
		}
	}
	wg.Wait()
	ctx.JSON(200, results)
}

func RunCmd(ctx *gin.Context) {
	cmd := ctx.PostForm("command")
	targetNodes := ctx.PostForm("targetNodes")
	rmsg := RMSG{}

	if len(cmd) == 0 || len(targetNodes) == 0 {
		logger.Brain.Println("RunCmd Args Error")
		rmsg.Msg = "ERORR: arguments"
		ctx.JSON(400, rmsg)
		return
	}
	nodes := []string{}
	err := json.Unmarshal([]byte(targetNodes), &nodes)
	if err != nil {
		rmsg.Msg = "ERROR: targetNodes"
		ctx.JSON(400, rmsg)
		return
	}

	results := make([]UploadResults, len(nodes))
	var wg sync.WaitGroup

	for i := range nodes {
		name := nodes[i]
		results[i].Name = name
		if addr, exists := model.GetNodeAddress(name); exists {
			wg.Add(1)
			go runCmd(addr, cmd, &wg, &results[i].Result)
		} else {
			results[i].Result = "NodeNotExists"
		}
	}
	wg.Wait()
	ctx.JSON(200, results)
}

func runCmd(addr string, cmd string, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "OK"

	if conn, err := net.Dial("tcp", addr); err != nil {
		return
	} else {
		defer conn.Close()

		message.SendMessage(conn, message.TypeCommandRun, []byte(cmd))
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeCommandResponse {
			logger.Tentacle.Println("TypeCommandResponse", err)
			*result = "NetError"
			return
		}

		var rmsg RMSG
		err = json.Unmarshal(raw, &rmsg)
		if err != nil {
			logger.Tentacle.Println("UnmarshalNodeState", err)
			*result = "MasterError"
			return
		}
		logger.Tentacle.Print(rmsg.Msg)
		*result = rmsg.Msg
	}
}

func runScript(addr string, payload []byte, wg *sync.WaitGroup, result *string) {
	defer wg.Done()
	*result = "OK"

	if conn, err := net.Dial("tcp", addr); err != nil {
		return
	} else {
		defer conn.Close()

		message.SendMessage(conn, message.TypeCommandRunScript, payload)
		mtype, raw, err := message.RecvMessage(conn)
		if err != nil || mtype != message.TypeCommandResponse {
			logger.Tentacle.Println("TypeCommandResponse", err)
			*result = "NetError"
			return
		}

		var rmsg RMSG
		err = json.Unmarshal(raw, &rmsg)
		if err != nil {
			logger.Tentacle.Println("UnmarshalNodeState", err)
			*result = "MasterError"
			return
		}
		logger.Tentacle.Print(rmsg.Msg)
		*result = rmsg.Msg
	}
}