package adapter

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/utils"
	"github.com/oklog/ulid/v2"

	"github.com/h2non/bimg"
)

type CompressAdapter interface {
	CompressImage(originalFile *multipart.FileHeader, uploadFile multipart.File, dirname string) (string, string, error)
	CompressImageToTempFile(originalFilename string, reader io.Reader) (string, string, error)
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
		log.Println("error in reading uploaded file for compression")
		return "", "", err
	}

	options := bimg.Options{
		Quality: a.compressQuality,
	}

	processed, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		log.Printf("error in bimg.Process: %v", err)
		return "", "", err
	}

	tmpDir := filepath.Join(os.TempDir(), dirname)
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return "", "", err
	}

	ext := filepath.Ext(originalFile.Filename)
	ulidStr := ulid.Make().String()
	filename := ulidStr + ext
	filePath := filepath.Join(tmpDir, filename)

	if err := bimg.Write(filePath, processed); err != nil {
		log.Printf("error writing compressed image: %v", err)
		return filename, "", err
	}

	return filename, filePath, nil
}
func (a *compressAdapter) CompressImageToTempFile(originalFilename string, reader io.Reader) (string, string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("error reading image data for compression: %v", err)
		return "", "", err
	}

	options := bimg.Options{Quality: a.compressQuality}
	processed, err := bimg.NewImage(data).Process(options)
	if err != nil {
		log.Printf("error processing image with bimg: %v", err)
		return "", "", err
	}

	tmpDir := filepath.Join(os.TempDir(), "compressed")
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return "", "", err
	}

	// ulidStr := ulid.Make().String()
	// ext := filepath.Ext(file.Filename) // e.g., .jpg
	// cleanFilename := strings.TrimSuffix(file.Filename, ext)
	// safeFilename := strings.ReplaceAll(cleanFilename, " ", "_")

	ulidStr := ulid.Make().String()
	ext := filepath.Ext(originalFilename)
	cleanFilename := strings.TrimSuffix(originalFilename, ext)
	filename := cleanFilename + "_" + ulidStr + ext
	fullPath := filepath.Join(tmpDir, filename)

	if err := bimg.Write(fullPath, processed); err != nil {
		log.Printf("error writing compressed image to file: %v", err)
		return "", "", err
	}

	return filename, fullPath, nil
}
