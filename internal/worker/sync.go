package worker

import (
	"github.com/0xPierre/git-saver/internal/config"
)

func Sync(cfg *config.Config) error {
	// First pull all needed repositories
	repositories, err := Pull(cfg)
	if err != nil {
		return err
	}

	// Then push to all targets
	err = Push(cfg, repositories)
	if err != nil {
		return err
	}

	return nil
}
