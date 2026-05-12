package target

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/0xPierre/git-saver/internal/gitops"
	"github.com/0xPierre/git-saver/internal/model"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type gitlabTarget struct {
	client   *gitlab.Client
	token    string
	userName string
}

func newGitlab(token, baseURL string) Target {
	opts := []gitlab.ClientOptionFunc{}
	if baseURL != "" {
		opts = append(opts, gitlab.WithBaseURL(baseURL))
	}

	client, err := gitlab.NewClient(token, opts...)
	if err != nil {
		log.Fatalf("gitlab: NewClient: %v", err)
	}

	user, _, err := client.Users.CurrentUser()
	if err != nil {
		log.Fatalf("gitlab: CurrentUser: %v", err)
	}

	return &gitlabTarget{client: client, token: token, userName: user.Username}
}

// visibilityFor maps the source-side Private flag to a GitLab visibility level.
// We never set Internal because the source model does not expose that signal.
func visibilityFor(private bool) gitlab.VisibilityValue {
	if private {
		return gitlab.PrivateVisibility
	}
	return gitlab.PublicVisibility
}

// pathSlug returns the GitLab URL slug for repo.Name. GitLab stores project
// paths in lowercase, so we must match that for both lookups (case-sensitive)
// and creation (otherwise GitLab returns "path has already been taken" when
// the normalised path collides with an existing one).
func (t *gitlabTarget) pathSlug(name string) string {
	return strings.ToLower(name)
}

func (t *gitlabTarget) projectPath(name string) string {
	return t.userName + "/" + t.pathSlug(name)
}

func (t *gitlabTarget) CreateRepoIfDoesntExist(repo *model.Repo) (string, error) {
	want := visibilityFor(repo.Private)
	slug := t.pathSlug(repo.Name)

	project, resp, err := t.client.Projects.GetProject(t.projectPath(repo.Name), nil)
	if err == nil {
		if project.Visibility != want {
			if _, _, err := t.client.Projects.EditProject(project.ID, &gitlab.EditProjectOptions{
				Visibility: &want,
			}); err != nil {
				return "", fmt.Errorf("updating visibility on %q: %w", repo.Name, err)
			}
			fmt.Printf("-- Visibility updated")
		}
		return project.HTTPURLToRepo, nil
	}

	if resp == nil || resp.StatusCode != http.StatusNotFound {
		return "", fmt.Errorf("getting %q: %w", repo.Name, err)
	}

	project, _, err = t.client.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:        &repo.Name,
		Path:        &slug,
		Description: &repo.Description,
		Visibility:  &want,
	})
	if err != nil {
		return "", fmt.Errorf("creating %q: %w", repo.Name, err)
	}
	return project.HTTPURLToRepo, nil
}

func (t *gitlabTarget) PushRepo(repo *model.Repo, remoteURL string) error {
	return gitops.ArchivePush(repo.LocalPath, remoteURL, gitops.Auth{
		Username: t.userName,
		Password: t.token,
	})
}
