package model

type Repo struct {
	Name        string
	Description string
	CloneURL    string
	Private     bool
	Source      string
	localPath   string
}
