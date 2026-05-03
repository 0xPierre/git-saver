package source

import (
	"context"
	"fmt"
	"os"

	"github.com/0xPierre/git-saver/internal/model"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/client"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
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

	// for loop to iterate overage pages
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

func (s *githubSource) CloneOrUpdate(cloneUrl string, dest string) error {
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		_, err = git.PlainClone(dest, &git.CloneOptions{
			URL:    cloneUrl,
			Mirror: true,
			ClientOptions: []client.Option{
				client.WithHTTPAuth(&http.BasicAuth{
					Username: "oauth2",
					Password: s.token,
				}),
			},
		})
		return err
	}

	return nil
}
