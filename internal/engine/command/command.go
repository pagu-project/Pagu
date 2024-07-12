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

const (
	CommandParamTypeString      uint8 = 3
	CommandParamTypeInteger     uint8 = 4
	CommandParamTypeBoolean     uint8 = 5
	CommandParamTypeUser        uint8 = 6
	CommandParamTypeChannel     uint8 = 7
	CommandParamTypeRole        uint8 = 8
	CommandParamTypeMentionable uint8 = 9
	CommandParamTypeNumber      uint8 = 10
	CommandParamTypeAttachment  uint8 = 11
)

type Args struct {
	Name     string
	Desc     string
	Type     uint8
	Optional bool
}

type HandlerFunc func(cmd *Command, appID entity.AppID, callerID string, args map[string]any) CommandResult

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
	User        *entity.User
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
		Handler: func(_ *Command, _ entity.AppID, _ string, _ map[string]any) CommandResult {
			return cmd.SuccessfulResult(cmd.HelpMessage())
		},
	}

	cmd.AddSubCommand(helpCmd)
}
