package client

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

func init()  {
	AddCMD(func(client *Client) CMD {
		return NewChannelCMD(client)
	})
}

type ChannelCMD struct {
	cmd     *kingpin.CmdClause
	channel *uint64
	client  *Client
}

func NewChannelCMD(client *Client) *ChannelCMD {
	return &ChannelCMD{
		client: client,
	}
}

func (c *ChannelCMD) Parse(app *kingpin.Application) {
	c.cmd = app.Command("channel", "进入聊天管道")
	c.channel = c.cmd.Arg("channelID","管道ID").Required().Uint64()
}

func (c *ChannelCMD) Execute() error  {
	SetCurrentChannel(*c.channel)
	return nil
}

func (c *ChannelCMD) Match(cmd string) bool  {
	if c.cmd.FullCommand() == cmd {
		return true
	}
	return false
}
