# go-git-stats

CLI tool for git statistics

## COMMANDS

### _repo_

Get all repositories.

```sh
# When Github access token is set to GGS_TOKEN (environment variable)
$ echo $GGS_TOKEN
> ghq_....

# Repositories for authenticated user with Github access token
$ ggs repo
# abbreviation command
$ ggs r

# Public repositories without Github access token
$ ggs repo -name kokoichi206
$ ggs r -n kokoichi206
```

### _stats_

Get statistics of a specific repository.

```sh
# Github access token should be set to GGS_TOKEN (environment variable)
# when you get private repopsitory stats.
$ echo $GGS_TOKEN
> ghq_....

$ ggs stats -name kokoichi206/go-git-stats
# abbreviation command
$ ggs s -name kokoichi206/go-git-stats
```

## INSTALLATION

Built binaries are available from GitHub Releases.
https://github.com/kokoichi206/go-git-stats/releases

## LICENSE

under [MIT License](./LICENSE).
