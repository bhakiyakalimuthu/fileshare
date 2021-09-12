package pkg

import (
	"bhakiyakalimuthu/fileshare/config"
	"go.uber.org/zap"
	"net/http"
)

type Controller struct {
	logger *zap.Logger
	cfg *config.Config
}

func NewController(logger *zap.Logger) *Controller {
	return &Controller{
		logger:logger,
		cfg : config.Get(),
	}
}


func(c *Controller) Init(){
	http.Handle("/", http.FileServer(http.Dir(c.cfg.FilePath)))
}