package pkg

import (
	"bhakiyakalimuthu/fileshare/config"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"go.uber.org/zap"
)

type Provider interface {
	GetFile(name string) (string, error)
	UploadFile(name string, data []byte) error
}

var _ Provider = (*Service)(nil)

type Service struct {
	logger         *zap.Logger
	SourceFilePath string
	DestFilePath   string
}

func NewService(logger *zap.Logger) Provider {
	return &Service{
		logger:         logger,
		SourceFilePath: config.Get().SourceFilePath,
		DestFilePath:   config.Get().DestFilePath,
	}
}

func (s *Service) GetFile(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("empty file path")
	}
	if !isFileExists(name) {
		return "", fmt.Errorf("file is not exist with name %s", name)
	}
	return name, nil
}
func (s *Service) UploadFile(name string, data []byte) error {

	filePath, err := s.ParseDestFilePath(name)
	if err != nil {
		return err
	}
	//s.logger.Info(fmt.Sprintf("Service - UploadFile called for file %s", absPath))
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("failed to open file %s", err)
	}
	count := 0
	n := len(data)
	for {
		w, err := f.Write(data[count:])
		if err != nil {
			return fmt.Errorf("failed to write data to file %s", err)
		}
		count += w
		if w >= n {
			s.logger.Info(fmt.Sprintf("Service - UploadFile completes writing to a file byte count %d", count))
			return nil
		}
	}
	return nil
}

func (s *Service) ParseDestFilePath(fName string) (string, error) {
	return validateFilePath(fName, s.DestFilePath)

}

func validateFilePath(fName, fPath string) (string, error) {

	absPath, err := filepath.Abs(fPath)
	if err != nil {
		return "", fmt.Errorf("failed to get abs file path %s", err)
	}
	return path.Join(absPath, filepath.Base(fName)), nil
}

func isFileExists(filepath string) bool {

	fileInfo, err := os.Stat(filepath)

	if os.IsNotExist(err) {
		return false
	}
	// Return false if the fileInfo says the file path is a directory.
	return !fileInfo.IsDir()
}
