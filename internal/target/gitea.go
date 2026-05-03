package target

import "github.com/0xPierre/git-saver/internal/model"

type giteaTargets struct {
	// client *gitea.Client
	token string
}

func newGitea(token string) Target {
	// client := github.NewClient(nil).WithAuthToken(token)

	return &giteaTargets{token: token}
}

func (t *giteaTargets) CreateRepoIfDoesntExist(name string) ([]model.Repo, error) {
	return nil, nil
}

func (t *giteaTargets) PushRepo(repo model.Repo) error {
	return nil
}
