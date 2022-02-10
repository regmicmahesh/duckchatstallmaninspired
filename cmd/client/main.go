package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var wg sync.WaitGroup

func readMessage(conn net.Conn, message chan string) (string, error) {
	for {
		buffer := make([]byte, 1024)
		length, err := conn.Read(buffer)
		if err != nil {
			return "", err
		}
		message <- (strings.TrimSpace(string(buffer[:length])))
	}
}

func writeMessage(conn net.Conn, message string) error {
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing to server:", err.Error())
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: client <server-address> <username>")
		os.Exit(1)
	}
	serverAddress := os.Args[1]
	username := os.Args[2]

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()
	Mx, My := ui.TerminalDimensions()

	conn, _ := net.Dial("tcp", serverAddress)
	defer conn.Close()

	receivedMessageList := widgets.NewList()
	receivedMessageList.Title = "List"
	receivedMessageList.Rows = []string{
		"[+] Connected to server [+]",
	}
	receivedMessageList.TextStyle = ui.NewStyle(ui.ColorYellow)
	receivedMessageList.WrapText = false
	receivedMessageList.SetRect(0, 0, Mx, My-3)
	ui.Render(receivedMessageList)

	sendMessageBox := widgets.NewParagraph()
	sendMessageBox.Text = ""
	sendMessageBox.SetRect(0, My-3, Mx, My)
	ui.Render(sendMessageBox)

	msg := make(chan string)
	go readMessage(conn, msg)

	uiEvents := ui.PollEvents()

	for {
		select {
		case recievedMsg := <-msg:
			receivedMessageList.Rows = append(receivedMessageList.Rows, recievedMsg)
			receivedMessageList.ScrollBottom()
			ui.Render(receivedMessageList)

		case e := <-uiEvents:
			switch e.Type {
			case ui.KeyboardEvent:
				switch {
				case e.ID == "<C-c>":
					return
				case len(e.ID) == 1:
					sendMessageBox.Text += e.ID
					ui.Render(sendMessageBox)
				case e.ID == "<Backspace>" && len(sendMessageBox.Text) > 0:
					sendMessageBox.Text = sendMessageBox.Text[:len(sendMessageBox.Text)-1]
					ui.Render(sendMessageBox)
				case e.ID == "<Space>":
					sendMessageBox.Text += " "
					ui.Render(sendMessageBox)
				case e.ID == "<Enter>" && len(sendMessageBox.Text) > 0:

					receivedMessageList.Rows = append(receivedMessageList.Rows, "You: "+sendMessageBox.Text)

					writeMessage(conn, username+": "+sendMessageBox.Text)
					sendMessageBox.Text = ""
					receivedMessageList.ScrollBottom()
					ui.Render(receivedMessageList)
					ui.Render(sendMessageBox)
				default:
					sendMessageBox.Text += e.ID
				}
			}
		}
	}
}
