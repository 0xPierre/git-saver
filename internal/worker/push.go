package worker

import (
	"fmt"

	"github.com/0xPierre/git-saver/internal/config"
	"github.com/0xPierre/git-saver/internal/model"
	"github.com/0xPierre/git-saver/internal/target"
)

func Push(cfg *config.Config, repos []model.Repo) error {
	fmt.Printf("Pushing...\n")

	if len(cfg.Targets) == 0 {
		return fmt.Errorf("no pushing targets in config")
	}

	for _, t := range cfg.Targets {
		fmt.Printf("[+] Target %v\n", t.Type)

		target, err := target.New(t)
		if err != nil {
			return fmt.Errorf("target %q: %w", t.Type, err)
		}

		for _, repo := range repos {
			fmt.Printf("- Uploading %v\n", repo.Name)
			target.CreateRepoIfDoesntExist(repo.Name)
		}

	}
	return nil
}
