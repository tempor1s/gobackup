# gobackup

[![Go Report Card](https://goreportcard.com/badge/github.com/tempor1s/gobackup)](https://goreportcard.com/report/github.com/tempor1s/gobackup)

A tool that allows you to backup all your GitHub, GitLab repos onto your local machine, and even upload them to a different repo manager!

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
- [x] Upload cloned repositories to GitLab 
- [x] Do the above concurrently
- [ ] Do the reverse, upload cloned repositories to GitHub
- [ ] Do the above concurrently
- [ ] Add backup command that will download the repo to memory (or disk I guess) and then instant upload it to new platform of choice
- [ ] Respect the privacy status of a cloned repo when we upload it again
- [ ] Set up Cron Job to periodically backup new repos / changes to old repos to local / other service like GitLab
