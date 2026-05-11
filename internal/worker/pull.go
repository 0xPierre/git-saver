package worker

import (
	"fmt"
	"strings"

	"github.com/0xPierre/git-saver/internal/config"
	"github.com/0xPierre/git-saver/internal/model"
	"github.com/0xPierre/git-saver/internal/source"
)

// filterRepositories filters the full repo list down to only the ones listed
// in source-repositories, unless sync-all-repo is true.
func filterRepositories(repos []model.Repo, cfg *config.Config) []model.Repo {
	if cfg.SyncAllRepos {
		return repos
	}

	allowed := make(map[string]bool, len(cfg.SourceRepositories))

	for _, sr := range cfg.SourceRepositories {
		allowed[sr] = true
	}

	var filteredRepos []model.Repo

	for _, repo := range repos {
		if allowed[repo.Name] {
			filteredRepos = append(filteredRepos, repo)
		}
	}

	return filteredRepos
}

// Pull iterates over all configured sources, lists their repositories,
// mirrors them into cfg.DefaultPullDirectory, and returns the full list of pulled repos.
func Pull(cfg *config.Config) ([]model.Repo, error) {
	fmt.Printf("Pulling...\n")

	var allRepos []model.Repo

	for _, s := range cfg.Sources {
		fmt.Printf("[+] Source: %v\n", s.Type)
		src, err := source.New(s)
		if err != nil {
			return nil, fmt.Errorf("source %q: %w", s.Type, err)
		}

		repos, err := src.ListRepos()
		if err != nil {
			return nil, fmt.Errorf("source %q: %w", s.Type, err)
		}

		filteredRepos := filterRepositories(repos, cfg)

		names := make([]string, len(filteredRepos))
		for i, repo := range filteredRepos {
			names[i] = repo.Name
		}
		fmt.Printf("[+] %d repos: %s\n", len(filteredRepos), strings.Join(names, ", "))

		// Cloning repos in default-pull-directory
		for _, repo := range filteredRepos {
			fmt.Printf("- Cloning %v\n", repo.Name)

			dest := cfg.DefaultPullDirectory + "/" + s.Type + "/" + repo.Name
			src.CloneOrUpdate(repo.CloneURL, dest)
			repo.LocalPath = dest

			allRepos = append(allRepos, repo)
		}

	}

	return allRepos, nil
}
