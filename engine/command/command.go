package command

import (
	"fmt"
	"slices"
)

type AppID int

const (
	AppIdCLI      AppID = 1
	AppIdDiscord  AppID = 2
	AppIdgRPC     AppID = 3
	AppIdHTTP     AppID = 4
	AppIdTelegram AppID = 5
)

func (appID AppID) String() string {
	switch appID {
	case AppIdCLI:
		return "CLI"
	case AppIdDiscord:
		return "Discord"
	case AppIdgRPC:
		return "gRPC"
	case AppIdHTTP:
		return "HTTP"
	case AppIdTelegram:
		return "Telegram"
	}

	return ""
}

func AllAppIDs() []AppID {
	return []AppID{
		AppIdCLI,
		AppIdDiscord,
		AppIdgRPC,
		AppIdHTTP,
		AppIdTelegram,
	}
}

type Args struct {
	Name     string
	Desc     string
	Optional bool
}

type Command struct {
	Emoji       string
	Color       string
	Name        string
	Desc        string
	Help        string
	Args        []Args //! should be nil for commands.
	AppIDs      []AppID
	SubCommands []Command
	Handler     func(cmd Command, source AppID, callerID string, args ...string) CommandResult
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
		Title:      fmt.Sprintf("%v %v", cmd.Desc, cmd.Emoji),
		Message:    fmt.Sprintf(message, a...),
		Successful: true,
	}
}

func (cmd *Command) FailedResult(message string, a ...interface{}) CommandResult {
	return CommandResult{
		Color:      cmd.Color,
		Title:      fmt.Sprintf("%v %v", cmd.Desc, cmd.Emoji),
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
		Title:      fmt.Sprintf("%v %v", cmd.Desc, cmd.Emoji),
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

func (cmd *Command) HasAppId(appID AppID) bool {
	return slices.Contains(cmd.AppIDs, appID)
}

func (cmd *Command) HasSubCommand() bool {
	return len(cmd.SubCommands) > 0 && cmd.SubCommands != nil
}

func (cmd *Command) HelpMessage() string {
	help := cmd.Help
	help += "\n\nAvailable commands:\n"
	for _, sc := range cmd.SubCommands {
		help += fmt.Sprintf("  %-12s %s\n", sc.Name, sc.Desc)
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
		Desc:   fmt.Sprintf("Help for %v command", cmd.Name),
		AppIDs: AllAppIDs(),
		Handler: func(_ Command, _ AppID, _ string, _ ...string) CommandResult {
			return cmd.SuccessfulResult(cmd.HelpMessage())
		},
	}

	cmd.AddSubCommand(helpCmd)
}
