package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var wg sync.WaitGroup
var nameToColor = make(map[string]string)
var intToColor = []string{
	"red",
	"blue",
	"black",
	"cyan",
	"yellow",
	"white",
	"clear",
	"green",
	"magenta",
}

func readMessage(conn net.Conn, message chan string) (string, error) {
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		name := conn.LocalAddr().String()
		_, ok := nameToColor[name]
		if !ok {
			rand.Seed(time.Now().UnixNano())
			nameToColor[name] = intToColor[rand.Intn(8)]
		}
		color := nameToColor[name]

		msg = strings.TrimSpace(msg)
		if msg != "1" {
			message <- (fmt.Sprintf("[%s](fg:%s)", msg, color))
		}
	}
}

func writeMessage(conn net.Conn, message string) error {
	writer := bufio.NewWriter(conn)
	_, err := writer.WriteString(message + "\n")
	writer.Flush()
	if err != nil {
		fmt.Println("Error writing to server:", err.Error())
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: client <server-address>")
		os.Exit(1)
	}
	serverAddress := os.Args[1]

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
					receivedMessageList.Rows = append(receivedMessageList.Rows, fmt.Sprintf("[You: %s](fg:magenta)", sendMessageBox.Text))
					writeMessage(conn, sendMessageBox.Text)
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
