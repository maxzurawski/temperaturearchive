package handlers

type SocketState int

const (
	RECEIVING SocketState = iota
	CLOSED
)
