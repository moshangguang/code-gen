package contract

type Connection interface {
	Name() string
	Type() string
	Connect() bool
	Metadata() interface{}
}
