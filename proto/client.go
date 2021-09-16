package proto

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

func NewUploadGrpcClient(host string, port int) (UploaderClient, func(), error) {

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithInsecure(),
		grpc.WithBalancerName(roundrobin.Name),
	)
	if err != nil {
		return nil, nil, err
	}
	closer := func(){ _ = conn.Close()}
	return NewUploaderClient(conn), closer, nil
}
