package files

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mushroomsir/mimetypes"
	"io"
	"os"
	"path/filepath"
)

type FileService struct {
	mediaPath string
}

func NewFilesSvc(mediaPath string) *FileService {
	if mediaPath == "" {
		panic(errors.New("media path can't be nil")) // <- be accurate
	}

	return &FileService{mediaPath: mediaPath}
}

func (receiver *FileService) Save(src io.Reader, contentType string) (name string, err error) {
	extensions := mimetypes.Extension(contentType)

	if len(extensions) == 0 {
		return "", errors.New("invalid extension")
	}

	uuidV4 := uuid.New().String()
	name = fmt.Sprintf("%s%s", uuidV4, "." + extensions)
	path := filepath.Join(receiver.mediaPath, name)

	dst, _ := os.Create(path)
	defer func() {
		err = dst.Close()
	}()
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}
	return name, nil
}
