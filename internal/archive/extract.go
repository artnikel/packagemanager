package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ExtractArchive(filename, destDir string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, header.Name)

		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}

		if _, err := io.Copy(destFile, tarReader); err != nil {
			destFile.Close()
			return err
		}
		destFile.Close()

		fmt.Printf("Unzipped file: %s\n", destPath)
	}

	return nil
}
