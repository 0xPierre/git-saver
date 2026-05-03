package source

import (
	"fmt"

	"github.com/0xPierre/git-saver/internal/config"
	"github.com/0xPierre/git-saver/internal/model"
)

type Source interface {
	ListRepos() ([]model.Repo, error)
	CloneOrUpdate(cloneUrl string, dest string) error
}

func New(s config.Source) (Source, error) {
	switch s.Type {
	case "github":
		return newGithub(s.Token), nil
	default:
		return nil, fmt.Errorf("unknow source type %q", s.Type)
	}
}
