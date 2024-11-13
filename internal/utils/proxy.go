package utils

import "os"

// SetProxy sets the HTTP and HTTPS proxy environment variables
func SetProxy() error {
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:7890")
	os.Setenv("HTTPS_PROXY", "http://127.0.0.1:7890")
	return nil
}
