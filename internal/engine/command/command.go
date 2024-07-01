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

type Args struct {
	Name     string
	Desc     string
	Optional bool
}

type HandlerFunc func(cmd Command, appID entity.AppID, callerID string, args ...string) CommandResult

type Command struct {
	Emoji       string
	Color       string
	Name        string
	Help        string
	Args        []Args //! should be nil for commands.
	AppIDs      []entity.AppID
	SubCommands []Command
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

func (cmd *Command) SuccessfulResult(message string, a ...interface{}) CommandResult {
	return CommandResult{
		Color:      cmd.Color,
		Title:      fmt.Sprintf("%v %v", cmd.Help, cmd.Emoji),
		Message:    fmt.Sprintf(message, a...),
		Successful: true,
	}
}

func (cmd *Command) FailedResult(message string, a ...interface{}) CommandResult {
	return CommandResult{
		Color:      cmd.Color,
		Title:      fmt.Sprintf("%v %v", cmd.Help, cmd.Emoji),
		Message:    fmt.Sprintf(message, a...),
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

func (cmd *Command) CheckArgs(input []string) error {
	minArg := len(cmd.Args)
	maxArg := len(cmd.Args)

	for _, arg := range cmd.Args {
		if arg.Optional {
			minArg--
		}
	}

	if len(input) < minArg || len(input) > maxArg {
		return fmt.Errorf("incorrect number of arguments, expected %d but got %d", minArg, len(input))
	}

	return nil
}

func (cmd *Command) HasAppId(appID entity.AppID) bool {
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

func (cmd *Command) AddSubCommand(subCmd Command) {
	if subCmd.HasSubCommand() {
		subCmd.AddHelpSubCommand()
	}

	cmd.SubCommands = append(cmd.SubCommands, subCmd)
}

func (cmd *Command) AddHelpSubCommand() {
	helpCmd := Command{
		Name:   "help",
		Help:   fmt.Sprintf("Help for %v command", cmd.Name),
		AppIDs: entity.AllAppIDs(),
		Handler: func(_ Command, _ entity.AppID, _ string, _ ...string) CommandResult {
			return cmd.SuccessfulResult(cmd.HelpMessage())
		},
	}

	cmd.AddSubCommand(helpCmd)
}
