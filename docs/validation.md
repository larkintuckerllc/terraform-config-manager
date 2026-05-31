# Validating Project Configurations

## Overview

When a project owner submits a change to `project.tf`, the terraform-config-manager validates the file to ensure it conforms to the platform team's standards. This validation is intended to run as part of a review build (e.g., a CI check on pull requests).

## Problem

The two-file ownership model gives project owners control over `project.tf`, but that control needs guardrails. Without validation, a project owner could introduce inline resources, reference unapproved modules, use unpinned module versions, or otherwise deviate from the platform team's expectations.

## Scope

Validation focuses exclusively on `project.tf`. All other files in the Terraform configuration are owned by the platform team and protected via CODEOWNERS — they don't need runtime validation.

## Rules

### Only module blocks are allowed

`project.tf` must contain zero or more `module` blocks. Any other block types — `resource`, `data`, `provider`, `terraform`, `locals`, `variable`, `output` — are rejected.

### Every module must have a source attribute

Module blocks without a `source` attribute are rejected.

### Module sources must come from the approved repository

Every module block must have a `source` attribute that references the approved Terraform modules repository:

```
git::https://github.com/larkintuckerllc/terraform-modules.git//<module-name>?ref=<tag>
```

Any other source (local paths, other Git repos, Terraform registry, etc.) is rejected.

### Module sources must be pinned to a tag

The `source` attribute must include a `?ref=` parameter pointing to a Git tag. Unpinned sources or branch references are rejected.

## What Is Not Validated

- **Module names.** Whether the referenced module subdirectory exists in the approved repository is not checked — `terraform init` will catch invalid module references.
- **Module parameter values.** Whether the inputs passed to a module are correct is left to `terraform validate` and `terraform plan`. The config manager validates structure, not semantics.
- **Other files.** `main.tf`, `.terraform-version`, and other platform-owned files are not validated — they are protected by CODEOWNERS.

## Usage

```sh
terraform-config-manager validate [-dir=<path>]
```

- `-dir` — path to the project's Terraform configuration directory (defaults to current directory)

In a CI review build, this is typically just `terraform-config-manager validate` run from the project repo root.