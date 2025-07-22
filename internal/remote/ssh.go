// Package remote provides SSH-based file upload and download functionality using SFTP
package remote

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/artnikel/packagemanager/internal/config"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// createSSHClient creates an SSH client
func createSSHClient(cfg *config.SSHConfig) (*ssh.Client, error) {
	var auth []ssh.AuthMethod

	if cfg.Password != "" {
		auth = append(auth, ssh.Password(cfg.Password))
	}

	if cfg.KeyPath != "" {
		key, err := os.ReadFile(cfg.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("error reading SSH key: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("error parsing SSH key: %w", err)
		}

		auth = append(auth, ssh.PublicKeys(signer))
	}

	if len(auth) == 0 {
		return nil, fmt.Errorf("authentication method (password or key) is not specified")
	}

	clientConfig := &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // #nosec G106 -- insecure host key callback used intentionally in trusted environment
	}

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	return ssh.Dial("tcp", address, clientConfig)
}

// UploadToServer uploads the file to the server via SSH
func UploadToServer(filename string, cfg *config.Config) error {
	client, err := createSSHClient(&cfg.SSH)
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			fmt.Printf("warning: failed to close ssh client: %v\n", err)
		}
	}()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer func() {
		if err := sftpClient.Close(); err != nil {
			fmt.Printf("warning: failed to close sftp client: %v\n", err)
		}
	}()

	err = sftpClient.MkdirAll(cfg.Path.RemotePackageDir)
	if err != nil {
		return fmt.Errorf("failed to create remote directory: %w", err)
	}

	localFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := localFile.Close(); err != nil {
			fmt.Printf("warning: failed to close local file: %v\n", err)
		}
	}()

	remotePath := filepath.Join(cfg.Path.RemotePackageDir, filename)

	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return fmt.Errorf("error creating a file on the server: %w", err)
	}
	defer func() {
		if err := remoteFile.Close(); err != nil {
			fmt.Printf("warning: failed to close remote file: %v\n", err)
		}
	}()

	if _, err := io.Copy(remoteFile, localFile); err != nil {
		return err
	}

	fmt.Printf("File %s uploaded to the server: %s\n", filename, remotePath)
	return nil
}

// DownloadFromServer downloads a file from the server via SSH
func DownloadFromServer(filename string, cfg *config.Config) error {
	client, err := createSSHClient(&cfg.SSH)
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			fmt.Printf("warning: failed to close ssh client: %v\n", err)
		}
	}()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer func() {
		if err := sftpClient.Close(); err != nil {
			fmt.Printf("warning: failed to close sftp client: %v\n", err)
		}
	}()

	remotePath := filepath.Join(cfg.Path.RemotePackageDir, filename)
	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := remoteFile.Close(); err != nil {
			fmt.Printf("warning: failed to close remote file: %v\n", err)
		}
	}()

	localFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := localFile.Close(); err != nil {
			fmt.Printf("warning: failed to close local file: %v\n", err)
		}
	}()

	if _, err := io.Copy(localFile, remoteFile); err != nil {
		return err
	}

	fmt.Printf("File %s downloaded from the server\n", filename)
	return nil
}
