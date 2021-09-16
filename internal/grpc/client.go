package grpc

import (
	"bhakiyakalimuthu/fileshare/internal/pkg"
	pb "bhakiyakalimuthu/fileshare/proto"
	"context"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
)

type Client struct {
	logger *zap.Logger
	server pb.UploaderClient
	svc    pkg.Provider
}

func NewClient(logger *zap.Logger, server pb.UploaderClient) *Client {
	return &Client{
		logger: logger,
		server: server,
	}
}

func (c *Client) Upload(ctx context.Context, file string) error {
	c.logger.Info(fmt.Sprintf("client upload started for file %s", file))

	// open a file
	f, err := os.Open(file)
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to open a file %s : %s", file, err))
		return fmt.Errorf("failed to open a file %s", err)
	}
	defer f.Close()
	// create uploader
	stream, err := c.server.Upload(ctx)
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to create stream uploader %s", err))
		return fmt.Errorf("failed to create stream uploader %s", err)
	}
	defer stream.CloseSend()
	// stream contents
	buf := make([]byte, 64*1024)
	counter := 0
	for {
		counter += 1
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				c.logger.Error(fmt.Sprintf("End of file %s", err))
				break
			}
			c.logger.Error(fmt.Sprintf("failed to read bytes from file %s : %s", file, err))
			return fmt.Errorf("failed to read bytes from file %s : %s", file, err)
		}
		c.logger.Error(fmt.Sprintf("bytes read from files %d : counter : %d", n, counter))
		if err := stream.Send(&pb.UploadRequest{
			Name:    file,
			Content: buf[:n],
		}); err != nil {
			c.logger.Error(fmt.Sprintf("failed to send file buffer chunks%s", err))
			return fmt.Errorf("failed to send file buffer chunks%s", err)
		}
	}
	out, err := stream.CloseAndRecv()
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to receive upstream status %s", err))
		return fmt.Errorf("failed to receive upstream status %s", err)
	}
	if out.StatusCode != pb.StatusCode_OK {
		c.logger.Error(fmt.Sprintf("failed to upload statuscode %s : %s : %s", out.StatusCode, out.Status, err))
		return fmt.Errorf("failed to upload statuscode %s : %s : %s", out.StatusCode, out.Status, err)
	}
	c.logger.Error(fmt.Sprintf("upload succeeded   statusCode - %s : status - %s", out.StatusCode, out.Status))
	return nil
}
