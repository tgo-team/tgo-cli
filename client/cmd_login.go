package client

import (
	"bufio"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
)

func init() {
	AddCMD(func(client *Client) CMD {
		return NewLoginCMD(client)
	})
}

// 注册命令
type LoginCMD struct {
	UID      *uint64
	Password *string
	cmd      *kingpin.CmdClause
	client   *Client

	test *kingpin.CmdClause
	app  *kingpin.Application
}

func NewLoginCMD(client *Client) *LoginCMD {
	return &LoginCMD{
		client: client,
	}
}

func (r *LoginCMD) Parse(app *kingpin.Application) {
	r.app = app
	r.cmd = app.Command("login", "登录")
	r.UID = r.cmd.Arg("uid", "请输入一个唯一数字").Required().Uint64()
	r.Password = r.cmd.Arg("password", "请输入密码").Required().String()
}

func (r *LoginCMD) Execute() error {
	err := r.client.Login(*r.UID, *r.Password)
	if err != nil {
		return err
	}
	fmt.Println("登录成功！")

	input := bufio.NewScanner(os.Stdin)
	for {
		channelID := GetCurrentChannel()
		showSendState()
		input.Scan()
		cmd := input.Text()
		if cmd == "" || strings.TrimSpace(cmd) == "" {
			continue
		}
		var cmdStr string
		if channelID > 0 {
			cmd = fmt.Sprintf("send %s -c %d", cmd, GetCurrentChannel())
		}
		cmdStr, err := r.app.Parse(r.splitCMD(cmd))
		if err != nil {
			fmt.Println(err)
			continue
		}
		Execute(cmdStr)
	}
	return nil
}



func (r *LoginCMD) Match(cmd string) bool {
	if r.cmd.FullCommand() == cmd {
		return true
	}
	return false
}

func (r *LoginCMD) splitCMD(cmd string) []string {
	cmds := strings.Split(cmd, " ")
	newCmds := make([]string, 0, len(cmds))
	for _, cmd := range cmds {
		noSpaceCmd := strings.TrimSpace(cmd)
		if noSpaceCmd == "" {
			continue
		}
		newCmds = append(newCmds, noSpaceCmd)
	}
	return newCmds
}
