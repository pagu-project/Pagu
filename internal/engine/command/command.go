package command

import (
	"fmt"
	"slices"

	"github.com/pagu-project/Pagu/internal/entity"
)

var (
	TargetMaskMain      = 1
	TargetMaskTest      = 2
	TargetMaskModerator = 4

	TargetMaskAll = TargetMaskMain | TargetMaskTest | TargetMaskModerator
)

type InputBox int

const (
	InputBoxText InputBox = iota
	InputBoxNumber
	InputBoxFile
	InputBoxAmount
	InputBoxToggle
)

func (i InputBox) Int() int {
	return int(i)
}

type Args struct {
	Name     string
	Desc     string
	InputBox InputBox
	Optional bool
}

type HandlerFunc func(caller *entity.User, cmd *Command, args map[string]string) CommandResult

type Command struct {
	Emoji       string
	Color       string
	Name        string
	Help        string
	Args        []Args // should be nil for commands.
	AppIDs      []entity.AppID
	SubCommands []*Command
	Middlewares []MiddlewareFunc
	Handler     HandlerFunc
	TargetFlag  int
}

type CommandResult struct {
	Color      string
	Title      string
	Error      string
	Message    string
	Successful bool
}

func (cmd *Command) SuccessfulResult(message string, a ...any) CommandResult {
	return CommandResult{
		Color:      cmd.Color,
		Title:      fmt.Sprintf("%v %v", cmd.Help, cmd.Emoji),
		Message:    fmt.Sprintf(message, a...),
		Successful: true,
	}
}

func (cmd *Command) FailedResult(message string, a ...any) CommandResult {
	return CommandResult{
		Color:      cmd.Color,
		Title:      fmt.Sprintf("%v %v", cmd.Help, cmd.Emoji),
		Message:    fmt.Sprintf(message, a...),
		Error:      message,
		Successful: false,
	}
}

func (cmd *Command) ErrorResult(err error) CommandResult {
	return cmd.FailedResult("An error occurred: %v", err.Error())
}

func (cmd *Command) HelpResult() CommandResult {
	return CommandResult{
		Color:      cmd.Color,
		Title:      fmt.Sprintf("%v %v", cmd.Help, cmd.Emoji),
		Message:    cmd.HelpMessage(),
		Successful: false,
	}
}

func (cmd *Command) HasAppID(appID entity.AppID) bool {
	return slices.Contains(cmd.AppIDs, appID)
}

func (cmd *Command) HasSubCommand() bool {
	return len(cmd.SubCommands) > 0 && cmd.SubCommands != nil
}

func (cmd *Command) HelpMessage() string {
	help := cmd.Help
	help += "\n\nAvailable commands:\n"
	for _, sc := range cmd.SubCommands {
		help += fmt.Sprintf("  %-12s %s\n", sc.Name, sc.Help)
	}
	return help
}

func (cmd *Command) AddSubCommand(subCmd *Command) {
	if subCmd.HasSubCommand() {
		subCmd.AddHelpSubCommand()
	}

	cmd.SubCommands = append(cmd.SubCommands, subCmd)
}

func (cmd *Command) AddHelpSubCommand() {
	helpCmd := &Command{
		Name:   "help",
		Help:   fmt.Sprintf("Help for %v command", cmd.Name),
		AppIDs: entity.AllAppIDs(),
		Handler: func(_ *entity.User, _ *Command, _ map[string]string) CommandResult {
			return cmd.SuccessfulResult(cmd.HelpMessage())
		},
	}

	cmd.AddSubCommand(helpCmd)
}
