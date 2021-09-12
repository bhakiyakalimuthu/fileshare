package proto

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

func NewUploadGrpcClient(host, port int) (UploaderClient, error) {

	conn, err := grpc.Dial(
		fmt.Sprintf("kubernetes:///%d:%d", host, port),
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	)
	if err != nil {
		return nil, err
	}

	return NewUploaderClient(conn), nil
}
