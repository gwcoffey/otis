package ms2

func Load() (Manuscript, error) {
	manuscript, err := newDirNode("manuscript/")
	return manuscript, err
}

func MustLoad() Manuscript {
	manuscript, err := Load()
	if err != nil {
		panic(err)
	}
	return manuscript
}
