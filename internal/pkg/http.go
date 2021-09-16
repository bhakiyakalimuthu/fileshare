package pkg

import (
	"bhakiyakalimuthu/fileshare/config"
	"fmt"
	"net/http"
	"path/filepath"

	"go.uber.org/zap"
)

type Controller struct {
	logger *zap.Logger
	cfg    *config.Config
	svc    Provider
}

func NewController(logger *zap.Logger, svc Provider) *Controller {
	return &Controller{
		logger: logger,
		cfg:    config.Get(),
		svc:    svc,
	}
}

func (c *Controller) Init() {
	//mux := http.NewServeMux()
	//mux.HandleFunc()
	//c.svc.GetFile()
	//fs := http.ServeFile(http.Dir(c.cfg.FilePath))

	absPath, err := filepath.Abs(c.cfg.DestFilePath)
	if err != nil {
		panic(fmt.Sprintf("Init failed to get abs path error %s", err))
	}
	c.logger.Info(fmt.Sprintf("controller initiated on port %s path %s", c.cfg.HttpServerPort, absPath))
	http.Handle("/", http.FileServer(http.Dir(absPath)))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", c.cfg.HttpServerPort), nil); err != nil {
		panic(fmt.Sprintf("ListenAndServer error %s", err))
	}
}
