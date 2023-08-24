package ms2_test

import (
	"gwcoffey/otis/shared/ms2"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadErrors(t *testing.T) {
	withExemplar("no-manuscript", func() {
		_, err := ms2.Load()
		if err == nil {
			t.Error("Load() error = nil; expected non-nil")
		}
	})

	withExemplar("flat", func() {
		_, err := ms2.Load()
		if err != nil {
			t.Errorf("Load() error = '%v'; expected nil", err)
		}
	})
}

func TestMustLoadErrors(t *testing.T) {
	withExemplar("no-manuscript", func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Error("Load() error = nil; expected non-nil")
			}
		}()
		ms2.MustLoad()
	})

	withExemplar("flat", func() {
		defer func() {
			err := recover()
			if err != nil {
				t.Errorf("Load() error = '%v'; expected nil", err)
			}
		}()
		ms2.MustLoad()
	})
}

func TestMustLoadFlat(t *testing.T) {
	withExemplar("flat", func() {
		manuscript := ms2.MustLoad()
		assertWorkCount(t, manuscript, 1)
		assertWorkMetadata(t, manuscript.Works()[0], "Flat Example", "Flat", "Geoff Coffey", "Coffey")
		assertWorkSceneCount(t, manuscript.Works()[0], 2)
		assertSceneMetadata(t, manuscript.Works()[0].Scenes()[0], 0)
		assertSceneMetadata(t, manuscript.Works()[0].Scenes()[1], 1)
		assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Scenes()[0])
		assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Scenes()[1])
		assertChapterCount(t, manuscript.Works()[0], 0)
		assertFolderCount(t, manuscript.Works()[0], 0)
	})
}

func assertWorkCount(t *testing.T, manuscript ms2.Manuscript, expected int) {
	if actual := len(manuscript.Works()); expected != actual {
		t.Fatalf("count of works = %d; expected %d", actual, expected)
	}
}

func assertWorkMetadata(t *testing.T, work ms2.Work, title string, runningTitle string, author string, surname string) {
	if expected, actual := title, work.Title(); expected != actual {
		t.Errorf("it.Works()[0].Title() = '%v'; expected '%v''", work.Title(), expected)
	}
	if expected, actual := runningTitle, work.RunningTitle(); expected != actual {
		t.Errorf("it.Works()[0].RunningTitle() = '%v'; expected '%v''", actual, expected)
	}
	if expected, actual := author, work.Author(); expected != actual {
		t.Errorf("it.Works()[0].Author() = '%v'; expected '%v''", actual, expected)
	}
	if expected, actual := surname, work.AuthorSurname(); expected != actual {
		t.Errorf("it.Works()[0].AuthorSurname() = '%v'; expected '%v''", actual, expected)
	}
}

func assertWorkSceneCount(t *testing.T, work ms2.Work, expected int) {
	if actual := len(work.Scenes()); expected != actual {
		t.Fatalf("count of scenes in work %s = %d; expected %d", work, actual, expected)
	}
}

func assertSceneMetadata(t *testing.T, scene ms2.Scene, number int) {
	if scene.Number() != number {
		t.Errorf("number of %s = %d; expected %d", scene, scene.Number(), number)
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

func assertFolderCount(t *testing.T, work ms2.Work, expected int) {
	if actual := len(work.Folders()); expected != actual {
		t.Fatalf("count of folders in work %s = %d; expected %d", work, actual, expected)
	}
}

func mustChdir(path string) {
	err := os.Chdir(path)
	if err != nil {
		panic(err)
	}
}

func withExemplar(name string, runner func()) {
	// record the working dir
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// defer changing back when we're done
	defer mustChdir(wd)

	// change to the test dir
	mustChdir(filepath.Join("../../test-data/manuscripts", name))

	// run the test
	runner()
}
