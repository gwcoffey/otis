package initcmd

import (
	_ "embed"
	"fmt"
	ms2 "gwcoffey/otis/ms"
	"gwcoffey/otis/oerr"
	"os"
)

type Args struct {
	ProjectPath *string `args:"positional"`
}

const (
	directoryPerms = 0755
	filePerms      = 0644
)

//go:embed template/otis.yml
var otisTemplate []byte

//go:embed template/00-scene.md
var sceneTemplate []byte

func Init(args *Args) {
	var existing ms2.Manuscript
	var err error

	if args.ProjectPath == nil {
		existing, err = ms2.LoadHere()
	} else {
		existing, err = ms2.Load(*args.ProjectPath)
	}
	if err != nil && !oerr.IsProjectNotFoundErr(err) {
		panic(err)
	}

	if err == nil {
		panic(fmt.Sprintf("you are already in an otis project at %s", existing.Path()))
	}

	err = os.WriteFile("otis.yml", otisTemplate, filePerms)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir("manuscript", directoryPerms)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("manuscript/00-scene.md", sceneTemplate, filePerms)
	if err != nil {
		panic(err)
	}
}
