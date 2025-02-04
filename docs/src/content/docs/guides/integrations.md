---
title: CI/CD Integrations
---

## Badges

It's possible to add badges to your documentation. DepsHub supports the following badge:

```
![dependencies](https://img.shields.io/github/actions/workflow/status/<your_organization>/<your_repository>/<your_workflow_file>.yml?branch=main&label=DepsHub&fedcba&logo=data:image/svg%2bxml;base64,PHN2ZyB3aWR0aD0iNzAiIGhlaWdodD0iNzgiIHZpZXdCb3g9IjAgMCA3MCA3OCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTMuOTAwMSA0NS42ODk5QzQuMDMyNDUgNDUuNzk1MSA0LjE2ODQ0IDQ1Ljg5ODYgNC4zMDgwOCA0Ni4wMDA2TDI1LjMwNzcgNjEuMzMwOEMzMC42NjA2IDY1LjIzODQgMzkuMzM5MiA2NS4yMzg0IDQ0LjY5MiA2MS4zMzA4TDY1LjY5MTcgNDYuMDAwNkM2NS44MzEzIDQ1Ljg5ODYgNjUuOTY3MyA0NS43OTUgNjYuMDk5NyA0NS42ODk5QzcxLjA0MDkgNDkuNjE1OCA3MC45MDQ5IDU1LjcyNDQgNjUuNjkxNyA1OS41MzAxTDQ0LjY5MjEgNzQuODYwM0MzOS4zMzkyIDc4Ljc2OCAzMC42NjA2IDc4Ljc2OCAyNS4zMDc4IDc0Ljg2MDNMNC4zMDgxMSA1OS41MzAxQy0wLjkwNTA2MyA1NS43MjQ0IC0xLjA0MTA3IDQ5LjYxNTggMy45MDAxIDQ1LjY4OTlaIiBmaWxsPSJ3aGl0ZSIvPgo8cGF0aCBkPSJNNC4zMDgwOCAzMi43NjU0Qy0xLjA0NDc1IDI4Ljg1NzcgLTEuMDQ0NzUgMjIuNTIyMSA0LjMwODA4IDE4LjYxNDVMMjUuMzA3NyAzLjI4NDI3QzMwLjY2MDYgLTAuNjIzNDAzIDM5LjMzOTIgLTAuNjIzNDAzIDQ0LjY5MiAzLjI4NDI3TDY1LjY5MTcgMTguNjE0NUM3MS4wNDQ1IDIyLjUyMjEgNzEuMDQ0NSAyOC44NTc3IDY1LjY5MTcgMzIuNzY1NEw0NC42OTIgNDguMDk1NkMzOS4zMzkyIDUyLjAwMzMgMzAuNjYwNiA1Mi4wMDMzIDI1LjMwNzcgNDguMDk1Nkw0LjMwODA4IDMyLjc2NTRaIiBmaWxsPSJ3aGl0ZSIvPgo8L3N2Zz4K)
```

Where:
- `<your_organization>` is the name of your GitHub organization.
- `<your_repository>` is the name of your GitHub repository.
- `<your_workflow_file>` is the name of your GitHub Actions workflow file.

For example, the following badge shows the status of the `depshub.yml` workflow in the `DepsHubHQ/depshub` repository:

```
![dependencies](https://img.shields.io/github/actions/workflow/status/depshubhq/depshub/depshub.yml?branch=main&label=DepsHub&fedcba&logo=data:image/svg%2bxml;base64,PHN2ZyB3aWR0aD0iNzAiIGhlaWdodD0iNzgiIHZpZXdCb3g9IjAgMCA3MCA3OCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTMuOTAwMSA0NS42ODk5QzQuMDMyNDUgNDUuNzk1MSA0LjE2ODQ0IDQ1Ljg5ODYgNC4zMDgwOCA0Ni4wMDA2TDI1LjMwNzcgNjEuMzMwOEMzMC42NjA2IDY1LjIzODQgMzkuMzM5MiA2NS4yMzg0IDQ0LjY5MiA2MS4zMzA4TDY1LjY5MTcgNDYuMDAwNkM2NS44MzEzIDQ1Ljg5ODYgNjUuOTY3MyA0NS43OTUgNjYuMDk5NyA0NS42ODk5QzcxLjA0MDkgNDkuNjE1OCA3MC45MDQ5IDU1LjcyNDQgNjUuNjkxNyA1OS41MzAxTDQ0LjY5MjEgNzQuODYwM0MzOS4zMzkyIDc4Ljc2OCAzMC42NjA2IDc4Ljc2OCAyNS4zMDc4IDc0Ljg2MDNMNC4zMDgxMSA1OS41MzAxQy0wLjkwNTA2MyA1NS43MjQ0IC0xLjA0MTA3IDQ5LjYxNTggMy45MDAxIDQ1LjY4OTlaIiBmaWxsPSJ3aGl0ZSIvPgo8cGF0aCBkPSJNNC4zMDgwOCAzMi43NjU0Qy0xLjA0NDc1IDI4Ljg1NzcgLTEuMDQ0NzUgMjIuNTIyMSA0LjMwODA4IDE4LjYxNDVMMjUuMzA3NyAzLjI4NDI3QzMwLjY2MDYgLTAuNjIzNDAzIDM5LjMzOTIgLTAuNjIzNDAzIDQ0LjY5MiAzLjI4NDI3TDY1LjY5MTcgMTguNjE0NUM3MS4wNDQ1IDIyLjUyMjEgNzEuMDQ0NSAyOC44NTc3IDY1LjY5MTcgMzIuNzY1NEw0NC42OTIgNDguMDk1NkMzOS4zMzkyIDUyLjAwMzMgMzAuNjYwNiA1Mi4wMDMzIDI1LjMwNzcgNDguMDk1Nkw0LjMwODA4IDMyLjc2NTRaIiBmaWxsPSJ3aGl0ZSIvPgo8L3N2Zz4K)
```

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
