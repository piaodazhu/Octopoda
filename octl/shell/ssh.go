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
	"strconv"
	"strings"
	"time"

	"github.com/piaodazhu/proxylite"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type SSHInfo struct {
	Addr     string
	Username string
	Password string
}

type SSHTerminal struct {
	Session *ssh.Session
	exitMsg string
	stdout  io.Reader
	stdin   io.Writer
	stderr  io.Reader
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

func DelSSH(nodename string) {
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
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		output.PrintFatalln("DELETE")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	output.PrintJSON(raw)
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
	if nodename == "master" {
		sshMaster()
		return
	}
	sshinfo, err := nameclient.SshinfoQuery(nodename)
	if err != nil {
		output.PrintFatalln("SshinfoQuery error:", err)
	}
	addr := fmt.Sprintf("%s:%d", sshinfo.Ip, sshinfo.Port)
	dossh(addr, sshinfo.Username, sshinfo.Password)
}

func sshMaster() {
	var username, passwd string
	fmt.Println("Please enter its username: ")
	fmt.Scanln(&username)

	fmt.Println("Please enter its password: ")
	pass, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		output.PrintFatalln("ReadPassword error:", err)
	}
	passwd = string(pass)
	dossh(config.GlobalConfig.Brain.Ip, username, passwd)
}

func dossh(addr, user, passwd string) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Dial your ssh server.
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("unable to connect: ", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.

	err = New(conn)
	if err != nil {
		fmt.Println(err)
	}
}

func (t *SSHTerminal) updateTerminalSize() {
	go func() {
		// SIGWINCH is sent to the process when the window size of the terminal has
		// changed.
		// sigwinchCh := make(chan os.Signal, 1)
		// signal.Notify(sigwinchCh, syscall.SIGWINCH)

		fd := int(os.Stdin.Fd())
		termWidth, termHeight, err := term.GetSize(fd)
		if err != nil {
			fmt.Println(err)
		}

		// for sigwinch := range sigwinchCh {
		for {
			// if sigwinch == nil {
			// 	return
			// }
			time.Sleep(time.Microsecond * 200)
			currTermWidth, currTermHeight, err := term.GetSize(fd)

			// Terminal size has not changed, don't do anything.
			if currTermHeight == termHeight && currTermWidth == termWidth {
				continue
			}

			t.Session.WindowChange(currTermHeight, currTermWidth)
			if err != nil {
				fmt.Printf("Unable to send window-change reqest: %s.", err)
				continue
			}

			termWidth, termHeight = currTermWidth, currTermHeight
		}
	}()

}

func (t *SSHTerminal) interactiveSession() error {

	defer func() {
		if t.exitMsg == "" {
			fmt.Fprintln(os.Stdout, "the connection was closed on the remote side on ", time.Now().Format(time.RFC822))
		} else {
			fmt.Fprintln(os.Stdout, t.exitMsg)
		}
	}()

	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, state)

	termWidth, termHeight, err := term.GetSize(fd)
	if err != nil {
		return err
	}

	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}

	err = t.Session.RequestPty(termType, termHeight, termWidth, ssh.TerminalModes{})
	if err != nil {
		return err
	}

	t.updateTerminalSize()

	t.stdin, err = t.Session.StdinPipe()
	if err != nil {
		return err
	}
	t.stdout, err = t.Session.StdoutPipe()
	if err != nil {
		return err
	}
	t.stderr, _ = t.Session.StderrPipe()

	go io.Copy(os.Stderr, t.stderr)
	go io.Copy(os.Stdout, t.stdout)
	go func() {
		buf := make([]byte, 128)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			if n > 0 {
				_, err = t.stdin.Write(buf[:n])
				if err != nil {
					fmt.Println(err)
					t.exitMsg = err.Error()
					return
				}
			}
		}
	}()

	err = t.Session.Shell()
	if err != nil {
		return err
	}
	err = t.Session.Wait()
	if err != nil {
		return err
	}
	return nil
}

func New(client *ssh.Client) error {

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	s := SSHTerminal{
		Session: session,
	}

	return s.interactiveSession()
}
