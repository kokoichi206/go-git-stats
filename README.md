# go-git-stats

CLI tool for git statistics

## COMMANDS

### _repo_

Get all repositories.
Sample files are found [here](./cmd/generate/_examples).

```sh
# When Github access token is set to GGS_TOKEN (environment variable)
$ echo $GSS_TOKEN
> ghq_....

# Repositories for authenticated user with Github access token
$ ggs repo
# abbreviation command
$ ggs r

# Public repositories without Github access token
$ ggs repo -name kokoichi206
$ ggs r -n kokoichi206
```

## INSTALLATION

Built binaries are available from GitHub Releases.
https://github.com/kokoichi206/go-git-stats/releases

## LICENSE

under [MIT License](./LICENSE).
