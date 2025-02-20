package sftpresume

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type FtpConfig struct {
	Server       string
	Port         int
	Username     string
	Password     string
	LocalFile    string
	RemoteFile   string
	MaxRetries   int           
	RetryDelay   time.Duration 
}


func UploadFile(cfg FtpConfig) error {
	var retries int

	for retries <= cfg.MaxRetries {
		err := upload(cfg)
		if err == nil {
			// Upload successful
			return nil
		}

		
		fmt.Printf("Upload failed: %v. Retrying (%d/%d)...\n", err, retries+1, cfg.MaxRetries)
		retries++
		time.Sleep(cfg.RetryDelay)
	}

	return fmt.Errorf("upload failed after %d attempts", cfg.MaxRetries)
}


func upload(cfg FtpConfig) error {
	// Establish an SSH connection
	sshConfig := &ssh.ClientConfig{
		User: cfg.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Server, cfg.Port), sshConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SFTP server: %w", err)
	}
	defer conn.Close()

	// Create an SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer client.Close()

	// Open the local file
	localFile, err := os.Open(cfg.LocalFile)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer localFile.Close()

	// Check the size of the remote file to determine where to resume
	var remoteSize int64
	remoteFileInfo, err := client.Stat(cfg.RemoteFile)
	if err == nil {
		remoteSize = remoteFileInfo.Size()
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat remote file: %w", err)
	}

	// Seek the local file to the point where the upload needs to resume
	_, err = localFile.Seek(remoteSize, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek local file: %w", err)
	}

	// Open the remote file in append mode
	remoteFile, err := client.OpenFile(cfg.RemoteFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %w", err)
	}
	defer remoteFile.Close()

	// Copy the remaining data
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}
