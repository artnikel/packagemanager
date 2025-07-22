package archive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/artnikel/packagemanager/internal/models"
)

func TestCreateArchive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "archive_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFiles := []string{"test1.txt", "test2.go", "exclude.tmp"}
	for _, filename := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		content := "test content for " + filename
		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	targets := []models.Target{
		{Path: filepath.Join(tmpDir, "*.txt")},
		{Path: filepath.Join(tmpDir, "*.go")},
		{Path: filepath.Join(tmpDir, "*.tmp"), Exclude: "*.tmp"},
	}

	archivePath := filepath.Join(tmpDir, "test.tar.gz")

	err = CreateArchive(archivePath, targets)
	if err != nil {
		t.Fatalf("CreateArchive failed: %v", err)
	}

	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		t.Error("Archive file was not created")
	}

	files := extractAndListFiles(t, archivePath)

	expectedFiles := map[string]bool{
		filepath.Join(tmpDir, "test1.txt"): false,
		filepath.Join(tmpDir, "test2.go"):  false,
	}

	for _, file := range files {
		if _, exists := expectedFiles[file]; exists {
			expectedFiles[file] = true
		}
	}

	for file, found := range expectedFiles {
		if !found {
			t.Errorf("Expected file %s not found in archive", file)
		}
	}

	for _, file := range files {
		if filepath.Base(file) == "exclude.tmp" {
			t.Error("Excluded file found in archive")
		}
	}
}

func extractAndListFiles(t *testing.T, archivePath string) []string {
	file, err := os.Open(archivePath)
	if err != nil {
		t.Fatalf("Failed to open archive: %v", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	var files []string

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Failed to read tar header: %v", err)
		}
		files = append(files, header.Name)
	}

	return files
}

func TestExtractArchive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "extract_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	srcDir := filepath.Join(tmpDir, "src")
	destDir := filepath.Join(tmpDir, "dest")

	if err := os.MkdirAll(srcDir, 0o0750); err != nil {
		t.Fatalf("Failed to create src dir: %v", err)
	}

	testFile := filepath.Join(srcDir, "test.txt")
	testContent := "test content"
	if err := os.WriteFile(testFile, []byte(testContent), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	archivePath := filepath.Join(tmpDir, "test.tar.gz")
	targets := []models.Target{{Path: testFile}}
	if err := CreateArchive(archivePath, targets); err != nil {
		t.Fatalf("CreateArchive failed: %v", err)
	}

	if err := ExtractArchive(archivePath, destDir); err != nil {
		t.Fatalf("ExtractArchive failed: %v", err)
	}

	extractedFile := filepath.Join(destDir, testFile)
	content, err := os.ReadFile(extractedFile)
	if err != nil {
		t.Fatalf("Failed to read extracted file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Expected content '%s', got '%s'", testContent, string(content))
	}
}
