package client

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	AddCMD(func(client *Client) CMD {
		return NewMsgCMD(client)
	})
}

type MsgCMD struct {
	cmd *kingpin.CmdClause
	message *string
	client  *Client
	channelID *uint64
}

func NewMsgCMD(client *Client) *MsgCMD {

	return &MsgCMD{client: client}
}

func (m *MsgCMD) Parse(app *kingpin.Application) {
	m.cmd = app.Command("send", "发送一个消息")
	m.message = m.cmd.Arg("message", "消息内容").Required().String()
	m.channelID = m.cmd.Flag("cid","管道ID，如果是发给个人管道ID则为个人的uid。").Short('c').Required().Uint64()

}

func (m *MsgCMD) Execute() error {
	return m.client.SendMsg(*m.channelID,[]byte(*m.message))
}

func (m *MsgCMD) Match(cmd string) bool {
	if m.cmd.FullCommand() == cmd {
		return true
	}
	return false
}
