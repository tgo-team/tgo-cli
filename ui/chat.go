package ui

import (
	"fmt"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"log"
)

type ChatUI struct {
	sessionList  *widgets.List
	chat *widgets.TextBox
	inputHeight int
	input *widgets.TextBox
	text string
}

func NewChatUI() *ChatUI {
	return &ChatUI{
		inputHeight: 4,
	}
}

func (c *ChatUI) Init() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	c.drawFrame()

	c.PollEvents()

}

func (c *ChatUI) drawFrame() {
	c.drawSessionList()

	c.drawChat()

	c.drawInput()
}

func (c *ChatUI) drawSessionList() {
	width, height := ui.TerminalDimensions()
	sessionListHeight := height
	sessionListWidth := width / 6

	c.sessionList = widgets.NewList()
	c.sessionList .Title = "最近会话"
	c.sessionList .Rows = []string{
		"[0] github.com/gizak/termui/v3",
		"[1] [你好，世界](fg:blue)",
		"[2] [こんにちは世界](fg:red)",
		"[3] [color](fg:white,bg:green) output",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] foo",
		"[8] bar",
		"[9] baz",
	}
	c.sessionList .TextStyle = ui.NewStyle(ui.ColorYellow)
	c.sessionList .WrapText = false
	c.sessionList .SetRect(0, 0, sessionListWidth, sessionListHeight)
	ui.Render(c.sessionList )
}

func (c *ChatUI) drawChat() {
	width, height := ui.TerminalDimensions()
	c.chat = widgets.NewTextBox()
	c.chat.InsertText("Borderless Text")
	c.chat.SetRect(c.sessionList.Max.X, 0, width, height-c.inputHeight)
	c.chat.Border = true
	ui.Render(c.chat )
}

func (c *ChatUI) drawInput()  {
	width, _ := ui.TerminalDimensions()
	c.input = widgets.NewTextBox()
	c.input.InsertText("")
	c.input.SetRect(c.sessionList.Max.X, c.sessionList.Max.Y, width, c.sessionList.Max.Y-c.inputHeight)
	c.input.Border = true
	c.input.ShowCursor =true
	ui.Render(c.input)
}

func (c *ChatUI) Render() {

}

func (c *ChatUI) refreshInputText()  {
	c.input.SetText(c.text)
}

func (c *ChatUI) PollEvents() {
	defer ui.Close()
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			if e.Type == ui.ResizeEvent {
				ui.Clear()
				c.drawFrame()
			}
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Left>":
				c.input.MoveCursorLeft()
			case "<Right>":
				c.input.MoveCursorRight()
			case "<Backspace>":
				c.input.Backspace()
				runslice := []rune(c.text)
				c.text = string(runslice[0:len(runslice)-1])
				c.refreshInputText()
			case "<Enter>":
				//c.input.InsertText("\n")
			case "<Tab>":
				//c.input.InsertText("\t")
			case "<Space>":
				c.text = fmt.Sprintf("%s%s",c.text," ")
				runslice := []rune(c.text)
				c.text = string(runslice)
				c.refreshInputText()
			default:
				c.text = fmt.Sprintf("%s%s",c.text,e.ID)
				runslice := []rune(c.text)
				c.text = string(runslice)

				c.refreshInputText()
			}
			ui.Render(c.input)
		}
	}
}
