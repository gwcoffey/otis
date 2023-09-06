package initcmd

import (
	_ "embed"
	"fmt"
	"gwcoffey/otis/shared/o"
	"gwcoffey/otis/shared/o/oerr"
	"os"
	"path/filepath"
)

type Args struct {
	WorkCount int `arg:"-w,--works" help:"how many works to include" default:"1"`
}

const (
	directoryPerms = 0755
	filePerms      = 0644
)

//go:embed template/otis.yml
var otisTemplate []byte

//go:embed template/work.yml
var workTemplate []byte

//go:embed template/00-scene.md
var sceneTemplate []byte

func Init(projectPath string, args *Args) {
	existingRoot, err := o.FindProjectRoot()
	if err != nil && !oerr.IsProjectNotFoundErr(err) {
		panic(err)
	}

	if err == nil {
		panic(fmt.Sprintf("you are already in an otis project at %s", existingRoot))
	}

	err = os.WriteFile("otis.yml", otisTemplate, filePerms)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir("manuscript", directoryPerms)
	if err != nil {
		panic(err)
	}

	if args.WorkCount == 1 {
		err = os.WriteFile("manuscript/work.yml", workTemplate, filePerms)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("manuscript/00-scene.md", sceneTemplate, filePerms)
		if err != nil {
			panic(err)
		}
	} else {
		for i := 0; i < args.WorkCount; i++ {
			dirname := fmt.Sprintf("%02d-work-%d", i, i+1)
			err = os.Mkdir(filepath.Join("manuscript", dirname), directoryPerms)
			if err != nil {
				panic(err)
			}

			err = os.WriteFile(filepath.Join("manuscript", dirname, "work.yml"), workTemplate, filePerms)
			if err != nil {
				panic(err)
			}

			err = os.WriteFile(filepath.Join("manuscript", dirname, "00-scene.md"), sceneTemplate, filePerms)
			if err != nil {
				panic(err)
			}
		}
	}
}
