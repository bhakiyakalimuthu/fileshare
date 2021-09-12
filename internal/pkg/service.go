package pkg
import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
)
type Provider interface {
	GetFile(name string) (string,error)
	UploadFile(name string, data []byte)error
}
var _ Provider = (*Service)(nil)
type Service struct {
	logger *zap.Logger
	filePath string
}

func NewService(logger *zap.Logger) Provider {
	return &Service{
		logger:logger,
	}
}

func(s *Service) GetFile(name string) (string,error){
	if name == ""{
		return "", fmt.Errorf("empty file path")
	}
	return name,nil
}

func(s *Service) UploadFile(name string, data []byte) error {

	f,err := os.Open(path.Join(s.filePath,filepath.Base(name)))
	defer f.Close()
	if err!=nil {
		return fmt.Errorf("failed to open file %s",err)
	}
	for {
		_, err= f.Write(data)
		if err!=nil {
			return fmt.Errorf("failed to write data to file %s",err)
		}
	}
	return nil
}