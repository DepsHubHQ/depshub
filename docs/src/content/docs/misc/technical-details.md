---
title: Technical Details
---

DepsHub is built as a CLI using Go and is available for multiple platforms, including Linux, macOS, and Windows.

## Data Sources

DepsHub uses the following data sources to provide insights:

- [pypi.org](https://pypi.org)
- [registry.npmjs.org](https://registry.npmjs.org)
- [crates.io](https://crates.io)
- [deps.dev](https://deps.dev)
- [hex.pm](https://hex.pm)

## Cache

DepsHub caches the data fetched from the data sources to improve the performance of the tool.
The cache is stored in the user's home directory under the `/.cache/depshub` directory.

- Windows: `%USERPROFILE%\.cache\depshub\`
- Linux/macOS: `~/.cache/depshub/`
