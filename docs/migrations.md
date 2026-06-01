# Migrating Configurations

## Overview

When the platform team needs to update Terraform configurations — bumping provider versions, updating module tags, changing Terraform version constraints — the terraform-config-manager applies these changes as a series of versioned migrations, similar to database migrations.

## Problem

Terraform configurations managed by the scaffolder will drift from the current platform standards over time. A project scaffolded at version `0.1.0` may need provider version bumps, new module versions, or structural changes introduced in later versions. Applying all changes at once is fragile — each change may depend on the state left by the previous one, and skipping intermediate steps risks breaking configs.

## Versioning

The terraform-config-manager version and the terraform-modules version are independent, similar to how Helm separates chart version from app version:

- **terraform-config-manager version** (`0.1.0`, `0.2.0`, ...) — the version of the manager and the config structure it produces. Tracked in `.terraform-config-manager-version`.
- **terraform-modules tags** (`v0.1.0`, `v0.2.0`, ...) — the version of the modules themselves, on their own release cadence.

A manager migration may bump a module tag, but the version numbers are not aligned. For example, manager version `0.3.0` might reference module tag `v0.2.0` — the two evolve independently.

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

The framework automatically updates `.terraform-config-manager-version` after each successful migration step — individual migrations do not need to handle this.

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

The `0.1.0 → 0.2.0` migration bumps the module tag in `project.tf` from `v0.1.0` to `v0.2.0`, picking up a new configurable `location` variable added to the `my-bucket` module. Since the variable has a default value, no changes to module parameters are needed — only the tag is updated.

The migration function is a single call to the shared `UpdateModuleTag` utility in `internal/hclutil`, which parses `project.tf`, finds all module sources referencing the old tag, and replaces them with the new tag.

The framework handles updating `.terraform-config-manager-version` automatically after each successful migration step.
