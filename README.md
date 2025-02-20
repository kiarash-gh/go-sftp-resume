# sftpresume

`sftpresume` is a Go package that enables SFTP file uploads with resume capability. If an upload is interrupted, it will resume from where it left off, avoiding redundant data transfer.

## Features
- Secure SFTP file uploads
- Resume incomplete file uploads
- Configurable retries and delay handling

## Installation

```sh
go get github.com/yourusername/sftpresume
```

## Usage

### Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/yourusername/sftpresume"
)

func main() {
	config := sftpresume.FtpConfig{
		Server:     "sftp.example.com",
		Port:       22,
		Username:   "user",
		Password:   "password",
		LocalFile:  "./localfile.txt",
		RemoteFile: "/remote/path/remote.txt",
		MaxRetries: 3,
		RetryDelay: 5 * time.Second,
	}

	err := sftpresume.UploadFile(config)
	if err != nil {
		fmt.Println("Upload failed:", err)
	} else {
		fmt.Println("Upload successful")
	}
}
```

## Configuration

| Field       | Type          | Description |
|------------|--------------|-------------|
| `Server`   | `string`     | SFTP server address |
| `Port`     | `int`        | SFTP port (default 22) |
| `Username` | `string`     | SFTP username |
| `Password` | `string`     | SFTP password |
| `LocalFile` | `string`    | Path to the local file |
| `RemoteFile` | `string`  | Path to the remote file |
| `MaxRetries` | `int`      | Maximum retry attempts |
| `RetryDelay` | `time.Duration` | Delay between retries |

## License

This project is licensed under the MIT License.

## Contributions

Feel free to open issues or submit pull requests to enhance the functionality!

