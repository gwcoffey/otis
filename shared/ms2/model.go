package ms2

import (
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// node represents a tree of filesystem objects rooted at `manuscript/`; all ms interfaces operate on
// node internally and provide a public interface to read the directory tree as a structured manuscript
type node struct {
	path     string
	workCfg  *workMeta
	sceneCfg *sceneMeta
	children []*node
	content  []byte
}

// workMeta represents the metadata for a work, read directly from `work.yml` in a directory
// (fields are exported to support YAML unmarshalling)
type workMeta struct {
	Title         string `yaml:"title"`
	RunningTitle  string `yaml:"runningTitle"`
	Author        string `yaml:"author"`
	AuthorSurname string `yaml:"authorSurname"`
}

// sceneMeta represents the metadata of a scene, sourced from the filename
type sceneMeta struct {
	number int
}

var numberPrefixPattern = regexp.MustCompile(`^\d+`)
var metaFilenames = map[string]bool{
	"work.yml":    true,
	"chapter.yml": true,
}

func (n *node) addWorkMeta(path string) (err error) {
	content, err := os.ReadFile(filepath.Join(path, "work.yml"))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	} else if err != nil {
		return
	}

	if err = yaml.Unmarshal(content, &n.workCfg); err != nil {
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
			child, err = newSceneNode(filepath.Join(path, entry.Name()))
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

func (n *node) addSceneMeta() (err error) {
	n.sceneCfg = new(sceneMeta)
	matches := numberPrefixPattern.FindStringSubmatch(filepath.Base(n.path))
	if len(matches) == 1 {
		n.sceneCfg.number, err = strconv.Atoi(matches[0])
		if err != nil {
			return err
		}
	} else {
		err = errors.New(fmt.Sprintf("scene %s is missing required scene number prefix", n.path))
	}

	return nil
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

func newDirNode(path string) (n *node, err error) {
	n = &node{path: path}
	if err = n.addWorkMeta(path); err != nil {
		return
	}
	if err = n.addChildren(path); err != nil {
		return
	}
	return
}

func newSceneNode(path string) (n *node, err error) {
	n = &node{path: path}
	if err = n.addSceneMeta(); err != nil {
		return
	}
	return
}
