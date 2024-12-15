package http

import (
	// External packages.
	"github.com/labstack/echo/v4"
	// Internal packages.
	"github.com/pagu-project/Pagu/config"
	// Importing Pagu/docs to include API documentation in the binary.
	_ "github.com/pagu-project/Pagu/docs"
	"github.com/pagu-project/Pagu/internal/engine"
	"github.com/pagu-project/Pagu/pkg/log"
	swagger "github.com/swaggo/echo-swagger"
)

type HTTPServer struct {
	handler HTTPHandler
	eServer *echo.Echo
	cfg     *config.HTTP
}

func NewHTTPServer(be *engine.BotEngine, cfg *config.HTTP) HTTPServer {
	return HTTPServer{
		handler: NewHTTPHandler(be),
		eServer: echo.New(),
		cfg:     cfg,
	}
}

func (hs *HTTPServer) Start() error {
	log.Info("Starting HTTP Server", "listen", hs.cfg.Listen)
	hs.eServer.GET("/swagger/*", swagger.WrapHandler)
	hs.eServer.GET("/api/v1/commands", hs.handler.GetCommands)
	hs.eServer.POST("/api/v1/calculate/reward", hs.handler.CalculateReward)
	hs.eServer.POST("/api/v1/calculate/fee", hs.handler.CalculateFee)
	hs.eServer.GET("/api/v1/calculate/help", hs.handler.CalculateHelp)
	hs.eServer.GET("/api/v1/network/health", hs.handler.NetworkHealth)
	hs.eServer.POST("/api/v1/network/node-info", hs.handler.NetworkNodeInfo)
	hs.eServer.POST("/api/v1/network/status", hs.handler.NetworkNodeStatus)
	hs.eServer.GET("/api/v1/network/help", hs.handler.NetworkHelp)
	hs.eServer.POST("/api/v1/voucher/claim", hs.handler.VoucherClaim)
	hs.eServer.POST("/api/v1/voucher/create-one", hs.handler.VoucherCreateOne)
	hs.eServer.POST("/api/v1/voucher/create-bulk", hs.handler.VoucherCreateBulk)
	hs.eServer.POST("/api/v1/voucher/status", hs.handler.VoucherStatus)
	hs.eServer.GET("/api/v1/voucher/help", hs.handler.VoucherHelp)
	hs.eServer.GET("/api/v1/market/price", hs.handler.MarketPrice)
	hs.eServer.GET("/api/v1/market/help", hs.handler.MarketHelp)
	hs.eServer.POST("/api/v1/phoenix/faucet", hs.handler.PhoenixFaucet)
	hs.eServer.GET("/api/v1/phoenix/status", hs.handler.PhoenixStatus)
	hs.eServer.GET("/api/v1/phoenix/help", hs.handler.PhoenixHelp)
	hs.eServer.GET("/api/v1/help", hs.handler.Help)
	return hs.eServer.Start(hs.cfg.Listen)
}

func (hs *HTTPServer) Stop() error {
	log.Info("Stopping HTTP Server")
	return hs.eServer.Close()
}
