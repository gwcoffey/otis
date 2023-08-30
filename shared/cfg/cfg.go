package cfg

import (
	"errors"
	"github.com/go-yaml/yaml"
	"os"
	"path/filepath"
)

type Author struct {
	Name     string  `yaml:"name"`
	Surname  *string `yaml:"surname"`
	RealName *string `yaml:"realname"`
}

type Config struct {
	Author  Author `yaml:"author"`
	Address string `yaml:"address"`
}

func ProjectPath() (path string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}

	for ; dir != filepath.Dir(dir); dir = filepath.Dir(dir) {
		path := filepath.Join(dir, "otis.yml")
		_, err = os.Stat(path)
		if err == nil || !os.IsNotExist(err) {
			return
		}
	}

	err = errors.New("this is not an otis project directory")
	return
}

func FindAndLoad() (cfg Config, err error) {
	projectPath, err := ProjectPath()
	if err != nil {
		return
	}

	cfgYaml, err := os.ReadFile(filepath.Join(projectPath, "otis.yml"))
	if err != nil {
		return
	}

	err = yaml.Unmarshal(cfgYaml, &cfg)
	return
}
