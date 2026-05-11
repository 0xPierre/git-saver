package source

import (
	"context"
	"fmt"
	"os"

	"github.com/0xPierre/git-saver/internal/gitops"
	"github.com/0xPierre/git-saver/internal/model"
	"github.com/google/go-github/v85/github"
)

var sourceName = "github"

type githubSource struct {
	client *github.Client
	token  string
}

func newGithub(token string) Source {
	client := github.NewClient(nil).WithAuthToken(token)

	return &githubSource{client: client, token: token}
}

func (s *githubSource) ListRepos() ([]model.Repo, error) {
	var allRepo []model.Repo

	opts := &github.RepositoryListByAuthenticatedUserOptions{ListOptions: github.ListOptions{PerPage: 100}}

	for {
		repos, resp, err := s.client.Repositories.ListByAuthenticatedUser(context.Background(), opts)
		if err != nil {
			return []model.Repo{}, fmt.Errorf("Listing %v repos: %w", sourceName, err)
		}

		for _, r := range repos {
			allRepo = append(allRepo, model.Repo{
				Name:        r.GetName(),
				Description: r.GetDescription(),
				CloneURL:    r.GetCloneURL(),
				Private:     r.GetPrivate(),
				Source:      sourceName,
			})
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allRepo, nil
}

func (s *githubSource) auth() gitops.Auth {
	return gitops.Auth{Username: "oauth2", Password: s.token}
}

func (s *githubSource) CloneOrUpdate(cloneUrl string, dest string) error {
	_, err := os.Stat(dest)
	if os.IsNotExist(err) {
		return gitops.MirrorClone(cloneUrl, dest, s.auth())
	}
	if err != nil {
		return fmt.Errorf("stat %q: %w", dest, err)
	}
	return gitops.MirrorFetch(dest, s.auth())
}
