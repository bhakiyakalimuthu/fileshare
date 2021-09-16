package main

import (
	"bhakiyakalimuthu/fileshare/config"
	"bhakiyakalimuthu/fileshare/internal/grpc"
	"bhakiyakalimuthu/fileshare/internal/pkg"

	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	config.InitConfig()
	l := logger()
	svc := pkg.NewService(l)
	ctrl := pkg.NewController(l, svc)
	go func() {
		ctrl.Init()
	}()
	s := grpc.NewServerGRPC(l, svc)
	sr := s.Init()
	s.Register(sr)
	err := s.Listen()
	if err != nil {
		panic(err)
	}
	defer s.Close()
}

func logger() *zap.Logger {
	enc := zap.NewDevelopmentEncoderConfig()
	enc.EncodeLevel = zapcore.CapitalColorLevelEncoder
	log := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(enc),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	log.Warn("logger setup done")
	return log
}
