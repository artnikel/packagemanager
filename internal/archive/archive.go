package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/artnikel/packagemanager/internal/models"
)

func CreateArchive(filename string, targets []models.Target) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	for _, target := range targets {
		matches, err := filepath.Glob(target.Path)
		if err != nil {
			return fmt.Errorf("error mask file search %s: %v", target.Path, err)
		}

		for _, match := range matches {
			if target.Exclude != "" {
				if excluded, _ := filepath.Match(target.Exclude, filepath.Base(match)); excluded {
					fmt.Printf("Skipping the file: %s (masked out %s)\n", match, target.Exclude)
					continue
				}
			}

			if err := AddFileToTar(tarWriter, match); err != nil {
				return fmt.Errorf("file addition error %s: %v", match, err)
			}
			fmt.Printf("Added file: %s\n", match)
		}
	}

	return nil
}

func AddFileToTar(tarWriter *tar.Writer, filename string) error {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return nil 
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	header := &tar.Header{
		Name: filename,
		Size: fileInfo.Size(),
		Mode: int64(fileInfo.Mode()),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tarWriter, file)
	return err
}
