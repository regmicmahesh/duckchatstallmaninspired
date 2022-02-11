package handler

import (
	"fmt"
	"regexp"

	i "github.com/regmicmahesh/term-chat/internal/common"
)


func HandleChangeNickname(ctx i.Context) {
	ctx.Server.BroadcastServerMessage(fmt.Sprintf("%s changed their name to %s", ctx.Client.GetUsername(), ctx.Get("name")))
	ctx.Client.SetUsername(ctx.Get("name"))
}

func JoinChat(ctx i.Context) {
	ctx.Server.AddClient(ctx.Client)
	ctx.Server.BroadcastServerMessage(fmt.Sprintf("%s joined the chat", ctx.Client.GetUsername()))
}

func QuitChat(ctx i.Context) {
	ctx.Server.BroadcastServerMessage(fmt.Sprintf("%s left the chat", ctx.Client.GetUsername()))
	ctx.Server.RemoveClient(ctx.Client)
}


func HandleWhisper(ctx i.Context) {
	target := ctx.Server.GetClientByUsername(ctx.Get("target"))
	if target != nil {
		ctx.Server.SendServerPrivateMessage(fmt.Sprintf("%s whispered to you: %s", ctx.Client.GetUsername(), ctx.Get("message")), target.(i.ClientInterface))
	} else {
		ctx.Server.SendServerPrivateMessage(fmt.Sprintf("%s is not in the chat.", target.GetUsername()), ctx.Client)
	}
}

func GetNumberOfUsers(ctx i.Context) {
	ctx.Server.SendServerPrivateMessage(fmt.Sprintf("%d users are in the chat", ctx.Server.GetNumberOfUsers()), ctx.Client)
}

func Whois(ctx i.Context) {
	user := ctx.Server.GetClientByUsername(ctx.Get("target"))
	ctx.Server.SendServerPrivateMessage(fmt.Sprintf("%s is connected from %s", user.GetUsername(), user.GetIPAddr()), ctx.Client)
}


type Command struct {
	pattern  *regexp.Regexp
	noOfArgs int
	handler  func(i.Context)
}

type CommandHandler struct {
	Server   i.ServerInterface
	commands []Command
}

func NewCommandHandler(server i.ServerInterface) *CommandHandler {
	ch := &CommandHandler{
		Server:   server,
		commands: []Command{},
	}
	server.RegisterCommandHandler(ch)
	return ch
}



func (ch *CommandHandler) RegisterCommand(pattern string, noOfArgs int, handler func(i.Context)) {
	ch.commands = append(ch.commands, Command{
		pattern:  regexp.MustCompile(pattern),
		noOfArgs: noOfArgs,
		handler:  handler,
	})
}

func (ch *CommandHandler) Handle(client i.ClientInterface, command string) bool {
	flag := false
	for _, cmd := range ch.commands {
		result := make(map[string]string)
		subs := cmd.pattern.FindStringSubmatch(command)
		if subs == nil {
			continue
		}

		if len(subs) != cmd.noOfArgs+1 {
			continue
		}
		for i, name := range cmd.pattern.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = subs[i]
			}
		}

		context := i.Context{
			Server: ch.Server,
			Client: client,
			Args:   result,
		}

		cmd.handler(context)
		flag = true
		break
	}

	return flag

}
func (cmdH *CommandHandler) InitCommandHandler() i.CommandHandlerInterface {

	cmdH.RegisterCommand("^/join$", 0, JoinChat)
	cmdH.RegisterCommand("^/nick (?P<name>[a-zA-Z0-9]+)$", 1, HandleChangeNickname)
	cmdH.RegisterCommand("^/quit$", 0, QuitChat)
	cmdH.RegisterCommand("^/whisper (?P<target>[a-zA-Z0-9]+) (?P<message>.+)$", 2, HandleWhisper)
	cmdH.RegisterCommand("^/users$", 0, GetNumberOfUsers)
	cmdH.RegisterCommand("^/whois (?P<target>[a-zA-Z0-9]+)$", 1, Whois)

	return cmdH

}

