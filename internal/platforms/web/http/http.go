package http

import (
	"github.com/labstack/echo/v4"
	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/engine"
	"github.com/pagu-project/Pagu/pkg/log"
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
	hs.eServer.GET("/commands", hs.handler.GetCommands)
	hs.eServer.POST("/api/v1/calculate/reward", hs.handler.CalculateReward)
	hs.eServer.POST("/api/v1/calculate/fee", hs.handler.CalculateFee)
	hs.eServer.POST("/api/v1/calculate/help", hs.handler.CalculateHelp)
	hs.eServer.POST("/api/v1/netwrok/health", hs.handler.NetworkHealth)
	hs.eServer.POST("/api/v1/netwrok/node-info", hs.handler.NetworkNodeInfo)
	hs.eServer.POST("/api/v1/netwrok/status", hs.handler.NetworkNodeStatus)
	hs.eServer.POST("/api/v1/netwrok/help", hs.handler.NetworkHelp)
	hs.eServer.POST("/api/v1/voucher/claim", hs.handler.VoucherClaim)
	hs.eServer.POST("/api/v1/voucher/create-one", hs.handler.VoucherCreateOne)
	hs.eServer.POST("/api/v1/voucher/create-bulk", hs.handler.VoucherCreateBulk)
	hs.eServer.POST("/api/v1/voucher/status", hs.handler.VoucherStatus)
	hs.eServer.POST("/api/v1/voucher/help", hs.handler.VoucherHelp)
	hs.eServer.POST("/api/v1/market/price", hs.handler.MarketPrice)
	hs.eServer.POST("/api/v1/market/help", hs.handler.MarketHelp)
	hs.eServer.POST("/api/v1/phoenix/faucet", hs.handler.PhoenixFaucet)
	hs.eServer.POST("/api/v1/phoenix/status", hs.handler.PhoenixStatus)
	hs.eServer.POST("/api/v1/phoenix/help", hs.handler.PhoenixHelp)
	hs.eServer.POST("/api/v1/help", hs.handler.Help)
	return hs.eServer.Start(hs.cfg.Listen)
}

func (hs *HTTPServer) Stop() error {
	log.Info("Stopping HTTP Server")
	return hs.eServer.Close()
}
