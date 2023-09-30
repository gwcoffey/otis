package ms

import (
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gwcoffey/otis/oerr"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func validateManuscript(m Manuscript) (err error) {
	if len(m.Chapters()) > 0 && m.Scenes()[0].Path() != m.Chapters()[0].Scenes()[0].Path() {
		err = errors.New(fmt.Sprintf("manuscript %s has scenes before the first chapter", m.Path()))
	}

	return
}

func MustLoad(path string) Manuscript {
	manuscript, err := Load(path)
	if err != nil {
		panic(err)
	}
	return manuscript
}

// Load loads the manuscript at the given path
func Load(path string) (ms Manuscript, err error) {
	yamlData, err := os.ReadFile(filepath.Join(path, "otis.yml"))
	if err != nil {
		return
	}

	var meta manuscriptMeta
	err = yaml.Unmarshal(yamlData, &meta)
	if err != nil {
		return
	}

	node, err := newRootNode(filepath.Join(path, "manuscript"))
	if err != nil {
		return
	}

	ms = &manuscript{path: path, meta: meta, node: node}
	if err = validateManuscript(ms); err != nil {
		return
	}
	return
}

// LoadHere loads the manuscript that contains the current working directory
func LoadHere() (Manuscript, error) {
	path, err := findProjectRoot()
	if err != nil {
		return nil, err
	}

	return Load(path)
}

func findProjectRoot() (path string, err error) {
	path, err = os.Getwd()
	if err != nil {
		return
	}

	for ; path != filepath.Dir(path); path = filepath.Dir(path) {
		_, err = os.Stat(filepath.Join(path, "otis.yml"))
		if err == nil || !os.IsNotExist(err) {
			return
		}
	}

	err = oerr.ProjectNotFound()
	return
}

func FindNextSceneNumber(folder Folder) int {
	sceneNumber := 0
	for _, scene := range folder.Scenes() {
		if sceneNumber <= scene.Number() {
			sceneNumber = scene.Number() + 1
		}
	}
	return sceneNumber
}

var nameReplaceRegex = regexp.MustCompile(`^\d\d`)

func MakeRoomForScene(scenes []Scene, sceneNumber int) (err error) {
	// iterate the scenes and move the specified scene number forward one spot
	// (and recursively make room for that move if needed)
	for i, scene := range scenes {
		if scene.Number() == sceneNumber {
			if len(scenes) > 1 {
				// recursively make room for the move we're about to make
				err = MakeRoomForScene(scenes[i+1:], scene.Number()+1)
				if err != nil {
					return
				}
			}

			// determine new filename by incrementing scene number prefix on existing name
			newName := nameReplaceRegex.ReplaceAllString(filepath.Base(scene.Path()), fmt.Sprintf("%02d", sceneNumber+1))

			// move the file
			err = os.Rename(scene.Path(), filepath.Join(filepath.Dir(scene.Path()), newName))
			if err != nil {
				return
			}
		}
	}
	return
}

func WordCount(m Manuscript) (count int, err error) {
	for _, scene := range m.Scenes() {
		var text string
		text, err = scene.Text()
		if err != nil {
			return
		}
		count += len(strings.Fields(text))
	}
	return
}

func ApproximateWordCount(m Manuscript) (result string, err error) {
	count, err := WordCount(m)
	if err != nil {
		return
	}

	count = int(math.Round(float64(count)/500.0)) * 500

	p := message.NewPrinter(language.English)
	result = p.Sprintf("%d", count)
	return
}
