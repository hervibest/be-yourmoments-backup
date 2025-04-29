package adapter

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/utils"

	"github.com/h2non/bimg"
)

type CompressAdapter interface {
	CompressImage(originalFile *multipart.FileHeader, uploadFile multipart.File, dirname string) (string, string, error)
}

type compressAdapter struct {
	compressQuality int
}

func NewCompressAdapter() CompressAdapter {
	compressQuality, _ := strconv.Atoi(utils.GetEnv("COMPRESS_QUALITY")) //75

	return &compressAdapter{
		compressQuality: compressQuality,
	}
}

func (a *compressAdapter) CompressImage(originalFile *multipart.FileHeader, uploadFile multipart.File, dirname string) (string, string, error) {
	buffer, err := io.ReadAll(uploadFile)
	if err != nil {
		log.Println("error in compress image")
		return "", "", err
	}

	filename := fmt.Sprint("Compressed_" + originalFile.Filename)

	options := bimg.Options{
		Quality: a.compressQuality,
	}

	processed, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		return filename, "", err
	}

	filePath := fmt.Sprintf("./%s/%s", dirname, filename)

	if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
		return filename, "", err
	}

	if err := bimg.Write(filePath, processed); err != nil {
		return filename, "", err
	}

	return filename, filePath, nil
}
