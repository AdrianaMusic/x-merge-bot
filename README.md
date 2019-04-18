# x-merge-bot
Auto merge pull request

### Usage

```yml
version: '2'
services:
  x-merge-bot:
    image: echoulen/x-merge-bot:latest
    environment:
      - TOKEN=<GITHUB_TOKEN>
      - REPO=repo1 repo2 repo3
      - OWNER=xteamstudio
```