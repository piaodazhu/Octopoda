package shell

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"octl/config"
	"octl/nameclient"
	"octl/output"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/piaodazhu/proxylite"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type SSHInfo struct {
	Addr     string
	Username string
	Password string
}

type proxyMsg struct {
	Code int
	Msg  string
	Data string
}

func SetSSH(nodename string) {
	// set username and password
	form := nameclient.SshInfoUploadParam{Type: "other", Name: nodename}

	fmt.Println("Please enter its username: ")
	fmt.Scanln(&form.Username)
	if form.Username == "" {
		output.PrintFatalln("username must not leave empty")
	}

	fmt.Println("Please enter its password: ")
	pass, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		output.PrintFatalln("ReadPassword error:", err)
	}
	form.Password = string(pass)

	var confirm string
	fmt.Println("Please enter [yes|no] to confirm: ")
	fmt.Scanln(&confirm)
	if confirm != "yes" && confirm != "y" {
		output.PrintInfoln("you cancelled setssh. Bye")
		os.Exit(0)
	}

	if _, err := nameclient.SshinfoQuery(nodename); err == nil {
		delSSH(nodename)
	}

	// ask tentacle to register ssh service
	URL := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.Ssh,
	)
	values := url.Values{}
	values.Add("name", nodename)
	res, err := http.PostForm(URL, values)
	if err != nil {
		output.PrintFatalln("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	// get mapped ip:port from response
	pmsg := proxyMsg{}
	if err := json.Unmarshal(raw, &pmsg); err != nil {
		output.PrintFatalln("Unmarshal proxyMsg: ", err)
	}
	if pmsg.Code != 0 {
		output.PrintFatalln("Register ssh proxy failed: ", pmsg)
	}
	ss := strings.Split(pmsg.Data, ":")
	if len(ss) != 2 {
		output.PrintFatalln("Parse ssh proxy response failed: ", pmsg)
	}
	form.Ip = ss[0]

	var port int
	if port, err = strconv.Atoi(ss[1]); err != nil {
		output.PrintFatalln("Register ssh proxy failed: ", pmsg)
	}
	form.Port = port

	// register sshinfo to nameserver
	err = nameclient.SshinfoRegister(&form)
	if err != nil {
		output.PrintFatalln("SshinfoRegister error:", err)
	}
	output.PrintInfoln("SshinfoRegister success")
}

func delSSH(nodename string) []byte {
	defer nameclient.NameDelete(nodename, "ssh")
	URL := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.Ssh,
		nodename,
	)

	req, err := http.NewRequest("DELETE", URL, nil)
	if err != nil {
		output.PrintFatalln("NewRequest: ", err)
		return nil
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		output.PrintFatalln("DELETE")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	return raw
}

func DelSSH(nodename string) {
	output.PrintJSON(delSSH(nodename))
}

func GetSSH() {
	entry, err := nameclient.NameQuery(config.GlobalConfig.Brain.Name + ".proxyliteFace")
	if err != nil {
		output.PrintFatalln("cannot query proxyliteFace: ", err)
	}

	infos, err := proxylite.DiscoverServices(fmt.Sprintf("%s:%d", entry.Ip, entry.Port))
	if err != nil {
		output.PrintFatalln("cannot discover services: ", err)
	}

	raw, _ := json.Marshal(infos)
	output.PrintJSON(raw)
}

func SSH(nodename string) {
	sshinfo, err := nameclient.SshinfoQuery(nodename)
	if err != nil {
		output.PrintFatalln("SshinfoQuery error:", err)
	}

	dossh(sshinfo.Username, sshinfo.Ip, sshinfo.Password, sshinfo.Port)
}

func dossh(user, ip, passwd string, port int) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(passwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatalf("SSH dial error: %s", err.Error())
	}
	if runtime.GOOS == "windows" {
		dossh_windows(user, ip, port)
		return
	}

	// 建立新会话
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("new session error: %s", err.Error())
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	modes := ssh.TerminalModes{
		ssh.ECHO: 1, // 禁用回显（0禁用，1启动）
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	fileDescriptor := int(os.Stdin.Fd())
	if term.IsTerminal(fileDescriptor) {
		originalState, err := term.MakeRaw(fileDescriptor)
		if err != nil {
			return
		}
		defer term.Restore(fileDescriptor, originalState)

		err = session.RequestPty("xterm-256color", 32, 160, modes)
		if err != nil {
			return
		}
	}

	if err = session.Shell(); err != nil {
		log.Fatalf("start shell error: %s", err.Error())
	}
	if err = session.Wait(); err != nil {
		log.Fatalf("return error: %s", err.Error())
	}
}

func dossh_windows(user, ip string, port int) {
	output.PrintWarningln("Windows user have to input ssh password even if the password has been set before.")
	cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", user, ip), "-p", fmt.Sprint(port))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}
