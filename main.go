package main

import "github.com/tgo-team/tgo-cli/client"

func main() {
	c := client.New(client.NewOptions())
	client.SetupCMD(c)

	//chat := ui.NewChatUI()
	//
	//chat.Init()

}
