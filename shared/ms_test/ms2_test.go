package ms

import (
	"gwcoffey/otis/shared/ms"
	"strings"
	"testing"
)

func TestLoadErrors(t *testing.T) {
	_, err := ms.Load("../../test-data/manuscripts/does-not-exist")
	if err == nil {
		t.Error("Load(...) error = nil; expected non-nil")
	}

	_, err = ms.Load("../../test-data/manuscripts/flat")
	if err != nil {
		t.Errorf("Load(...) error = '%v'; expected nil", err)
	}
}

func TestMustLoadErrors(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Error("MustLoad(...) error = nil; expected non-nil")
		}
	}()
	ms.MustLoad("../../test-data/manuscripts/does-not-exist")

	defer func() {
		err := recover()
		if err != nil {
			t.Errorf("MustLoad(...) error = '%v'; expected nil", err)
		}
	}()
	ms.MustLoad("../../test-data/manuscripts/flat")
}

func TestMustLoadFlat(t *testing.T) {
	manuscript := ms.MustLoad("../../test-data/manuscripts/flat")
	assertWorkCount(t, manuscript, 1)
	assertWorkMetadata(t, manuscript.Works()[0], "Flat Example", "Flat", "Geoff Coffey", "Coffey")
	assertSceneCount(t, manuscript.Works()[0], 2)
	assertSceneMetadata(t, manuscript.Works()[0].Scenes()[0], 0, "Scene")
	assertSceneMetadata(t, manuscript.Works()[0].Scenes()[1], 1, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Scenes()[0])
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Scenes()[1])
	assertChapterCount(t, manuscript.Works()[0], 0)
	assertFolderCount(t, manuscript.Works()[0], 0)
}

func TestMustLoadChapters(t *testing.T) {
	manuscript := ms.MustLoad("../../test-data/manuscripts/chapters")
	assertWorkCount(t, manuscript, 1)
	assertWorkMetadata(t, manuscript.Works()[0], "Chapters Example", "Chapters", "Geoff Coffey", "Coffey")
	assertSceneCount(t, manuscript.Works()[0], 0)
	assertChapterCount(t, manuscript.Works()[0], 2)
	assertChapterMetadata(t, manuscript.Works()[0].Chapters()[0], "My Epilogue", 0)
	assertChapterMetadata(t, manuscript.Works()[0].Chapters()[1], "My Chapter", 1)
	assertSceneCount(t, manuscript.Works()[0].Chapters()[0], 1)
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[0].Scenes()[0], 0, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Chapters()[0].Scenes()[0])
	assertSceneCount(t, manuscript.Works()[0].Chapters()[1], 1)
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[1].Scenes()[0], 0, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Chapters()[1].Scenes()[0])
	assertFolderCount(t, manuscript.Works()[0], 2)
	assertFolderMetadata(t, manuscript.Works()[0].Folders()[0], 0, "Epi")
	assertSceneCount(t, manuscript.Works()[0].Folders()[0], 1)
	assertFolderMetadata(t, manuscript.Works()[0].Folders()[1], 1, "Ch")
	assertSceneCount(t, manuscript.Works()[0].Folders()[1], 1)
}

func TestMustLoadMultiWork(t *testing.T) {
	manuscript := ms.MustLoad("../../test-data/manuscripts/multi-work")

	assertWorkCount(t, manuscript, 2)

	assertWorkMetadata(t, manuscript.Works()[0], "Multi Work Book 1 Example", "Multi 1", "Geoff Coffey", "Coffey")
	assertSceneCount(t, manuscript.Works()[0], 2)
	assertSceneMetadata(t, manuscript.Works()[0].Scenes()[0], 0, "Scene")
	assertSceneMetadata(t, manuscript.Works()[0].Scenes()[1], 1, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Scenes()[0])
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Scenes()[1])

	assertWorkMetadata(t, manuscript.Works()[1], "Multi Work Book 2 Example", "Multi 2", "Geoff Coffey", "Coffey")
	assertSceneCount(t, manuscript.Works()[1], 2)
	assertSceneMetadata(t, manuscript.Works()[1].Scenes()[0], 0, "Scene")
	assertSceneMetadata(t, manuscript.Works()[1].Scenes()[1], 1, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[1].Scenes()[0])
	assertSceneTextStartsWithLorem(t, manuscript.Works()[1].Scenes()[1])

	assertChapterCount(t, manuscript.Works()[0], 0)

	assertFolderCount(t, manuscript.Works()[0], 0)
}

func TestLoadOutline(t *testing.T) {
	manuscript := ms.MustLoad("../../test-data/manuscripts/outline")
	assertWorkCount(t, manuscript, 1)
	assertWorkMetadata(t, manuscript.Works()[0], "Outline Example", "Outline", "Geoff Coffey", "Coffey")

	assertFolderCount(t, manuscript.Works()[0], 3)

	assertFolderMetadata(t, manuscript.Works()[0].Folders()[0], 0, "Act 1")
	assertFolderCount(t, manuscript.Works()[0].Folders()[0], 2)
	assertFolderMetadata(t, manuscript.Works()[0].Folders()[0].Folders()[0], 0, "At home")
	assertSceneCount(t, manuscript.Works()[0].Folders()[0].Folders()[0], 2)
	assertSceneMetadata(t, manuscript.Works()[0].Folders()[0].Folders()[0].Scenes()[0], 0, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Folders()[0].Folders()[0].Scenes()[0])
	assertSceneMetadata(t, manuscript.Works()[0].Folders()[0].Folders()[0].Scenes()[1], 1, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Folders()[0].Folders()[0].Scenes()[1])

	assertFolderMetadata(t, manuscript.Works()[0].Folders()[0].Folders()[1], 1, "On the way")
	assertSceneCount(t, manuscript.Works()[0].Folders()[0].Folders()[1], 2)
	assertSceneMetadata(t, manuscript.Works()[0].Folders()[0].Folders()[1].Scenes()[0], 0, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Folders()[0].Folders()[1].Scenes()[0])
	assertSceneMetadata(t, manuscript.Works()[0].Folders()[0].Folders()[1].Scenes()[1], 1, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Folders()[1].Folders()[0].Scenes()[1])

	assertFolderMetadata(t, manuscript.Works()[0].Folders()[1], 1, "Act 2")
	assertFolderCount(t, manuscript.Works()[0].Folders()[1], 2)
	assertFolderMetadata(t, manuscript.Works()[0].Folders()[1].Folders()[0], 0, "Begin to climb")
	assertSceneCount(t, manuscript.Works()[0].Folders()[1].Folders()[0], 2)
	assertSceneMetadata(t, manuscript.Works()[0].Folders()[1].Folders()[0].Scenes()[0], 0, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Folders()[1].Folders()[0].Scenes()[0])
	assertSceneMetadata(t, manuscript.Works()[0].Folders()[1].Folders()[0].Scenes()[1], 1, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Folders()[1].Folders()[0].Scenes()[1])

	assertFolderMetadata(t, manuscript.Works()[0].Folders()[1].Folders()[1], 1, "Pinnacle")
	assertSceneCount(t, manuscript.Works()[0].Folders()[1].Folders()[1], 2)
	assertSceneMetadata(t, manuscript.Works()[0].Folders()[1].Folders()[1].Scenes()[0], 0, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Folders()[1].Folders()[1].Scenes()[0])
	assertSceneMetadata(t, manuscript.Works()[0].Folders()[1].Folders()[1].Scenes()[1], 1, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Folders()[1].Folders()[1].Scenes()[1])

	assertFolderMetadata(t, manuscript.Works()[0].Folders()[2], 2, "Act 3")
	assertFolderCount(t, manuscript.Works()[0].Folders()[2], 2)
	assertFolderMetadata(t, manuscript.Works()[0].Folders()[2].Folders()[0], 0, "Realization")
	assertSceneCount(t, manuscript.Works()[0].Folders()[2].Folders()[0], 2)
	assertSceneMetadata(t, manuscript.Works()[0].Folders()[2].Folders()[0].Scenes()[0], 0, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Folders()[2].Folders()[0].Scenes()[0])
	assertSceneMetadata(t, manuscript.Works()[0].Folders()[2].Folders()[0].Scenes()[1], 1, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Folders()[2].Folders()[0].Scenes()[1])

	assertFolderMetadata(t, manuscript.Works()[0].Folders()[2].Folders()[1], 1, "Resolution")
	assertSceneCount(t, manuscript.Works()[0].Folders()[2].Folders()[1], 2)
	assertSceneMetadata(t, manuscript.Works()[0].Folders()[2].Folders()[1].Scenes()[0], 0, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Folders()[2].Folders()[1].Scenes()[0])
	assertSceneMetadata(t, manuscript.Works()[0].Folders()[2].Folders()[1].Scenes()[1], 1, "Scene")
	assertSceneTextStartsWithLorem(t, manuscript.Works()[0].Folders()[2].Folders()[1].Scenes()[1])

	assertChapterCount(t, manuscript.Works()[0], 5)
	assertSceneCount(t, manuscript.Works()[0].Chapters()[0], 2)
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[0].Scenes()[0], 0, "Scene")
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[0].Scenes()[1], 1, "Scene")

	assertSceneCount(t, manuscript.Works()[0].Chapters()[1], 2)
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[1].Scenes()[0], 0, "Scene")
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[1].Scenes()[1], 1, "Scene")

	assertSceneCount(t, manuscript.Works()[0].Chapters()[2], 4)
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[2].Scenes()[0], 0, "Scene")
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[2].Scenes()[1], 1, "Scene")
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[2].Scenes()[2], 0, "Scene")
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[2].Scenes()[3], 1, "Scene")

	assertSceneCount(t, manuscript.Works()[0].Chapters()[3], 2)
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[3].Scenes()[0], 0, "Scene")
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[3].Scenes()[1], 1, "Scene")

	assertSceneCount(t, manuscript.Works()[0].Chapters()[4], 2)
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[4].Scenes()[0], 0, "Scene")
	assertSceneMetadata(t, manuscript.Works()[0].Chapters()[4].Scenes()[1], 1, "Scene")
}

func TestLoadInvalidOrphanScenes(t *testing.T) {
	_, err := ms.Load("../../test-data/manuscripts/invalid/orphan-scenes")
	if err == nil {
		t.Errorf("load orphan-scenes did not produce error")
	} else if expected := "has scenes before the first chapter"; !strings.HasSuffix(err.Error(), expected) {
		t.Fatalf("load orphan-scenes error = %v; expected ...%v", err, expected)
	}

	_, err = ms.Load("../../test-data/manuscripts/invalid/no-works")
	if err == nil {
		t.Errorf("load no-works did not produce error")
	} else if expected := "has no works"; !strings.HasSuffix(err.Error(), expected) {
		t.Fatalf("load orphan-scenes error = %v; expected ...%v", err, expected)
	}

}
