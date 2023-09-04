package wordcount

import (
	"errors"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gwcoffey/otis/shared/cfg"
	"gwcoffey/otis/shared/ms"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const maxWidth = 40
const indentSize = "  "

type Args struct {
	Work      *string `arg:"positional" help:"count only the specified work in a multi-work manuscript"`
	ByChapter bool    `arg:"--chapter" help:"count by chapter rather than by folder"`
}

type printBy int

const (
	byFolder printBy = iota
	byChapter
)

func sceneWordCount(scene ms.Scene) (count int, err error) {
	text, err := scene.Text()
	if err != nil {
		return
	}
	count = len(strings.Fields(text))
	return
}

func workWordCount(work ms.Work) (count int, err error) {
	for _, scene := range work.AllScenes() {
		var scount int
		scount, err = sceneWordCount(scene)
		if err != nil {
			return
		}
		count += scount
	}
	return
}

func truncate(str string) string {
	result := str
	if utf8.RuneCountInString(result) > maxWidth {
		result = result[0:maxWidth-1] + "…"
	}
	return result
}

func printWork(work ms.Work, by printBy) (err error) {
	count, err := workWordCount(work)
	if err != nil {
		return
	}
	printLine(truncate(work.Title()), count, true)

	for _, scene := range work.Scenes() {
		err = printScene(scene, indentSize)
		if err != nil {
			return
		}
	}

	switch by {
	case byFolder:
		for _, folder := range work.Folders() {
			err = printFolder(folder, indentSize)
			if err != nil {
				return
			}
		}
	case byChapter:
		for _, chapter := range work.Chapters() {
			err = printChapter(chapter, indentSize)
			if err != nil {
				return
			}
		}
	}
	return
}

func printFolder(folder ms.Folder, indent string) (err error) {
	fcount, err := folderWordCount(folder)
	if err != nil {
		return
	}

	label := fmt.Sprintf("%02d. %s", folder.Number()+1, folder.PrettyFileName())

	printLine(truncate(indent+label), fcount, true)

	for _, scene := range folder.Scenes() {
		err = printScene(scene, indent+indentSize)
		if err != nil {
			return
		}
	}

	for _, child := range folder.Folders() {
		err = printFolder(child, indent+indentSize)
		if err != nil {
			return
		}
	}

	return
}

func folderWordCount(folder ms.Folder) (count int, err error) {
	for _, scene := range folder.AllScenes() {
		var scount int
		scount, err = sceneWordCount(scene)
		if err != nil {
			return
		}
		count += scount
	}
	return
}

func printChapter(chapter ms.Chapter, indent string) (err error) {
	ccount, err := chapterWordCount(chapter)
	if err != nil {
		return
	}

	var label string
	if chapter.Number() != nil {
		label = fmt.Sprintf("% 2d. %s", *chapter.Number(), chapter.Title())
	} else {
		label = fmt.Sprintf("    %s", chapter.Title())
	}

	printLine(truncate(indent+label), ccount, false)
	return
}

func chapterWordCount(chapter ms.Chapter) (count int, err error) {
	var scount int
	for _, scene := range chapter.Scenes() {
		scount, err = sceneWordCount(scene)
		if err != nil {
			return
		}
		count += scount
	}
	return
}

func printScene(scene ms.Scene, indent string) (err error) {
	scount, err := sceneWordCount(scene)
	if err != nil {
		return
	}
	label := fmt.Sprintf("%02d. %s", scene.Number()+1, scene.PrettyFileName())
	printLine(truncate(indent+label), scount, false)
	return
}

func printLine(label string, count int, emphasize bool) {
	out := message.NewPrinter(language.English)
	format := fmt.Sprintf("%%-%d.%ds : %%7d\n", maxWidth, maxWidth)
	if emphasize {
		format = fmt.Sprintf("\033[94m%s\033[0m", format)
	}
	_, err := out.Printf(format, label, count)
	if err != nil {
		panic(err)
	}
}

func selectWorks(args *Args, manuscript ms.Manuscript) (works []ms.Work, err error) {
	works = manuscript.Works()

	if args.Work != nil {
		for _, work := range works {
			if filepath.Base(work.Path()) == *args.Work {
				works = []ms.Work{work}
				return
			}
		}
		err = errors.New(fmt.Sprintf("no such work: %s", *args.Work))
	}

	return
}

func WordCount(config cfg.Config, args *Args) {
	var manuscript ms.Manuscript
	var err error

	manuscript, err = ms.Load(filepath.Join(config.ProjectRoot, "manuscript"))
	if err != nil {
		panic(err)
	}

	works, err := selectWorks(args, manuscript)
	if err != nil {
		panic(err)
	}

	by := byFolder
	if args.ByChapter {
		by = byChapter
	}

	for _, work := range works {
		err := printWork(work, by)
		if err != nil {
			panic(err)
		}
	}
}
