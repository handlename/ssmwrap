package app

type Exporter interface {
	Export(value string) error
}
