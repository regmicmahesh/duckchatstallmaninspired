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
var ipToColor = make(map[string]string)
var possibleColors = []string{
	"red",
	"blue",
	"black",
	"cyan",
	"yellow",
	"white",
	"clear",
	// "green", // HardCoded for server messages
	// "magenta", // HardCoded for self
}

func readMessage(conn net.Conn, message chan string) (string, error) {
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		ip := conn.LocalAddr().String()
		_, ok := ipToColor[ip]
		if !ok {
			rand.Seed(time.Now().UnixNano())
			ipToColor[ip] = possibleColors[rand.Intn(len(possibleColors))]
		}
		color := ipToColor[ip]

		msg = strings.TrimSpace(msg)
		if msg != "1" {
			if strings.Split(msg, ":")[0] == "Server" {
				message <- (fmt.Sprintf("[%s](fg:green,mod:bold)", msg))
			} else {
				message <- (fmt.Sprintf("[%s](fg:%s)", msg, color))
			}
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
	receivedMessageList.Title = "Messages"
	receivedMessageList.TitleStyle.Fg = ui.ColorMagenta
	receivedMessageList.Rows = []string{
		"[+] Connected to server [+] ",
	}
	receivedMessageList.TextStyle = ui.NewStyle(ui.ColorYellow)
	receivedMessageList.WrapText = false
	receivedMessageList.SetRect(0, 0, Mx, My-3)
	ui.Render(receivedMessageList)

	sendMessageBox := widgets.NewParagraph()
	sendMessageBox.Title = "Enter Your Message"
	sendMessageBox.TitleStyle.Fg = ui.ColorMagenta
	sendMessageBox.Text = ""
	sendMessageBox.SetRect(0, My-3, Mx, My)
	ui.Render(sendMessageBox)

	msg := make(chan string)
	go readMessage(conn, msg)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Millisecond * 17).C

	for {
		select {
		case recievedMsg := <-msg:
			receivedMessageList.Rows = append(receivedMessageList.Rows, recievedMsg)
			receivedMessageList.ScrollBottom()
			ui.Render(receivedMessageList)

		case <-ticker:
			ui.Render(sendMessageBox, receivedMessageList)

		case e := <-uiEvents:
			switch e.Type {
			case ui.KeyboardEvent:
				switch e.ID {
				case "<C-c>":
					return
				case "<C-n>":
					receivedMessageList.ScrollDown()
				case "<C-p>":
					receivedMessageList.ScrollUp()
				case "<Space>":
					sendMessageBox.Text += " "
				case "<Backspace>":
					if len(sendMessageBox.Text) > 0 {
						sendMessageBox.Text = sendMessageBox.Text[:len(sendMessageBox.Text)-1]
					}
				case "<Enter>":
					if len(sendMessageBox.Text) > 0 {
						receivedMessageList.Rows = append(receivedMessageList.Rows, fmt.Sprintf("[You: %s](fg:magenta)", sendMessageBox.Text))
						writeMessage(conn, sendMessageBox.Text)
						sendMessageBox.Text = ""
						receivedMessageList.ScrollBottom()
					}
				default:
					if len(e.ID) != 1 && e.ID[0] == '<' {
						// Do nothing
					} else {
						sendMessageBox.Text += e.ID
					}
				}

			case ui.MouseEvent:
				switch e.ID {
				case "<MouseWheelUp>":
					receivedMessageList.ScrollUp()
				case "<MouseWheelDown>":
					receivedMessageList.ScrollDown()
				}

			case ui.ResizeEvent:
				payload := e.Payload.(ui.Resize)
				Mx, My := payload.Width, payload.Height

				receivedMessageList.SetRect(0, 0, Mx, My-3)
				sendMessageBox.SetRect(0, My-3, Mx, My)
			}
		}
	}
}
