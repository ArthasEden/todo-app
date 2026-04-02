package core_http_server

import "time"

type Config struct {
	Addr            string `envconfig:"HTTP_ADDR" required:"true"`
	ShutdownTimeout time.Duration
}
