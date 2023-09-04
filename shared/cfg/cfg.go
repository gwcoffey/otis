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
	ProjectRoot string
	Author      Author `yaml:"author"`
	Address     string `yaml:"address"`
}

func FindProjectRoot() (path string, err error) {
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

	err = errors.New("this is not an otis project directory")
	return
}

func Load(projectPath string) (config Config, err error) {
	cfgYaml, err := os.ReadFile(filepath.Join(projectPath, "otis.yml"))
	if err != nil {
		return
	}

	config.ProjectRoot = projectPath

	err = yaml.Unmarshal(cfgYaml, &config)
	return
}

// TmpDir returns the path to the temporary directory for build artifacts, creating it if necessary
func (c *Config) TmpDir() (path string, err error) {
	path = filepath.Join(c.ProjectRoot, ".build")
	// make sure it exists
	err = os.MkdirAll(path, os.ModePerm)
	return
}

// DistDir returns the path to the distribution directory, creating it if necessary
func (c *Config) DistDir() (path string, err error) {
	path = filepath.Join(c.ProjectRoot, "dist")
	// make sure it exists
	err = os.MkdirAll(path, os.ModePerm)
	return
}
