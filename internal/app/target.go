package app

type Target interface {
	Export(value string) error
}
