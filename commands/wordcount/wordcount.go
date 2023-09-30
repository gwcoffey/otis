package wordcount

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	ms2 "gwcoffey/otis/ms"
	"strings"
	"unicode/utf8"
)

const maxWidth = 40
const indentSize = "  "

type Args struct {
	ProjectPath *string `arg:"positional" help:"path to the otis project"`
	ByChapter   bool    `arg:"--chapter,-c" help:"count by chapter rather than by folder"`
}

type printBy int

const (
	byFolder printBy = iota
	byChapter
)

func sceneWordCount(scene ms2.Scene) (count int, err error) {
	text, err := scene.Text()
	if err != nil {
		return
	}
	count = len(strings.Fields(text))
	return
}

func truncate(str string) string {
	result := str
	if utf8.RuneCountInString(result) > maxWidth {
		result = result[0:maxWidth-1] + "…"
	}
	return result
}

func printManuscript(m ms2.Manuscript, by printBy) (err error) {
	count, err := ms2.WordCount(m)
	if err != nil {
		return
	}
	printLine(truncate(m.Title()), count, true)

	switch by {
	case byFolder:
		for _, folder := range m.Folders() {
			err = printFolder(folder, indentSize)
			if err != nil {
				return
			}
		}
	case byChapter:
		for _, chapter := range m.Chapters() {
			err = printChapter(chapter, indentSize)
			if err != nil {
				return
			}
		}
	}
	return
}

func printFolder(folder ms2.Folder, indent string) (err error) {
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

func folderWordCount(folder ms2.Folder) (count int, err error) {
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

func printChapter(chapter ms2.Chapter, indent string) (err error) {
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

func chapterWordCount(chapter ms2.Chapter) (count int, err error) {
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

func printScene(scene ms2.Scene, indent string) (err error) {
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

func WordCount(args *Args) {
	var manuscript ms2.Manuscript
	var err error

	if args.ProjectPath == nil {
		manuscript, err = ms2.LoadHere()
	} else {
		manuscript, err = ms2.Load(*args.ProjectPath)
	}
	if err != nil {
		panic(err)
	}

	by := byFolder
	if args.ByChapter {
		by = byChapter
	}

	err = printManuscript(manuscript, by)
	if err != nil {
		panic(err)
	}
}
