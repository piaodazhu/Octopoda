package shell

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/workgroup"
	"github.com/piaodazhu/Octopoda/protocols"

	"github.com/piaodazhu/proxylite"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func SetSSH(nodename string) {
	// set username and password
	var username, password string
	fmt.Println("Please enter its username: ")
	fmt.Scanln(&username)
	if username == "" {
		output.PrintFatalln("username must not leave empty")
	}

	fmt.Println("Please enter its password: ")
	pass, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		output.PrintFatalln("ReadPassword error:", err)
	}
	password = string(pass)

	var confirm string
	fmt.Println("Please enter [yes|no] to confirm: ")
	fmt.Scanln(&confirm)
	if confirm != "yes" && confirm != "y" {
		output.PrintInfoln("you cancelled setssh. Bye")
		os.Exit(0)
	}

	// ask tentacle to register ssh service
	URL := fmt.Sprintf("https://%s/%s%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Ssh,
	)
	values := url.Values{}
	values.Add("name", nodename)
	values.Add("username", username)
	values.Add("password", password)

	req, _ := http.NewRequest("POST", URL, bytes.NewBufferString(values.Encode()))
	workgroup.SetHeader(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := httpclient.BrainClient.Do(req)
	if err != nil {
		output.PrintFatalln("PostForm")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	// get mapped ip:port from response
	pmsg := protocols.ProxyMsg{}
	if err := json.Unmarshal(raw, &pmsg); err != nil {
		output.PrintFatalln("Unmarshal proxyMsg: ", err)
	}
	if pmsg.Code != 0 {
		output.PrintInfoln(pmsg)
		return
	}
	output.PrintInfoln("SshinfoRegister success")
}

func delSSH(nodename string) []byte {
	URL := fmt.Sprintf("https://%s/%s%s?name=%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Ssh,
		nodename,
	)

	req, err := http.NewRequest("DELETE", URL, nil)
	if err != nil {
		output.PrintFatalln("NewRequest: ", err)
		return nil
	}
	workgroup.SetHeader(req)
	res, err := httpclient.BrainClient.Do(req)
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
	entry, err := httpclient.NameQuery(config.GlobalConfig.Brain.Name + ".proxyliteFace")
	if err != nil {
		output.PrintFatalln("cannot query proxyliteFace: ", err)
	}

	infos, err := proxylite.DiscoverServices(entry.Value)
	if err != nil {
		output.PrintFatalln("cannot discover services: ", err)
	}

	raw, _ := json.Marshal(infos)
	output.PrintJSON(raw)
}

func SSH(nodename string) {
	url := fmt.Sprintf("https://%s/%s%s?name=%s",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Ssh,
		nodename,
	)
	req, _ := http.NewRequest("GET", url, nil)
	workgroup.SetHeader(req)
	res, err := httpclient.BrainClient.Do(req)
	if err != nil {
		output.PrintFatalln("Get")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		output.PrintFatalln("ssh info of this node not found:", nodename)
	}
	raw, _ := io.ReadAll(res.Body)
	info := protocols.SSHInfo{}
	if err = json.Unmarshal(raw, &info); err != nil {
		output.PrintFatalln("Unmarshal:", err)
	}

	dossh(info.Username, info.Ip, info.Password, int(info.Port))
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
		ssh.ECHO:          1,     // 禁用回显（0禁用，1启动）
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	fileDescriptor := int(os.Stdin.Fd())
	originalState, err := term.MakeRaw(fileDescriptor)
	if err != nil {
		return
	}

	defer term.Restore(fileDescriptor, originalState)
	err = session.RequestPty("xterm-256color", 32, 160, modes)
	if err != nil {
		return
	}

	if err = session.Shell(); err != nil {
		log.Fatalf("start shell error: %s", err.Error())
	}
	session.Wait()
}

func dossh_windows(user, ip string, port int) {
	output.PrintWarningln("Windows user have to input ssh password even if the password has been set before.")
	cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", user, ip), "-p", fmt.Sprint(port))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}
