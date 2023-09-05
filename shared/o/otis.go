package o

import (
	"errors"
	"github.com/go-yaml/yaml"
	"gwcoffey/otis/shared/ms"
	"os"
	"path/filepath"
)

type authorMeta struct {
	Name     string  `yaml:"name"`
	Surname  *string `yaml:"surname"`
	RealName *string `yaml:"realname"`
}

type otisMeta struct {
	Author       authorMeta `yaml:"author"`
	AddressLines string     `yaml:"address"`
}

type otis struct {
	root       string
	otisMeta   otisMeta
	manuscript ms.Manuscript
}

type Otis interface {
	ProjectRoot() string
	TmpDir() (string, error)
	DistDir() (string, error)
	AuthorName() string
	AuthorSurname() *string
	AuthorRealName() *string
	Address() string
	Manuscript() (ms.Manuscript, error)
}

func FindProjectRoot() (path string, err error) {
	path, err = os.Getwd()
	if err != nil {
		return
	}

	for ; path != filepath.Dir(path); path = filepath.Dir(path) {
		_, err = os.Stat(filepath.Join(path, "o.yml"))
		if err == nil || !os.IsNotExist(err) {
			return
		}
	}

	err = errors.New("this is not an o project directory")
	return
}

func Load(projectPath string) (_ Otis, err error) {
	otisYml, err := os.ReadFile(filepath.Join(projectPath, "o.yml"))
	if err != nil {
		return
	}

	var result otis
	result.root = projectPath

	err = yaml.Unmarshal(otisYml, &result.otisMeta)
	if err != nil {
		return
	}
	return &result, nil
}

// ProjectRoot returns the path to the project directory
func (o *otis) ProjectRoot() string {
	return o.root
}

// TmpDir returns the path to the temporary directory for build artifacts, creating it if necessary
func (o *otis) TmpDir() (path string, err error) {
	path = filepath.Join(o.root, ".build")
	// make sure it exists
	err = os.MkdirAll(path, os.ModePerm)
	return
}

// DistDir returns the path to the distribution directory, creating it if necessary
func (o *otis) DistDir() (path string, err error) {
	path = filepath.Join(o.root, "dist")
	// make sure it exists
	err = os.MkdirAll(path, os.ModePerm)
	return
}

// AuthorName returns the project-level configured author name
func (o *otis) AuthorName() string {
	return o.otisMeta.Author.Name
}

// AuthorSurname returns the project-level configured author surname
func (o *otis) AuthorSurname() *string {
	return o.otisMeta.Author.Surname
}

// AuthorRealName returns the project-level configured author real name
func (o *otis) AuthorRealName() *string {
	return o.otisMeta.Author.RealName
}

// Address returns the project-level configured address
func (o *otis) Address() string {
	return o.otisMeta.AddressLines
}

// Manuscript returns the model representing the manuscript in this project
func (o *otis) Manuscript() (_ ms.Manuscript, err error) {
	if o.manuscript == nil {
		o.manuscript, err = ms.Load(filepath.Join(o.root, "manuscript"))
		if err != nil {
			return
		}
	}
	return o.manuscript, nil
}
