package session

type State string

const (
	StateOnline       State = "online"
	StateOffline      State = "offline"
	StateReconnecting State = "reconnecting"
)
