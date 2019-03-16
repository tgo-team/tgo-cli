package client

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
)

func init()  {
	AddCMD(func(client *Client) CMD {

		return NewRegisterCMD(client)
	})
}

// 注册命令
type RegisterCMD struct {
	UID *uint64
	Password *string
	cmd *kingpin.CmdClause
	client *Client
}

func NewRegisterCMD(client *Client) *RegisterCMD  {
	return &RegisterCMD{
		client:client,
	}
}

func (r *RegisterCMD) Parse(app *kingpin.Application)  {
	r.cmd    = app.Command("register", "注册一个用户 uid必须为数字")
	r.UID = r.cmd.Arg("uid", "请输入一个唯一数字").Required().Uint64()
	r.Password = r.cmd.Arg("password", "请输入密码").Required().String()
}

func (r *RegisterCMD) Execute() error  {
	status,err := r.client.Register(*r.UID,*r.Password)
	if err!=nil {
		return err
	}
	if status == 2 {
		return fmt.Errorf("注册失败！-> 状态码[%d]",status)
	}
	if status == 3 {
		return fmt.Errorf("客户端已存在！")
	}
	return nil
}

func (r *RegisterCMD) Match(cmd string) bool  {
	if r.cmd.FullCommand() == cmd {
		return true
	}
	return false
}