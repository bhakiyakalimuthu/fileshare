package grpc

import "go.uber.org/zap"

type Uploader interface {

}

type Client struct {
	logger zap.Logger

}