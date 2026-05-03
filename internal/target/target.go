package target

import (
	"fmt"

	"github.com/0xPierre/git-saver/internal/config"
	"github.com/0xPierre/git-saver/internal/model"
)

type Target interface {
	CreateRepoIfDoesntExist(name string) ([]model.Repo, error)
	PushRepo(repo model.Repo) error
}

func New(t config.Target) (Target, error) {
	switch t.Type {
	case "gitea":
		return newGitea(t.Token), nil
	default:
		return nil, fmt.Errorf("unknow source type %q", t.Type)
	}
}
