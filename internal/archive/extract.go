package archive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const filePerm = 0o750

// ExtractArchive unzips the tar.gz archive
func ExtractArchive(filename, destDir string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("warning: failed to close file: %v\n", err)
		}
	}()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer func() {
		if err := gzReader.Close(); err != nil {
			fmt.Printf("warning: failed to close gzip reader: %v\n", err)
		}
	}()

	tarReader := tar.NewReader(gzReader)

	if err := os.MkdirAll(destDir, filePerm); err != nil {
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

		destPath, err := isPathSafe(destDir, header.Name)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(destPath), filePerm); err != nil {
			return err
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}

		const maxFileSize = 100 * 1024 * 1024

		limitedReader := io.LimitReader(tarReader, maxFileSize)
		if _, err := io.Copy(destFile, limitedReader); err != nil {
			_ = destFile.Close()
			return fmt.Errorf("file copy failed: %w", err)
		}

		_ = destFile.Close()

		fmt.Printf("Unzipped file: %s\n", destPath)
	}

	return nil
}

func isPathSafe(destDir, filePath string) (string, error) {
	destPath := filepath.Join(destDir, filePath)
	absDestPath, err := filepath.Abs(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute destDir: %w", err)
	}

	if !strings.HasPrefix(absDestPath, absDestDir+string(os.PathSeparator)) && absDestPath != absDestDir {
		return "", fmt.Errorf("illegal file path: %s", absDestPath)
	}

	return absDestPath, nil
}
