all: build print

# semver string from git
GIT_SEMVER_STR := $(shell git describe --tags --dirty)

# pass semver string to compiler
FLAGS := -ldflags "-X main.GitSemverStr=$(GIT_SEMVER_STR)"

# export cmd on windows
ifdef OS
   EXPORT = set
else
   ifeq ($(shell uname), Linux)
      EXPORT = export
   endif
endif

# environment variables for each target
ENV_WIN_X64 = GOOS=windows GOARCH=amd64
ENV_DAR_X64 = GOOS=darwin GOARCH=amd64
ENV_DAR_ARM = GOOS=darwin GOARCH=arm64

# output filename for each target
TARGET_WIN_X64 = windows_x64.exe
TARGET_DAR_X64 = darwin_x64
TARGET_DAR_ARM = darwin_arm


build: # build for each target
	@$(EXPORT) ${ENV_WIN_X64} && go build $(FLAGS) -o ./bin/${TARGET_WIN_X64} .
	@$(EXPORT) ${ENV_DAR_X64} && go build $(FLAGS) -o ./bin/${TARGET_DAR_X64} .
	@$(EXPORT) ${ENV_DAR_ARM} && go build $(FLAGS) -o ./bin/${TARGET_DAR_ARM} .

print:
	@echo Built binaries sucessfully for ${GIT_SEMVER_STR}