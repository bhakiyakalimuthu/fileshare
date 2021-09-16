package main

import (
	"bhakiyakalimuthu/fileshare/cmd/cli/client"
	"fmt"
	"strings"

	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	l := logger()
	var rootCmd = &cobra.Command{
		Use:   "grpc image uploader",
		Short: "Runs grpc uploader>",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			newArgs, err := parseInputs(args)
			if err != nil {
				l.Error("Upload failed", zap.Error(err))
				return
			}
			client.StartClient(l, newArgs)
		},
	}
	if err := rootCmd.Execute(); err != nil {
		l.Error("Root command execute failed", zap.Error(err))
	}
}

func logger() *zap.Logger {
	enc := zap.NewDevelopmentEncoderConfig()
	enc.EncodeLevel = zapcore.CapitalColorLevelEncoder
	log := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(enc),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	return log
}

func parseInputs(args []string) ([]string, error) {
	fmt.Printf("parse inputs %v", args)
	if len(args) == 0 {
		return nil, fmt.Errorf("please provide atleast one file to be uploaded")
	}
	s := make([]string, 0, len(args))
	for _, arg := range args {
		if strings.HasSuffix(arg, ".jpg") || strings.HasSuffix(arg, ".png") || strings.HasSuffix(arg, ".gif") {
			if strings.HasPrefix(arg, "./") {
				arg = strings.TrimPrefix(arg, "./")
			}
			s = append(s, arg)
			continue
		}
		return nil, fmt.Errorf(" %s not accepted - only files with jpg,png,gif extensions are accepted", arg)
	}
	return s, nil
}
