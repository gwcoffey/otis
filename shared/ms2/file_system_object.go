package ms2

type FileSystemObject interface {
	Path() string
	Number() int
	PrettyFileName() string
}
