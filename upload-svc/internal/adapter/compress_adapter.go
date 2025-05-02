package adapter

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hervibest/be-yourmoments-backup/upload-svc/internal/helper/utils"
	"github.com/oklog/ulid/v2"

	"github.com/h2non/bimg"
)

type CompressAdapter interface {
	CompressImage(originalFile *multipart.FileHeader, uploadFile multipart.File, dirname string) (string, string, error)
	CompressImageToTempFile(originalFilename string, reader io.Reader) (string, error)
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
	// Streaming → buffer baca sebagian demi sebagian
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

	// Buat file sementara di dir yang aman
	tmpDir := filepath.Join(os.TempDir(), dirname)
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return "", "", err
	}

	// Gunakan nama unik berbasis ulid
	filename := "Compressed_" + ulid.Make().String() + "_" + originalFile.Filename
	filePath := filepath.Join(tmpDir, filename)

	if err := bimg.Write(filePath, processed); err != nil {
		log.Printf("error writing compressed image: %v", err)
		return filename, "", err
	}

	return filename, filePath, nil
}

func (a *compressAdapter) CompressImageToTempFile(originalFilename string, reader io.Reader) (string, error) {
	// Baca seluruh konten dari reader (karena bimg perlu []byte)
	data, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("error reading image data for compression: %v", err)
		return "", err
	}

	// Kompres menggunakan bimg
	options := bimg.Options{
		Quality: a.compressQuality,
	}

	processed, err := bimg.NewImage(data).Process(options)
	if err != nil {
		log.Printf("error processing image with bimg: %v", err)
		return "", err
	}

	// Simpan ke file temporer
	tmpDir := filepath.Join(os.TempDir(), "compressed")
	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return "", err
	}

	// Buat nama file unik
	// Penamaan best practice
	filename := "compressed_" + ulid.Make().String() + filepath.Ext(originalFilename) // .jpg, .png
	fullPath := filepath.Join(tmpDir, filename)

	if err := bimg.Write(fullPath, processed); err != nil {
		log.Printf("error writing compressed image to file: %v", err)
		return "", err
	}

	return fullPath, nil
}
