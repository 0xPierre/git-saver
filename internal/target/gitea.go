package target

import (
	"fmt"
	"log"
	"net/http"

	"code.gitea.io/sdk/gitea"
	"github.com/0xPierre/git-saver/internal/gitops"
	"github.com/0xPierre/git-saver/internal/model"
)

func IsNotFound(resp *gitea.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusNotFound
}

type giteaTargets struct {
	client   *gitea.Client
	token    string
	userName string
}

func newGitea(token, url string) Target {
	client, err := gitea.NewClient(url, gitea.SetToken(token))
	if err != nil {
		log.Fatal(err)
	}

	user, _, err := client.GetMyUserInfo()
	if err != nil {
		log.Fatalf("gitea: GetMyUserInfo: %v", err)
	}

	return &giteaTargets{token: token, client: client, userName: user.UserName}
}

func (t *giteaTargets) CreateRepoIfDoesntExist(repo *model.Repo) (string, error) {
	giteaRepo, resp, err := t.client.GetRepo(t.userName, repo.Name)
	if err == nil {
		// When the repo visibility change, we want to update it on the target
		if giteaRepo.Private != repo.Private {
			private := repo.Private
			if _, _, err := t.client.EditRepo(t.userName, repo.Name, gitea.EditRepoOption{
				Private: &private,
			}); err != nil {
				return "", fmt.Errorf("updating privacy on %q: %w", repo.Name, err)
			}
			fmt.Printf("-- Privacy updated")
		}
		return giteaRepo.CloneURL, nil
	}

	if !IsNotFound(resp) {
		return "", err
	}

	giteaRepo, _, err = t.client.CreateRepo(gitea.CreateRepoOption{
		Name:        repo.Name,
		Description: repo.Description,
		Private:     repo.Private,
	})
	if err != nil {
		return "", err
	}
	return giteaRepo.CloneURL, nil
}

func (t *giteaTargets) PushRepo(repo *model.Repo, remoteURL string) error {
	return gitops.ArchivePush(repo.LocalPath, remoteURL, gitops.Auth{
		Username: t.userName,
		Password: t.token,
	})
}
