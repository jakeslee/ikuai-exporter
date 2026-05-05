package options

import "github.com/sirupsen/logrus"

type ServerOptions struct {
	URL          string `json:"url"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Level        string `json:"level"`
	InsecureSkip bool   `json:"insecureSkip"`
}

func NewServerOptions() *ServerOptions {
	return &ServerOptions{
		URL:          "http://10.0.1.253",
		Username:     "test",
		Password:     "test123",
		Level:        logrus.InfoLevel.String(),
		InsecureSkip: true,
	}
}
