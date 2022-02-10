package commands

type Action func(args []interface{}) error

type Command struct {
	Name    string
	Command string
	Action
}

var ChangeNicknameCommand = Command{
	Name:    "change nickname",
	Command: "/nick",
	Action:  nil,
}

func NewChangeNicknameCommand(action Action) Command {
	return Command{
		Name:    "change nickname",
		Command: "/nick",
		Action:  action,
	}
}

func handle(command Command, args []interface{}) error {
	return command.Action(args)
}

func HandleCommand(msg string, args []interface{}) error {
	switch(msg){
	case ChangeNicknameCommand.Command:
		return handle(ChangeNicknameCommand, args)
	}
	return nil
}
