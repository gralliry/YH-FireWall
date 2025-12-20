package connection

type Direction int

const (
	Inbound  Direction = 0
	Outbound Direction = 1
	Forward  Direction = 2
	Unknown  Direction = 3
)
