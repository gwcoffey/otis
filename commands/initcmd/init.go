package initcmd

import (
	_ "embed"
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

func Init(args *Args) (err error) {
	var existing ms2.Manuscript

	if args.ProjectPath == nil {
		existing, err = ms2.LoadHere()
	} else {
		existing, err = ms2.Load(*args.ProjectPath)
	}
	if err != nil && !oerr.IsProjectNotFoundErr(err) {
		return
	}

	if err == nil {
		return oerr.AlreadyAProject(existing.Path())
	}

	err = os.WriteFile("otis.yml", otisTemplate, filePerms)
	if err != nil {
		return
	}

	err = os.Mkdir("manuscript", directoryPerms)
	if err != nil {
		return
	}

	err = os.WriteFile("manuscript/00-scene.md", sceneTemplate, filePerms)
	if err != nil {
		return
	}

	return nil
}
