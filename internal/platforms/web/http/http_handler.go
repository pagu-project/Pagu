package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pagu-project/Pagu/internal/engine"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/log"
)

type HTTPHandler struct {
	engine *engine.BotEngine
}

func NewHTTPHandler(be *engine.BotEngine) HTTPHandler {
	return HTTPHandler{
		engine: be,
	}
}

func (hh *HTTPHandler) GetCommands(c echo.Context) error {
	log.Info("New request received for GetCommands.")
	commands := hh.engine.Commands()
	response := make([]*CommandResponse, len(commands))
	for i, cmd := range commands {
		response[i] = mapCommandToResponse(cmd)
	}
	return c.JSON(http.StatusOK, GetCommandsResponse{
		Data: response,
	})
}

func (hh *HTTPHandler) CalculateReward(c echo.Context) error {
	r := new(CalculateRewardRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	commands := []string{"calculate", "reward"}
	args := make(map[string]string)
	args["stake"] = r.Stake
	args["days"] = r.Days
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, args)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) CalculateFee(c echo.Context) error {
	r := new(CalculateFeeRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	commands := []string{"calculate", "fee"}
	args := make(map[string]string)
	args["amount"] = r.Amount
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, args)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) CalculateHelp(c echo.Context) error {
	commands := []string{"calculate", "help"}
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, nil)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) NetworkHealth(c echo.Context) error {
	commands := []string{"network", "health"}
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, nil)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) NetworkNodeInfo(c echo.Context) error {
	r := new(NetworkNodeInfoRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	commands := []string{"network", "node-info"}
	args := make(map[string]string)
	args["validator_address"] = r.ValidatorAddress
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, args)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) NetworkNodeStatus(c echo.Context) error {
	commands := []string{"network", "status"}
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, nil)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) NetworkHelp(c echo.Context) error {
	commands := []string{"network", "help"}
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, nil)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) VoucherClaim(c echo.Context) error {
	r := new(VoucherClaimRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	commands := []string{"voucher", "claim"}
	args := make(map[string]string)
	args["code"] = r.Code
	args["address"] = r.Address
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, args)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) VoucherCreateOne(c echo.Context) error {
	r := new(VoucherCreateOneRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	commands := []string{"voucher", "create-one"}
	args := make(map[string]string)
	args["amount"] = r.Amount
	args["valid-months"] = r.ValidMonths
	args["recipient"] = r.Recipient
	args["description"] = r.Description
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, args)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) VoucherCreateBulk(c echo.Context) error {
	r := new(VoucherCreateBulkRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	commands := []string{"voucher", "create-bulk"}
	args := make(map[string]string)
	args["file"] = r.File
	args["notify"] = r.Notify
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, args)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) VoucherStatus(c echo.Context) error {
	r := new(VoucherStatusRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	commands := []string{"voucher", "status"}
	args := make(map[string]string)
	args["code"] = r.Code
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, args)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) VoucherHelp(c echo.Context) error {
	commands := []string{"voucher", "help"}
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, nil)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) MarketPrice(c echo.Context) error {
	commands := []string{"market", "price"}
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, nil)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) MarketHelp(c echo.Context) error {
	commands := []string{"market", "help"}
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, nil)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) PhoenixFaucet(c echo.Context) error {
	r := new(PhoenixFaucetRequest)
	if err := c.Bind(r); err != nil {
		return err
	}
	commands := []string{"phoenix", "faucet"}
	args := make(map[string]string)
	args["address"] = r.Address
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, args)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) PhoenixStatus(c echo.Context) error {
	commands := []string{"phoenix", "status"}
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, nil)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) PhoenixHelp(c echo.Context) error {
	commands := []string{"phoenix", "help"}
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, nil)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func (hh *HTTPHandler) Help(c echo.Context) error {
	commands := []string{"help"}
	cmdResult := hh.engine.Run(entity.AppIDHTTP, c.RealIP(), commands, nil)
	return c.JSON(http.StatusOK, BasicResponse{
		Result: CommandResult{
			Color:      cmdResult.Color,
			Title:      cmdResult.Title,
			Error:      cmdResult.Error,
			Message:    cmdResult.Message,
			Successful: cmdResult.Successful,
		},
	})
}

func mapCommandToResponse(cmd *command.Command) *CommandResponse {
	if cmd == nil {
		return nil
	}
	argsResponses := make([]CommandArgsResponse, len(cmd.Args))
	for i, arg := range cmd.Args {
		argsResponses[i] = CommandArgsResponse{
			Name:     arg.Name,
			Desc:     arg.Desc,
			InputBox: arg.InputBox,
			Optional: arg.Optional,
		}
	}
	subCommandResponses := make([]*CommandResponse, len(cmd.SubCommands))
	for i, subCmd := range cmd.SubCommands {
		subCommandResponses[i] = mapCommandToResponse(subCmd)
	}
	return &CommandResponse{
		Name:        cmd.Name,
		Description: cmd.Help,
		Args:        argsResponses,
		SubCommands: subCommandResponses,
	}
}
