package handler

import (
	"fmt"
	"regexp"

	"github.com/regmicmahesh/term-chat/internal/common"
)

type Command struct {
	pattern  *regexp.Regexp
	noOfArgs int
	handler  func(common.Context)
}

type CommandHandler struct {
	Dispatcher interface{}
	commands   []Command
}

func NewCommandHandler(server common.ServerInterface) *CommandHandler {
	ch := &CommandHandler{
		Dispatcher: server,
		commands:   []Command{},
	}
	server.RegisterCommandHandler(ch)
	return ch
}

func (ch *CommandHandler) RegisterCommand(pattern string, noOfArgs int, handler func(common.Context)) {
	ch.commands = append(ch.commands, Command{
		pattern:  regexp.MustCompile(pattern),
		noOfArgs: noOfArgs,
		handler:  handler,
	})
}

func parseRegex(pattern *regexp.Regexp, cmd string) map[string]string {
	result := make(map[string]string)
	subs := pattern.FindStringSubmatch(cmd)
	if subs == nil || len(subs) != pattern.NumSubexp()+1 {
		return nil
	}

	for i, name := range pattern.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = subs[i]
		}
	}

	return result

}

func (ch *CommandHandler) Handle(client common.ClientInterface, command string) bool {
	flag := false
	for _, cmd := range ch.commands {

		parsed := parseRegex(cmd.pattern, command)

		if parsed == nil {
			continue
		}

		context := common.Context{
			Server: ch.Dispatcher.(common.ServerInterface),
			Client: client,
			Args:   parsed,
		}

		cmd.handler(context)
		flag = true
		break
	}

	return flag

}
func (cmdH *CommandHandler) InitCommandHandler() common.CommandHandlerInterface {

	cmdH.RegisterCommand("^/join (?P<name>[a-zA-Z0-9]+)$", 1, JoinChat)
	cmdH.RegisterCommand("^/nick (?P<name>[a-zA-Z0-9]+)$", 1, HandleChangeNickname)
	cmdH.RegisterCommand("^/quit$", 0, QuitChat)
	cmdH.RegisterCommand("^/whisper (?P<target>[a-zA-Z0-9]+) (?P<message>.+)$", 2, HandleWhisper)
	cmdH.RegisterCommand("^/count", 0, GetNumberOfUsers)
	cmdH.RegisterCommand("^/users" , 0, GetUsers)
	cmdH.RegisterCommand("^/whois (?P<target>[a-zA-Z0-9]+)$", 1, Whois)
	cmdH.RegisterCommand("^/register (?P<password>[a-zA-Z0-9]+)$", 1, HandleRegister)
	cmdH.RegisterCommand("^/login (?P<user>[a-zA-Z0-9]+) (?P<password>[a-zA-Z0-9]+)$", 2, HandleLogin)

	return cmdH

}

func HandleRegister(ctx common.Context) {
	username := ctx.Client.GetUsername()

	if username == "Server" {
		ctx.Server.SendServerPrivateMessage(fmt.Sprintf(common.REGISTER_FAILED, username), ctx.Client)
		return
	}

	if ctx.Server.IsUserRegistered(username) {
		ctx.Server.SendServerPrivateMessage(fmt.Sprintf(common.USER_ALREADY_REGISTERED, username), ctx.Client)
		return
	}

	ctx.Server.RegisterUser(username, ctx.Get("password"))

}

func HandleLogin(ctx common.Context) {
	if ctx.Server.IsUserCredentialsValid(ctx.Get("user"), ctx.Get("password")) {
		ctx.Client.SetUsername(ctx.Get("user"))
		ctx.Server.BroadcastServerMessage(fmt.Sprintf(common.LOGIN_SUCCESS, ctx.Get("user")))
	} else {
		ctx.Server.SendServerPrivateMessage(fmt.Sprintf(common.LOGIN_FAILED, ctx.Get("user")), ctx.Client)
	}

}

func HandleChangeNickname(ctx common.Context) {

	if ctx.Server.IsUserRegistered(ctx.Get("name")) {
		ctx.Server.SendServerPrivateMessage(fmt.Sprintf(common.USER_ALREADY_REGISTERED, ctx.Get("name")), ctx.Client)
		return
	}

	if ctx.Server.GetClientByUsername(ctx.Get("name")) != nil {
		ctx.Server.SendServerPrivateMessage(fmt.Sprintf(common.USER_ALREADY_IN_CHAT, ctx.Get("name")), ctx.Client)
		return
	}

	ctx.Server.BroadcastServerMessage(fmt.Sprintf(common.CHANGE_USERNAME, ctx.Client.GetUsername(), ctx.Get("name")))
	ctx.Client.SetUsername(ctx.Get("name"))
}

func JoinChat(ctx common.Context) {
	ctx.Client.SetUsername(ctx.Get("name"))
	ctx.Server.AddClient(ctx.Client)
	ctx.Server.BroadcastServerMessage(fmt.Sprintf(common.JOINED_CHAT, ctx.Client.GetUsername()))
}

func QuitChat(ctx common.Context) {
	ctx.Server.BroadcastServerMessage(fmt.Sprintf(common.LEFT_CHAT, ctx.Client.GetUsername()))
	ctx.Server.RemoveClient(ctx.Client)
}

func HandleWhisper(ctx common.Context) {
	target := ctx.Server.GetClientByUsername(ctx.Get("target"))
	if target != nil {
		ctx.Server.SendServerPrivateMessage(fmt.Sprintf(common.WHISPERED, ctx.Client.GetUsername(), ctx.Get("message")), target.(common.ClientInterface))
	} else {
		ctx.Server.SendServerPrivateMessage(fmt.Sprintf(common.NOT_IN_CHAT, target.GetUsername()), ctx.Client)
	}
}

func GetNumberOfUsers(ctx common.Context) {
	ctx.Server.SendServerPrivateMessage(fmt.Sprintf(common.USER_COUNT, ctx.Server.GetNumberOfUsers()), ctx.Client)
}

func GetUsers(ctx common.Context) {
	ctx.Server.SendServerPrivateMessage(fmt.Sprintf(common.USERS_LIST, ctx.Server.GetUsers()), ctx.Client)
}

func Whois(ctx common.Context) {
	user := ctx.Server.GetClientByUsername(ctx.Get("target"))
	if user != nil {
		ctx.Server.SendServerPrivateMessage(fmt.Sprintf(common.WHOIS, user.GetUsername(), user.GetIPAddr()), ctx.Client)
	} else {
		ctx.Server.SendServerPrivateMessage(fmt.Sprintf(common.NOT_IN_CHAT, ctx.Get("target")), ctx.Client)
	}
}
