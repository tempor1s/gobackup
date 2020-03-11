# gobackup

[![Go Report Card](https://goreportcard.com/badge/github.com/tempor1s/gobackup)](https://goreportcard.com/report/github.com/tempor1s/gobackup)

A tool that allows you to backup all your GitHub, GitLab & BitBucket repos onto your local machine, and even upload them to a different repo manager!

### Table of Contents

1. [Installation]("#installation")
2. [Usage]("#usage")
3. [Milestones]("#milestones")

## Installation

A quick guide on how to install the tool. Not currently available.

```bash
brew tap tempor1s/gobackup
brew install gobackup
```

## Usage

TODO

## Milestones

- [x] Clone a single repository to your computer through CLI.
- [x] Clone multiple repositories to local computer
- [x] Clone multiple reposiories using concurrency
- [x] Add support for cloning GitLab repos
- [ ] Add support for cloning BitBucket repos
- [ ] Upload cloned repositories to other services like GitLab or BitBucket
- [ ] Do the above concurrently
- [ ] Do the reverse, pull repos from BitBucket / GitLab and clone them or upload them to GitHub. Ideally want to be platform agnostic.
- [ ] Set up Cron Job to periodically backup new repos / changes to old repos to local / other service like GitLab or BitBucket
