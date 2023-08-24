// Package ms provides primitives for working with the manuscript content
package ms

import (
	"errors"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Type node represents a filesystem node within the `manuscript` directory.
// It conforms to both Dir and Scene interfaces
type node struct {
	// the path to the filesystem object
	Path string

	// the parent node
	ParentNode *node

	// all child nodes
	ChildNodes []*node

	// the node type
	Type NodeType

	// the content of a file this node represents (if loaded)
	Content *[]byte
}

type NodeType int

const (
	NodeTypeDir NodeType = iota
	NodeTypeScene
)

const markdownExt = ".md"
const kebabSeparator = "-"

type INode interface {
	// Parent Function returns the parent directory
	Parent() Dir

	// Filename returns the name (including extension) of the filesystem object
	// that underlies this node
	Filename() string

	// Name returns the pretty name by converting snake-case to camelcase
	Name() string
}

type Dir interface {
	INode

	// SubDirs function returns the direct child directories
	SubDirs() []Dir

	// SubDir returns the direct child directory with the given name
	SubDir(name string) *Dir

	// Scenes returns the scenes within this directory
	Scenes() []Scene

	// AllScenes returns all scenes within this directory (including subdirectories)
	AllScenes() []Scene
}

type Scene interface {
	INode

	// Text returns the textual content (in raw Markdown) of the scene
	// Calling this causes the scene to be loaded if it isn't already, and the
	// error returned is any error from the reader
	Text() (*string, error)

	// Number returns the scene number, which is the prefix of the filename
	Number() int
}

func readDir(path string) *node {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	var children []*node
	var thisNode = node{Path: path, ChildNodes: children, Type: NodeTypeDir}

	for _, entry := range entries {
		subPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			nextNode := readDir(subPath)
			children = append(children, nextNode)
		} else if filepath.Ext(subPath) == markdownExt {
			nextNode := node{Path: subPath, ParentNode: &thisNode, ChildNodes: []*node{}, Type: NodeTypeScene}
			children = append(children, &nextNode)
		}
	}

	thisNode.ChildNodes = children

	return &thisNode
}

func (n *node) DumpToString() string {
	sb := strings.Builder{}
	switch n.Type {
	case NodeTypeDir:
		sb.WriteString("node(dir)@")
	case NodeTypeScene:
		sb.WriteString("node(scene)@")
	default:
		panic(fmt.Sprintf("unexpected node type: %s", n.Type))
	}
	sb.WriteString(n.Filename())
	if len(n.ChildNodes) > 0 {
		sb.WriteString("{")
		for idx, child := range n.ChildNodes {
			if idx > 0 {
				sb.WriteRune(',')
			}
			sb.WriteString(child.DumpToString())
		}
		sb.WriteString("}")
	}

	return sb.String()
}

// BEGIN INode interface

func (n *node) Parent() Dir {
	return n.ParentNode
}

func (n *node) Filename() string {
	return filepath.Base(n.Path)
}

func (n *node) Name() string {
	switch n.Type {
	case NodeTypeDir:
		return kebabToSentence(n.Filename())
	case NodeTypeScene:
		namePart := n.Filename()
		re := regexp.MustCompile(`^\d+-?(.*).md`)
		matches := re.FindStringSubmatch(namePart)
		if len(matches) == 2 {
			namePart = matches[1]
		}
		return fmt.Sprintf("%02d. %s", n.Number(), kebabToSentence(namePart))
	}

	panic(errors.New(fmt.Sprintf("No such type: %s", n.Type)))
}

func kebabToSentence(str string) string {
	words := strings.Split(str, kebabSeparator)
	words[0] = cases.Title(language.English).String(words[0])
	return strings.Join(words, " ")
}

// BEGIN Dir interface

func (n *node) Scenes() []Scene {
	var result []Scene
	for _, child := range n.ChildNodes {
		if child.Type == NodeTypeScene {
			result = append(result, child)
		}
	}
	return result
}

func (n *node) AllScenes() []Scene {
	var result []Scene
	for _, child := range n.ChildNodes {
		if child.Type == NodeTypeScene {
			result = append(result, child)
		} else if child.Type == NodeTypeDir {
			result = append(result, child.AllScenes()...)
		}
	}
	return result
}

func (n *node) SubDirs() []Dir {
	var result []Dir
	for _, child := range n.ChildNodes {
		if child.Type == NodeTypeDir {
			result = append(result, child)
		}
	}
	return result
}

func (n *node) SubDir(name string) *Dir {
	for _, dir := range n.SubDirs() {
		if dir.Filename() == name {
			return &dir
		}
	}

	return nil
}

// BEGIN Scene interface

func (n *node) load() error {
	if n.Content == nil {
		bytes, err := os.ReadFile(n.Path)
		if err != nil {
			return err
		}
		n.Content = &bytes
	}
	return nil
}

func (n *node) Text() (*string, error) {
	err := n.load()
	if err != nil {
		return nil, err
	}
	text := string(*n.Content)
	return &text, nil
}

func (n *node) Number() int {
	re := regexp.MustCompile(`^\d+`)
	matches := re.FindStringSubmatch(n.Filename())
	if len(matches) == 1 {
		num, err := strconv.Atoi(matches[0])
		if err != nil {
			panic(err)
		}
		return num
	}
	return 0
}

// BEGIN public interface

func Load() Dir {
	x := readDir("manuscript/")
	return x
}
