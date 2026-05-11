package target

import (
	"fmt"

	"github.com/0xPierre/git-saver/internal/config"
	"github.com/0xPierre/git-saver/internal/model"
)

type Target interface {
	CreateRepoIfDoesntExist(repo *model.Repo) (string, error)
	PushRepo(repo *model.Repo, remoteURL string) error
}

func New(t config.Target) (Target, error) {
	switch t.Type {
	case "gitea":
		return newGitea(t.Token, t.URL), nil
	case "gitlab":
		return newGitlab(t.Token, t.URL), nil
	default:
		return nil, fmt.Errorf("unknow source type %q", t.Type)
	}
}
