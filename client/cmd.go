package client

import (
	"fmt"
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