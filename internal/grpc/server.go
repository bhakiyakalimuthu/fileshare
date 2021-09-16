package grpc

import (
	"bhakiyakalimuthu/fileshare/config"
	"bhakiyakalimuthu/fileshare/internal/pkg"
	pb "bhakiyakalimuthu/fileshare/proto"
	"context"
	"fmt"
	"io"
	"net"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server interface {
	Listen() error
	Close()
	Init() *grpc.Server
	Register(grpc.ServiceRegistrar)
}

var _ Server = (*ServerGRPC)(nil)

var _ pb.UploaderServer = &ServerGRPC{}

type ServerGRPC struct {
	logger *zap.Logger
	server *grpc.Server
	cfg    *config.Config
	svc    pkg.Provider
	pb.UnimplementedUploaderServer
}

func NewServerGRPC(logger *zap.Logger, svc pkg.Provider) Server {
	return &ServerGRPC{
		logger: logger,
		cfg:    config.Get(),
		svc:    svc,
	}
}

func (s *ServerGRPC) Upload(stream pb.Uploader_UploadServer) (err error) {

	s.logger.Warn("Server - upload started")
	var data *pb.UploadRequest
	firstStream := true
	for {
		data, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if data.Name == "" {
			stream.SendAndClose(&pb.UploadResponse{
				Status:     "File is not provided",
				StatusCode: pb.StatusCode_NotOK,
			})
			return
		}
		if firstStream {
			if err = s.svc.UploadFile(data.Name, data.Content); err != nil {
				s.logger.Error("failed to upload file", zap.Error(err))
				stream.SendAndClose(&pb.UploadResponse{
					Status:     "Failed to upload file",
					StatusCode: pb.StatusCode_NotOK,
				})
				return
			}
			firstStream = false
		}
	}
	if err := stream.SendAndClose(&pb.UploadResponse{
		Status:     "Upload file succeeded",
		StatusCode: pb.StatusCode_OK,
	}); err != nil {
		s.logger.Error("failed to send status", zap.Error(err))
		return fmt.Errorf("failed to send status %s", err)
	}
	return nil
}

func (s *ServerGRPC) Listen() error {

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GrpcServerPort))
	if err != nil {
		return fmt.Errorf("failed to listen %s", err)
	}

	s.logger.Warn(fmt.Sprintf("grpc server started listening on %s", l.Addr().String()))
	if err := s.server.Serve(l); err != nil {
		return err
	}
	s.logger.Warn("grpc server started listening")
	return nil
}
func (s *ServerGRPC) Init() *grpc.Server {
	s.server = grpc.NewServer(middleware.WithStreamServerChain(s.getStreamInterceptors()...))
	return s.server
}

func (s *ServerGRPC) Register(server grpc.ServiceRegistrar) {
	s.logger.Warn("grpc server registration succeeded")
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
