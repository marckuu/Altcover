package interfaces

type Producer interface {
	Produce(message []byte) error
	Close() error
}
