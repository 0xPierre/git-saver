package worker

import (
	"fmt"

	"github.com/0xPierre/git-saver/internal/config"
	"github.com/0xPierre/git-saver/internal/model"
	"github.com/0xPierre/git-saver/internal/target"
	"golang.org/x/sync/errgroup"
)

func Push(cfg *config.Config, repos []model.Repo) error {
	fmt.Printf("Pushing...\n")

	if len(cfg.Targets) == 0 {
		return fmt.Errorf("no pushing targets in config")
	}

	for _, t := range cfg.Targets {
		fmt.Printf("[+] Target %v\n", t.Type)

		pushTarget, err := target.New(t)
		if err != nil {
			return fmt.Errorf("target %q: %w", t.Type, err)
		}

		g := new(errgroup.Group)
		g.SetLimit(20)

		for _, repo := range repos {
			fmt.Printf("- Uploading %v\n", repo.Name)

			g.Go(func() error {
				url, err := pushTarget.CreateRepoIfDoesntExist(&repo)
				if err != nil {
					return fmt.Errorf("creating %q on %s: %w", repo.Name, t.Type, err)
				}

				if err := pushTarget.PushRepo(&repo, url); err != nil {
					return err
				}

				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return err
		}

	}
	return nil
}
