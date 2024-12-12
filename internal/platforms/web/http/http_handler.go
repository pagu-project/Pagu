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

// GetCommands retrieves the list of commands associated with this object.
// @Summary Get available commands
// @Description Retrieve a list of all available commands
// @Tags Commands
// @Produce json
// @Success 200 {object} GetCommandsResponse
// @Router /commands [get].
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

// CalculateReward calculates the reward for a node.
// @Summary Calculate reward
// @Description Calculate the reward based on the provided stake and number of days
// @Tags Calculate
// @Accept json
// @Produce json
// @Param request body CalculateRewardRequest true "Calculate Reward Request"
// @Success 200 {object} BasicResponse "Calculation result"
// @Router /calculate/reward [post].
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

// CalculateFee calculates transaction fee based on provided amount.
// @Summary Calculate fee
// @Description Calculate the fee based on the provided amount
// @Tags Calculate
// @Accept json
// @Produce json
// @Param request body CalculateFeeRequest true "Calculate Fee Request"
// @Success 200 {object} BasicResponse "Calculation result"
// @Router /calculate/fee [post].
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

// CalculateHelp returns help for calcuate commands.
// @Summary Calculate help information
// @Description Calculate help information
// @Tags Calculate
// @Accept json
// @Produce json
// @Success 200 {object} BasicResponse "Successful calculation result"
// @Router /calculate/help [get].
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

// NetworkHealth returns the health status.
// @Summary Get network health
// @Description Retrieve the health status of the network
// @Tags Network
// @Produce json
// @Success 200 {object} BasicResponse
// @Router /network/health [get].
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

// NetworkNodeInfo returns info of a node.
// @Summary Get Network Node Information
// @Description Retrieves information about a network node based on the provided validator address.
// @Tags Network
// @Accept json
// @Produce json
// @Param request body NetworkNodeInfoRequest true "Validator Address"
// @Success 200 {object} BasicResponse "The result of the command"
// @Router /network/node-info [post].
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

// NetworkNodeStatus returns status of network.
// @Summary Get network node status
// @Description Retrieves the current status of the network node.
// @Tags Network
// @Accept json
// @Produce json
// @Param body object true "Request body"
// @Success 200 {object} BasicResponse "Successfully retrieved network node status"
// @Router /network/status [post].
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

// NetworkHelp return help data for network commands.
// @Summary Provides network help information
// @Description Provides network help information
// @Tags Network
// @Accept  json
// @Produce  json
// @Success 200 {object} BasicResponse{result=CommandResult} "Successfully retrieved network help"
// @Router /network/help [get].
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

// VoucherClaim claim a voucher.
// @Summary Claim a voucher using the provided code and address.
// @Description This endpoint allows users to claim a voucher by providing a valid code and their address.
// @Tags Voucher
// @Accept json
// @Produce json
// @Param voucherClaimRequest body VoucherClaimRequest true "Voucher Claim Request"
// @Success 200 {object} BasicResponse "Voucher claim successful"
// @Router /voucher/claim [post].
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

// VoucherCreateOne creates a voucher.
// @Summary Creates a new voucher
// @Description Creates a voucher.
// @Tags Voucher
// @Accept json
// @Produce json
// @Param voucherCreateOneRequest body VoucherCreateOneRequest true "Voucher details"
// @Success 200 {object} BasicResponse "Voucher creation successful"
// @Router /voucher/create-one [post].
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

// VoucherCreateBulk crreates bulk vouchers.
// @Summary Create vouchers in bulk
// @Description This API allows you to create multiple vouchers.
// @Tags Voucher
// @Accept json
// @Produce json
// @Param file formData string true "File containing voucher data"
// @Param notify formData string false "Notification flag (optional)"
// @Success 200 {object} BasicResponse "Command result containing the status of the bulk voucher creation"
// @Router /vouchers/create-bulk [post].
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

// VoucherStatus checks the status of a voucher.
// @Summary Check the status of a voucher
// @Description Accepts a voucher code and returns the status of the voucher.
// @Tags Voucher
// @Accept json
// @Produce json
// @Param request body VoucherStatusRequest true "Voucher Status Request"
// @Success 200 {object} BasicResponse "Voucher status response"
// @Router /voucher/status [post].
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

// VoucherHelp returns help information for the voucher command.
// @Summary Get help information for the voucher command
// @Description Executes the 'voucher help' command and returns the result.
// @Tags Voucher
// @Accept json
// @Produce json
// @Success 200 {object} BasicResponse "Successful response with command result"
// @Router /voucher/help [get].
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

// MarketPrice returns market price.
// @Summary Get Market Price
// @Description Retrieves the current market price information
// @Tags Market
// @Accept  json
// @Produce  json
// @Success 200 {object} BasicResponse{result=CommandResult} "Successfully retrieved the market price"
// @Router /market/price [get].
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

// MarketHelp returns market help information.
// @Summary Get help for the market command
// @Description This endpoint processes the "market help" command and returns relevant details about the market command.
// @Tags Market
// @Accept  json
// @Produce  json
// @Success 200 {object} BasicResponse{Result=CommandResult} "Command help details"
// @Router /market/help [get].
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

// PhoenixFaucet initiates a faucet request for the Phoenix.
// @Summary Initiates a faucet request for the Phoenix
// @Description Initiates a faucet request for the Phoenix
// @Tags Phoenix
// @Accept json
// @Produce json
// @Param request body PhoenixFaucetRequest true "Faucet request with address"
// @Success 200 {object} BasicResponse "Successful response containing faucet status"
// @Router /phoenix/faucet [post].
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

// PhoenixStatus returns phonix status.
// @Summary Get Phoenix Status
// @Description Get Phoenix Status
// @Tags Phoenix
// @Accept  json
// @Produce  json
// @Success 200 {object} BasicResponse "Status of the Phoenix application"
// @Router /phoenix/status [get].
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

// PhoenixHelp returns phoenix help information.
// @Summary Executes Phoenix help command
// @Description Executes Phoenix help command.
// @Tags Phoenix
// @Accept json
// @Produce json
// @Success 200 {object} BasicResponse{result=CommandResult}
// @Router /phoenix/help [get].
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

// Help returns help information.
// @Summary Executes a help command and returns the result.
// @Description This endpoint runs the 'help' command and returns the result.
// @Tags Help
// @Accept json
// @Produce json
// @Success 200 {object} BasicResponse{Result=CommandResult} "Help command result"
// @Router /help [get].
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
