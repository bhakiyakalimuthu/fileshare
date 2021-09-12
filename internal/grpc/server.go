package grpc

import (
	"bhakiyakalimuthu/fileshare/config"
	"bhakiyakalimuthu/fileshare/internal/pkg"
	pb "bhakiyakalimuthu/fileshare/proto"
	"context"
	"fmt"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"net"
)

type Server interface {
	Listen() error
	Close()
	Register(grpc.ServiceRegistrar)
}

var _ Server = (*ServerGRPC)(nil)

var _ pb.UploaderServer = &ServerGRPC{}

type ServerGRPC struct {
	logger *zap.Logger
	server *grpc.Server
	cfg *config.Config
	svc pkg.Service
	pb.UnimplementedUploaderServer
}
func NewServerGRPC(logger *zap.Logger) Server{
	return &ServerGRPC{
		logger: logger,
		server: nil,
		cfg: config.Get(),
	}
}

func(s *ServerGRPC) Upload(stream pb.Uploader_UploadServer) (err error) {

	var data *pb.UploadRequest
	for {
		data, err = stream.Recv()
		if err!=nil {
			if err== io.EOF{
				break
			}
			return err
		}
	}
	if data.Name == ""{
		if err = s.svc.UploadFile(data.Name, data.Content);err!=nil {
			stream.SendAndClose(&pb.UploadResponse{
				Status:     "Failed to upload file",
				StatusCode: pb.StatusCode_NotOK,
			})
			return
		}
	}
	if err := stream.SendAndClose(&pb.UploadResponse{
		Status:     "Upload file succeeded",
		StatusCode: pb.StatusCode_OK,
	});err!=nil {
		return fmt.Errorf("failed to send status")
	}
	return nil
}

func (s *ServerGRPC) Listen() error {
	l,err := net.Listen("tcp",fmt.Sprintf(":%s",s.cfg.GrpcServerPort))
	if err!=nil {
		return fmt.Errorf("failed to listen %s",err)
	}

	s.server = grpc.NewServer(middleware.WithStreamServerChain(s.getStreamInterceptors()...))
	if err:= s.server.Serve(l);err!=nil {
		return err
	}
	return nil
}

func (s *ServerGRPC) Register(server grpc.ServiceRegistrar) {
	pb.RegisterUploaderServer(server, s)
}

func (s *ServerGRPC) Close() {
	s.server.GracefulStop()

}

func (s *ServerGRPC) getStreamInterceptors() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		grpcrecovery.StreamServerInterceptor(
			grpcrecovery.WithRecoveryHandlerContext(func(ctx context.Context, p interface{}) (err error) {
				s.logger.Error("panic in gRPC call", zap.Stack("stack"))
				return status.Errorf(codes.Internal, "%v", p)
			}),
		),
	}

}
