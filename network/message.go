package network

type GetStatusMessage struct{}
type StatusMessage struct {
	ID            string
	Version       uint32
	CurrentHeigth uint32
}
type GetBlockMessage struct {
	From uint32
	To   uint32
}
