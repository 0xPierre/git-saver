package worker

import (
	"fmt"

	"github.com/0xPierre/git-saver/internal/config"
	"github.com/0xPierre/git-saver/internal/model"
)

func Push(cfg *config.Config, repos []model.Repo) error {
	fmt.Printf("Pushing...")
	return nil
}
