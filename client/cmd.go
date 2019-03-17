package client

import (
	"fmt"
	"github.com/tgo-team/tgo-talk/tgo/packets"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

type CMD interface {
	Parse(app *kingpin.Application)
	Execute() error
	Match(cmd string) bool
}

var cmdList = make([]CMD,0)

func SetupCMD(client *Client)  {
	cmdList = getAllCmd(client)
	for _, cmd := range cmdList {
		cmd.Parse(app)
	}
	cmdStr := kingpin.MustParse(app.Parse(os.Args[1:]))
	Execute(cmdStr)
}
func getAllCmd(client *Client) []CMD {
	cmds := make([]CMD,0,len(cmdFuncs))
	for _,cmdFunc :=range cmdFuncs {
		cmds = append(cmds,cmdFunc(client))
	}
	return cmds
}
var (
	app = kingpin.New("tc", "this is tgo talk")
)


type cmdFunc func(client *Client) CMD
var cmdFuncs  =make([]cmdFunc,0)

func AddCMD(cmdFunc cmdFunc)  {
	cmdFuncs = append(cmdFuncs,cmdFunc)
}

func Execute(cmdStr string)  {
	for _, cmd := range cmdList {
		if cmd.Match(cmdStr) {
			err := cmd.Execute()
			if err!=nil {
				fmt.Println(err.Error())
			}
		}
	}
}



// SetCurrentChannel 设置当前管道
var currentID uint64
func SetCurrentChannel(channelID uint64)  {
	currentID = channelID
}
// GetCurrentChannel 获取当前管道
func GetCurrentChannel() uint64  {
	return currentID
}

// showSendState 显示发送状态
func showSendState() {
	channelID := GetCurrentChannel()
	if channelID > 0 {
		fmt.Print(fmt.Sprintf("%d>", channelID))
	} else {
		fmt.Print(">")
	}
}

// showRevMsg 显示接受的消息
func showRevMsg(msgPacket *packets.MessagePacket,) {
	fmt.Println(fmt.Sprintf("收到[%d]的消息->%s",msgPacket.From,string(msgPacket.Payload)))
	fmt.Print(">")
}