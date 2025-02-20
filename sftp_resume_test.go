package sftpresume

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Uploader interface {
	Upload(cfg FtpConfig) error
}

type MockUploader struct {
	mock.Mock
}

func (m *MockUploader) Upload(cfg FtpConfig) error {
	args := m.Called(cfg)
	return args.Error(0)
}

func UploadFileWithUploader(cfg FtpConfig, uploader Uploader) error {
	var retries int

	for retries <= cfg.MaxRetries {
		err := uploader.Upload(cfg)
		if err == nil {
			return nil
		}

		fmt.Printf("Upload failed: %v. Retrying (%d/%d)...\n", err, retries+1, cfg.MaxRetries)
		retries++
		time.Sleep(cfg.RetryDelay)
	}

	return fmt.Errorf("upload failed after %d attempts", cfg.MaxRetries)
}

func TestUploadFile_Success(t *testing.T) {
	cfg := FtpConfig{
		Server:     "test.com",
		Port:       22,
		Username:   "user",
		Password:   "pass",
		LocalFile:  "local.txt",
		RemoteFile: "remote.txt",
		MaxRetries: 3,
		RetryDelay: time.Millisecond * 10,
	}

	mockUploader := new(MockUploader)
	mockUploader.On("Upload", cfg).Return(nil)

	err := UploadFileWithUploader(cfg, mockUploader)

	assert.NoError(t, err)
	mockUploader.AssertExpectations(t)
}

func TestUploadFile_Failure(t *testing.T) {
	cfg := FtpConfig{
		Server:     "test.com",
		Port:       22,
		Username:   "user",
		Password:   "pass",
		LocalFile:  "local.txt",
		RemoteFile: "remote.txt",
		MaxRetries: 3,
		RetryDelay: time.Millisecond * 10,
	}

	mockUploader := new(MockUploader)
	mockUploader.On("Upload", cfg).Return(errors.New("upload failed"))

	err := UploadFileWithUploader(cfg, mockUploader)

	assert.Error(t, err)
	assert.Equal(t, "upload failed after 3 attempts", err.Error())
	mockUploader.AssertExpectations(t)
}
