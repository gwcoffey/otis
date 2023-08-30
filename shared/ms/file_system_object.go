package ms

type FileSystemObject interface {
	Path() string
	Number() int
	PrettyFileName() string
}
