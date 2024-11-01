package main

import (
	"github.com/DjordjeVuckovic/weather-radar/pkg/server"
	"time"
)

func main() {
	gst := server.WithGracefulShutdownTimeout(5 * time.Second)
	s := server.NewServer(":1312", gst)
	if err := s.Start(); err != nil {
		panic(err)
	}
}
