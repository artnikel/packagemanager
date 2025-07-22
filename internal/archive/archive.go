// Package archive provides functions to create and extract tar.gz archives from specified file targets
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

// CreateArchive creates a tar.gz archive from masked files
func CreateArchive(filename string, targets []models.Target) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("warning: failed to close file: %v\n", err)
		}
	}()

	gzWriter := gzip.NewWriter(file)
	defer func() {
		if err := gzWriter.Close(); err != nil {
			fmt.Printf("warning: failed to close gzip writer: %v\n", err)
		}
	}()
	tarWriter := tar.NewWriter(gzWriter)
	defer func() {
		if err := tarWriter.Close(); err != nil {
			fmt.Printf("warning: failed to close tar writer: %v\n", err)
		}
	}()

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

			if err := addFileToTar(tarWriter, match); err != nil {
				return fmt.Errorf("file addition error %s: %v", match, err)
			}
			fmt.Printf("Added file: %s\n", match)
		}
	}

	return nil
}

// addFileToTar adds the file to a tar archive
func addFileToTar(tarWriter *tar.Writer, filename string) error {
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
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("warning: failed to close file: %v\n", err)
		}
	}()

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
