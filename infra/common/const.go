package common

import "time"

const (
	ModuleCore = "core"
)

const (
	TimeLayout = "2006-01-02 15:04:05"
	DateLayout = "2006-01-02"
)

const (
	// WriteWait Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second
	// PongWait Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second
	// PingPeriod Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * 9) / 10
	// MaxMessageSize Maximum message size allowed from peer.
	MaxMessageSize = 512
)
