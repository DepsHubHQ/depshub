---
title: CI/CD Integrations
---

## Azure DevOps

## CircleCI

## GitHub Actions

You can find the official DepsHub GitHub Action in this repository: [depshub-action](https://github.com/DepsHubHQ/github-action).

Example:

```yaml
name: Run DepsHub

on:
  push:
    branches:
      - main
  pull_request:

jobs:
  depshub-lint:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Run DepsHub
        uses: DepsHubHQ/depshub-action@v1
        with:
          path: "./src" # optional
          config: "./depshub-config.yml" # optional
```

More options can be found in the repository's README.

## GitLab CI/CD

## Jenkins

## TeamCity

## Travis CI
