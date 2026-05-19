# Git-Saver

Keep repositories in sync, backup, and mirror them across multiple git hosting providers.

## Features

- Pull repositories from GitHub
- Push to multiple targets: GitLab, Gitea, or a local directory
- Sync all repositories from a source, or a selected list
- Configuration via a single YAML file

## Usage

Create a `config.yml` file:

```yaml
sources:
  - type: github
    url: https://github.com
    token: <your-github-token>

sync-all-repo: false # if set to true, `source-repositories` is not used.
source-repositories:
  - my-repo
  - another-repo
default-pull-directory: tmp-pull

targets:
  - type: gitlab
    url: https://gitlab.com
    token: <your-gitlab-token>
  - type: gitea
    url: https://gitea.com
    token: <your-gitea-token>
  - type: local
    dir: local
```

Build and run:

```sh
go build -o git-saver
./git-saver sync
```

Use a custom config path with `-c`:

```sh
./git-saver sync -c /path/to/config.yml
```
