# Migrating Configurations

## Overview

When the platform team needs to update Terraform configurations — bumping provider versions, updating module tags, changing Terraform version constraints — the terraform-config-manager applies these changes as a series of versioned migrations, similar to database migrations.

## Problem

Terraform configurations managed by the scaffolder will drift from the current platform standards over time. A project scaffolded at version `0.1.0` may need provider version bumps, new module versions, or structural changes introduced in later versions. Applying all changes at once is fragile — each change may depend on the state left by the previous one, and skipping intermediate steps risks breaking configs.

## Approach

Each version of the terraform-config-manager defines a migration from the previous version. Migrations are applied incrementally, stepping through each version in order.

Given a config at version `0.1.0` and a manager at version `0.4.0`:

```
0.1.0 → 0.2.0  (e.g., bump provider version constraint)
0.2.0 → 0.3.0  (e.g., bump module tag in project.tf)
0.3.0 → 0.4.0  (e.g., bump Terraform version)
```

The manager reads `.terraform-config-manager-version` to determine the config's current version, then applies each migration in sequence until the config matches the manager's version. The version file is updated after each successful step.

### What migrations can change

Migrations can modify any platform-owned file:

- **`main.tf`** — Terraform version constraints, provider version constraints, provider configuration
- **`.terraform-version`** — pinned Terraform version for tfenv
- **`.terraform-config-manager-version`** — updated after each migration step

Migrations can also modify the project-owned file:

- **`project.tf`** — module `source` tags, module parameters (e.g., adding new required variables, renaming parameters, or removing deprecated ones to match breaking changes in module versions)

These are platform-initiated changes to a project-owned file. Migrations do not add or remove module blocks — they only modify existing ones to stay compatible with updated module versions.

## Usage

```sh
terraform-config-manager migrate [-dir=<path>] [-target-version=<version>]
```

- `-dir` — path to the project's Terraform configuration directory (defaults to current directory)
- `-target-version` — version to migrate to (defaults to the manager's current version)

The command reports each migration step as it is applied and exits with an error if any step fails, leaving the config at the last successfully applied version.

## Migration Definitions

Each migration is defined in code as a function that takes the config directory and applies a specific set of changes using `hclwrite` to parse and surgically edit the HCL files. Migrations are registered in version order and are immutable once released — a published migration is never modified, only new ones are appended.

## Example

A migration from `0.1.0` to `0.2.0` that bumps the Google provider version constraint:

1. Parse `main.tf`
2. Find the `required_providers` block
3. Update the `google` provider's `version` attribute from `~> 7.0` to `~> 8.0`
4. Write the modified `main.tf`

The framework handles updating `.terraform-config-manager-version` automatically after each successful migration step.
