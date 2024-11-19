package utils

import (
	"errors"
	"os"
)

// SetProxy 设置 HTTP 和 HTTPS 代理环境变量
func SetProxy() error {
	httpProxy := os.Getenv("HTTP_PROXY")
	httpsProxy := os.Getenv("HTTPS_PROXY")

	if httpProxy == "" && httpsProxy == "" {
		return errors.New("proxy environment variables HTTP_PROXY and HTTPS_PROXY are not set")
	}

	if httpProxy != "" {
		if err := os.Setenv("HTTP_PROXY", httpProxy); err != nil {
			return errors.New("failed to set HTTP_PROXY: " + err.Error())
		}
	}

	if httpsProxy != "" {
		if err := os.Setenv("HTTPS_PROXY", httpsProxy); err != nil {
			return errors.New("failed to set HTTPS_PROXY: " + err.Error())
		}
	}

	return nil
}
