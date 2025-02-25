package frontend

type Selection int

const (
	Interface Selection = iota
	Peer
	Add
	Edit
	Delete
	Quit
)

type PeerKey string
