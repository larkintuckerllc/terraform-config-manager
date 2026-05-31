# Scaffolding New Configurations

## Overview

Scaffolding is the process of generating new Terraform configuration files from parameters rather than writing them by hand. This is the starting point — producing the foundational files that a user can immediately `terraform init` and `terraform apply`.

## Problem

Setting up a new Terraform project involves creating several boilerplate files that follow a predictable structure. Teams often copy these from an existing project and manually adjust them, which leads to drift in conventions, forgotten required fields, and inconsistent structure across projects.

## Scope

The scaffolder takes a single input — a **GCP project ID** — and generates a minimal, working Terraform configuration for managing resources in that project.

Design decisions:

- **Output directory per project.** The scaffolder creates a directory named after the GCP project ID (e.g., `my-project-id/`) inside the working directory. All generated files are placed there. This gives each project its own isolated Terraform root, with its own state, and makes it natural to manage multiple projects side by side.
- **Local state only.** No remote backend configuration. This keeps the initial setup simple.
- **Module-based resources.** Resources are defined in reusable Terraform modules hosted in a separate [Git repository](https://github.com/larkintuckerllc/terraform-modules), pinned to a specific Git tag.
- **Two-file ownership model.** The generated Terraform configuration is split across two files to support distinct ownership via Git CODEOWNERS:
  - `main.tf` — owned by the terraform-config-manager team. Contains the `terraform` block, `required_providers`, and `provider` configuration. Project owners cannot modify this file.
  - `project.tf` — owned by the project owner. Contains only module calls. This is the only file project owners need to modify to manage their resources, and is the focus of future config validation.

## Prerequisites

- **tfenv** must be installed and available on PATH. The scaffolder checks for this before generating any files and exits with an error if it is not found.

## Usage

```sh
terraform-config-manager scaffold -project=<gcp-project-id> [-output-dir=<path>]
```

- `-project` — GCP project ID (required)
- `-output-dir` — directory to create the project folder in (defaults to current directory)

## Approach

Using `hclwrite`, we programmatically construct HCL files from the input parameters. The `hclwrite` package operates at the token level, producing clean, properly formatted HCL that looks hand-written.

Key `hclwrite` capabilities used:

- `hclwrite.NewEmptyFile()` — create a new HCL file from scratch
- `Body.AppendNewBlock()` — add `terraform`, `provider`, `resource`, and other block types
- `Body.SetAttributeValue()` — set simple attributes (strings, numbers, bools)
- `File.Bytes()` — emit the final formatted HCL output

## What Gets Generated

Given project ID `my-project-id`, the scaffolder creates `my-project-id/` containing:

```
my-project-id/
├── .gitignore
├── .terraform-config-manager-version
├── .terraform-version
├── main.tf
└── project.tf
```

### `.terraform-version`

Pins the Terraform version for `tfenv` to auto-switch when entering the directory.

```
1.15.5
```

### `.gitignore`

Excludes Terraform working files and state from version control.

```
.terraform/
*.tfstate
*.tfstate.backup
*.tfplan
.terraform.lock.hcl
```

### `.terraform-config-manager-version`

The version of `terraform-config-manager` that generated this configuration. Same format as `.terraform-version` — a single version string.

```
0.1.0
```

This enables detecting stale configurations by comparing the file's version against the running manager's version.

### `main.tf`

Platform configuration owned by the terraform-config-manager team: `terraform` block, `required_providers`, and `provider` setup. Project owners should not modify this file. Contents will evolve as the manager changes platform settings.

### `project.tf`

Project-specific module calls owned by the project owner. This is the only file project owners modify to manage their resources. Module references point to the [terraform-modules](https://github.com/larkintuckerllc/terraform-modules) repository, pinned to a specific Git tag. Contents will evolve as the manager adds support for new modules. See the generated file for current details.

## Validating the Output

After scaffolding, the generated configuration can be validated and applied:

```sh
cd my-project-id
terraform init
terraform validate
terraform plan
terraform apply
```

GCP authentication (e.g., `gcloud auth application-default login`) is required before `terraform plan` or `terraform apply`. When done testing, `terraform destroy` cleans up the created resources.
