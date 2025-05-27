package frontend

type Selection int

const (
	Interface Selection = iota
	Peer
	Add
	Edit
	Delete
	SwitchMode
	Quit
)

type PeerKey string
