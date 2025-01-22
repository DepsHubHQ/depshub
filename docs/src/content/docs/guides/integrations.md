---
title: CI/CD Integrations
---

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

More options can be found in the repository's [README](https://github.com/DepsHubHQ/github-action).

## Other

DepsHub is available as a CLI tool. You can install it on your CI/CD system as described in the [installation](/installation) guide.
