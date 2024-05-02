package app

type Exporter interface {
	Address() string
	Export(value string) error
}
