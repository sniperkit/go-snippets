package main

type File struct {
	Name        string
	CommitCount int
}

type Files []File

func (f Files) Len() int {
	return len(f)
}

func (f Files) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f Files) Less(i, j int) bool {
	return f[i].CommitCount > f[j].CommitCount
}
