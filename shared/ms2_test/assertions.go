package ms2

import (
	"gwcoffey/otis/shared/ms2"
	"strings"
	"testing"
)

func assertWorkCount(t *testing.T, manuscript ms2.Manuscript, expected int) {
	if actual := len(manuscript.Works()); expected != actual {
		t.Fatalf("count of works = %d; expected %d", actual, expected)
	}
}

func assertWorkMetadata(t *testing.T, work ms2.Work, title string, runningTitle string, author string, surname string) {
	if expected, actual := title, work.Title(); expected != actual {
		t.Errorf("title of %s = %v; expected %v'", work, actual, expected)
	}
	if expected, actual := runningTitle, work.RunningTitle(); expected != actual {
		t.Errorf("running title of %s = %v; expected %v", work, actual, expected)
	}
	if expected, actual := author, work.Author(); expected != actual {
		t.Errorf("author of %s = %v; expected %v", work, actual, expected)
	}
	if expected, actual := surname, work.AuthorSurname(); expected != actual {
		t.Errorf("author surname of %s = %v; expected %v", work, actual, expected)
	}
}

func assertSceneCount(t *testing.T, scener ms2.Scener, expected int) {
	if actual := len(scener.Scenes()); expected != actual {
		t.Fatalf("count of scenes in %s = %v; expected %v", scener, actual, expected)
	}
}

func assertSceneMetadata(t *testing.T, scene ms2.Scene, number int, prettyName string) {
	if scene.Number() != number {
		t.Errorf("number of %s = %d; expected %d", scene, scene.Number(), number)
	}
	if scene.PrettyFileName() != prettyName {
		t.Errorf("pretty filename of %s = %v; expected %v", scene, scene.PrettyFileName(), prettyName)
	}
}

func assertSceneTextStartsWithLorem(t *testing.T, scene ms2.Scene) {
	text1, err := scene.Text()
	if err != nil {
		panic(err)
	}
	if expected, actual := "Lorem", strings.Fields(text1)[0]; expected != actual {
		t.Errorf("scene %s starts with '%v'; expected '%v'", scene.Path(), actual, expected)
	}
}

func assertChapterCount(t *testing.T, work ms2.Work, expected int) {
	if actual := len(work.Chapters()); expected != actual {
		t.Fatalf("count of chapters in work %s = %d; expected %d", work, actual, expected)
	}
}

func assertChapterMetadata(t *testing.T, chapter ms2.Chapter, title string, numbered bool) {
	if chapter.Title() != title {
		t.Errorf("chapter %s title = %s; expected %s", chapter, chapter.Title(), title)
	}

	if chapter.Numbered() != numbered {
		t.Errorf("chapter %s numbered = %v; expected %v", chapter, chapter.Numbered(), numbered)
	}
}

func assertFolderCount(t *testing.T, folderer ms2.Folderer, expected int) {
	if actual := len(folderer.Folders()); expected != actual {
		t.Fatalf("count of folders in work %s = %d; expected %d", folderer, actual, expected)
	}
}

func assertFolderMetadata(t *testing.T, folder ms2.Folder, number int, prettyName string) {
	if folder.Number() != number {
		t.Errorf("number of %s = %v; expecting %v", folder, folder.Number(), number)
	}
	if folder.PrettyFileName() != prettyName {
		t.Errorf("pretty name of %s = %v; expecting %v", folder, folder.PrettyFileName(), prettyName)
	}
}
