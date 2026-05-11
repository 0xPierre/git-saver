// Package gitops centralises the bare-mirror git operations used by sources
// and targets: cloning, refreshing, and pushing with mirror/archive semantics.
package gitops

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing/client"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

// Auth carries HTTP basic auth for a remote. Leave Password empty to push or fetch anonymously
type Auth struct {
	Username string
	Password string
}

func (a Auth) clientOptions() []client.Option {
	if a.Password == "" {
		return nil
	}
	return []client.Option{
		client.WithHTTPAuth(&http.BasicAuth{
			Username: a.Username,
			Password: a.Password,
		}),
	}
}

// mirrorRefSpec mirrors every ref under refs/,branches, tags, notes, and
// any provider-specific namespace
var mirrorRefSpec = []config.RefSpec{"+refs/*:refs/*"}

// MirrorClone bare-clones url into path. Use when path does not yet exist.
func MirrorClone(url, path string, auth Auth) error {
	_, err := git.PlainClone(path, &git.CloneOptions{
		URL:           url,
		Mirror:        true,
		ClientOptions: auth.clientOptions(),
	})
	if err != nil {
		return fmt.Errorf("cloning %q into %q: %w", url, path, err)
	}
	return nil
}

// MirrorFetch refreshes an existing bare mirror at path so it stays an exact
// copy of the source: force-update on rebases, prune refs deleted upstream.
func MirrorFetch(path string, auth Auth) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("opening %q: %w", path, err)
	}

	err = r.Fetch(&git.FetchOptions{
		RemoteName:    "origin",
		RefSpecs:      mirrorRefSpec,
		Force:         true,
		Prune:         true,
		ClientOptions: auth.clientOptions(),
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("fetching %q: %w", path, err)
	}
	return nil
}

// ArchivePush pushes the bare repo at path to remoteURL with archive semantics:
// force-update existing refs (so upstream force-pushes are propagated) but
// never delete refs that disappeared upstream, the target accumulates history.
func ArchivePush(path, remoteURL string, auth Auth) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("opening %q: %w", path, err)
	}

	err = r.Push(&git.PushOptions{
		RemoteURL:     remoteURL,
		RefSpecs:      mirrorRefSpec,
		Force:         true,
		ClientOptions: auth.clientOptions(),
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("pushing to %q: %w", remoteURL, err)
	}
	return nil
}
