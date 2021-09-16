package client

import (
	"bhakiyakalimuthu/fileshare/config"
	"bhakiyakalimuthu/fileshare/internal/grpc"
	"bhakiyakalimuthu/fileshare/proto"
	"context"

	"go.uber.org/zap"
)

func StartClient(l *zap.Logger, fileName []string) {
	//TODO: loop through inputs
	config.InitConfig()
	client, close, err := proto.NewUploadGrpcClient(config.Get().Host, config.Get().GrpcServerPort)
	defer close()

	if err != nil {
		panic(err)
	}
	c := grpc.NewClient(l, client)
	c.Upload(context.Background(), fileName[0])
}
