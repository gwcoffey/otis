package ms

import (
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"gwcoffey/otis/msfs"
	"gwcoffey/otis/text"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// node represents a tree of filesystem objects rooted at `manuscript/`; all ms objects operate on
// node internally and provide a public interface to read the directory tree as a structured manuscript
type node struct {
	isDir       bool
	path        string
	chapterMeta *chapterMeta
	children    []*node
	content     []byte
	fileNumber  int
}

// chapterMeta represents the metadata for a chapter, read directly from the `chapter.yml` in the
// directory represented by the node (fields are exported to support YAML unmarshalling)
type chapterMeta struct {
	Title    string `yaml:"title"`
	Numbered *bool  `yaml:"numbered"`
}

var metaFilenames = map[string]bool{
	"chapter.yml": true,
}

type FileSystemObject interface {
	Path() string
	Number() int
	PrettyFileName() string
}

func (n *node) addChapterMeta(path string) (err error) {
	content, err := os.ReadFile(filepath.Join(path, "chapter.yml"))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	} else if err != nil {
		return
	}

	if err = yaml.Unmarshal(content, &n.chapterMeta); err != nil {
		return
	}

	return
}

func (n *node) addChildren(path string) (err error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		var child *node
		if entry.IsDir() {
			child, err = newDirNode(filepath.Join(path, entry.Name()))
		} else if filepath.Ext(entry.Name()) == ".md" {
			child, err = newFileNode(filepath.Join(path, entry.Name()))
		} else if metaFilenames[entry.Name()] || strings.HasPrefix(entry.Name(), ".") {
			// ignore these
			child = nil
		} else {
			err = errors.New(fmt.Sprintf("unexpected file in manuscript: %s", filepath.Join(path, entry.Name())))
		}
		if err != nil {
			return
		}
		if child != nil {
			n.children = append(n.children, child)
		}
	}

	return
}

func (n *node) setFileNumber() {
	var err error
	n.fileNumber, err = msfs.FileNumber(n.path)
	if err != nil {
		n.fileNumber = -1
	}
	return
}

func (n *node) loadContent() (err error) {
	if n.content == nil {
		n.content, err = os.ReadFile(n.path)
		if err != nil {
			return err
		}
	}
	return
}

var filenamePattern = regexp.MustCompile(`^\d+-?(.*)?`)

func (n *node) prettyFileName() string {
	namePart := filepath.Base(n.path)
	matches := filenamePattern.FindStringSubmatch(namePart)
	if len(matches) == 2 {
		namePart = matches[1]
	}
	return text.KebabToSentence(strings.TrimSuffix(namePart, filepath.Ext(namePart)))
}

// walk recursively walks the node tree breadth-first from this node down, calling the supplied
// function on each node (including the root)
func (n *node) walk(fn func(*node)) {
	fn(n)
	for _, child := range n.children {
		child.walk(fn)
	}
}

func (n *node) folders() (folders []Folder) {
	for _, child := range n.children {
		if child.isDir {
			folders = append(folders, &folder{node: child})
		}
	}
	return
}

func newRootNode(path string) (n *node, err error) {
	n = &node{isDir: true, path: path}
	if err = n.addChapterMeta(path); err != nil {
		return
	}
	if err = n.addChildren(path); err != nil {
		return
	}
	return
}

func newDirNode(path string) (n *node, err error) {
	if n, err = newRootNode(path); err != nil {
		return
	}
	n.setFileNumber()
	return
}

func newFileNode(path string) (n *node, err error) {
	n = &node{isDir: false, path: path}
	n.setFileNumber()
	return
}
