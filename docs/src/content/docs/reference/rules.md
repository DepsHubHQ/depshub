---
title: Rules list
description: The list of rules that can be used in the configuration file.
---

The list of rules that can be used in the configuration file.
Some of the rules have additional options that can accept values.

### allowed-licenses

Set the allowed licenses for the manifest file.

| Type          | Default Value           |
| ------------- | ----------------------- |
| Array<string> | `["MIT", "Apache-2.0"]` |

### lockfile

Checks if the lockfile is present.

### max-libyear

Set the maximum allowed [libyear](https://libyear.com/) for the manifest file.

| Type   | Default Value |
| ------ | ------------- |
| Number | `25`          |

### max-major-updates

Set the maximum **percentage** of major updates for the manifest file.

| Type                | Default Value |
| ------------------- | ------------- |
| Number (Percentage) | `20.0`        |

### max-minor-updates

Set the maximum **percentage** of minor updates for the manifest file.

| Type                | Default Value |
| ------------------- | ------------- |
| Number (Percentage) | `40.0`        |

### max-package-age

Set the maximum allowed package age for the manifest file in months.

| Type            | Default Value |
| --------------- | ------------- |
| Number (Months) | `12`          |

### max-patch-updates

Set the maximum **percentage** of patch updates for the manifest file.

| Type                | Default Value |
| ------------------- | ------------- |
| Number (Percentage) | `60`          |

### min-weekly-downloads

Set the minimum allowed package weekly downloads for the manifest file.

| Type   | Default Value |
| ------ | ------------- |
| Number | `1000`        |

### no-any-tag

Forbids the usage of the **any** tags (`*`, `latest` or empty version ` `) in the manifest file.

### no-deprecated

Forbids the usage of deprecated packages in the manifest file.

### no-duplicates

Forbids the usage of duplicate packages in the manifest file.

### no-multiple-versions

Forbids the usage of multiple versions of the same package in the repository.

### no-pre-release

Forbids the usage of pre-release (-alpha, -beta etc) packages in the manifest file.

### no-unstable

Forbids the usage of unstable (<1.0.0) packages in the manifest file.

### sorted

Checks if all the dependencies in the manifest file are sorted alphabetically.
