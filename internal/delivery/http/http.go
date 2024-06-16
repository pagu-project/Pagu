package http

import (
	"net/http"
	"strings"

	"github.com/pagu-project/Pagu/internal/entity"

	"github.com/pagu-project/Pagu/internal/engine"
	"github.com/pagu-project/Pagu/pkg/log"

	"github.com/labstack/echo/v4"
	"github.com/pagu-project/Pagu/config"
)

type HTTPServer struct {
	handler HTTPHandler
	eServer *echo.Echo
	cfg     config.HTTP
}

type HTTPHandler struct {
	engine *engine.BotEngine
}

func NewHTTPServer(be *engine.BotEngine, cfg config.HTTP) HTTPServer {
	return HTTPServer{
		handler: HTTPHandler{
			engine: be,
		},
		eServer: echo.New(),
		cfg:     cfg,
	}
}

func (hs *HTTPServer) Start() error {
	log.Info("Starting HTTP Server", "listen", hs.cfg.Listen)
	hs.eServer.POST("/run", hs.handler.Run)
	return hs.eServer.Start(hs.cfg.Listen)
}

type RunRequest struct {
	Command string `json:"command"`
}

type RunResponse struct {
	Result string `json:"result"`
}

func (hh *HTTPHandler) Run(c echo.Context) error {
	r := new(RunRequest)
	if err := c.Bind(r); err != nil {
		return err
	}

	beInput := []string{}

	tokens := strings.Split(r.Command, " ")
	beInput = append(beInput, tokens...)

	cmdResult := hh.engine.Run(entity.AppIdHTTP, c.RealIP(), beInput)

	return c.JSON(http.StatusOK, RunResponse{
		Result: cmdResult.Message,
	})
}

func (hs *HTTPServer) Stop() error {
	log.Info("Stopping HTTP Server")
	return hs.eServer.Close()
}
