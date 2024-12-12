package http

import (
	"github.com/pagu-project/Pagu/internal/engine/command"
)

type CommandArgsResponse struct {
	Name     string           `json:"name"`
	Desc     string           `json:"description"`
	InputBox command.InputBox `json:"type"`
	Optional bool             `json:"optional"`
}

type CommandResponse struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Args        []CommandArgsResponse `json:"args"`
	SubCommands []*CommandResponse    `json:"subCommands"`
}

type CommandResult struct {
	Color      string `json:"color"`
	Title      string `json:"title"`
	Error      string `json:"error"`
	Message    string `json:"message"`
	Successful bool   `json:"successful"`
}

type GetCommandsResponse struct {
	Data []*CommandResponse `json:"data"`
}

type CalculateRewardRequest struct {
	Stake string `json:"stake"`
	Days  string `json:"days"`
}

type CalculateFeeRequest struct {
	Amount string `json:"amount"`
}

type NetworkNodeInfoRequest struct {
	ValidatorAddress string `json:"validatorAddress"`
}

type VoucherClaimRequest struct {
	Code    string `json:"code"`
	Address string `json:"address"`
}

type VoucherCreateOneRequest struct {
	Amount      string `json:"amount"`
	ValidMonths string `json:"validMonths"`
	Recipient   string `json:"recipient"`
	Description string `json:"description"`
}

type VoucherCreateBulkRequest struct {
	File   string `json:"file"`
	Notify string `json:"notify"`
}

type VoucherStatusRequest struct {
	Code string `json:"code"`
}

type PhoenixFaucetRequest struct {
	Address string `json:"address"`
}

type BasicResponse struct {
	Result CommandResult `json:"result"`
}
