package client

import "time"

type Options struct {
	TCPAddress string
	UDPAddress string
	MaxHeartbeatInterval time.Duration
}

func NewOptions() *Options {

	return &Options{
		TCPAddress: "0.0.0.0:6666",
		UDPAddress: "0.0.0.0:5555",
		MaxHeartbeatInterval: 60 * time.Second,
	}
}
