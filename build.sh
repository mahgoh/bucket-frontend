# export git semver string to environment
export GIT_SEMVER_STR=$(git describe --tags --dirty)

# pass env var to compiler
go build -ldflags "-X main.GitSemverStr=$GIT_SEMVER_STR"