package cfg

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"os"
	"path/filepath"
)

type Title struct {
	Full    string  `yaml:"full"`
	Running *string `yaml:"running"`
}

type Author struct {
	Name     string  `yaml:"name"`
	Surname  *string `yaml:"surname"`
	Realname *string `yaml:"realname"`
}

type ConfigData struct {
	Title   Title  `yaml:"title"`
	Author  Author `yaml:"author"`
	Address string `yaml:"address"`
}

type Config interface {
	FullTitle() string
	RunningTitle() *string
	AuthorName() string
	AuthorSurname() *string
	AuthorRealName() *string
	AddressLines() string
}

func findAndGet() ([]byte, error) {
	var err error

	dir, err := os.Getwd() // Get the current working directory
	if err != nil {
		return nil, err
	}

	for ; dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		path := filepath.Join(dir, "otis.yml")
		content, readErr := os.ReadFile(path)
		if readErr == nil {
			return content, nil
		} else if !os.IsNotExist(readErr) {
			return nil, readErr
		}
	}

	return nil, fmt.Errorf("otis.yml not found")
}

func (c ConfigData) FullTitle() string {
	return c.Title.Full
}

func (c ConfigData) RunningTitle() *string {
	return c.Title.Running
}

func (c ConfigData) AuthorName() string {
	return c.Author.Name
}

func (c ConfigData) AuthorSurname() *string {
	return c.Author.Surname
}

func (c ConfigData) AuthorRealName() *string {
	return c.Author.Realname
}

func (c ConfigData) AddressLines() string {
	return c.Address
}

func FindAndLoad() Config {
	var cfg ConfigData
	cfgYaml, err := findAndGet()
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(cfgYaml, &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
