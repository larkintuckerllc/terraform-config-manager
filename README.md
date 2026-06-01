# terraform-config-manager

A Go application for programmatically managing Terraform configurations using HashiCorp's official [hcl/v2](https://github.com/hashicorp/hcl) library.

## Purpose

Terraform configurations are typically written and maintained by hand. This works well for small projects, but as infrastructure grows, teams often face repetitive scaffolding, inconsistent conventions, and tedious bulk updates across many `.tf` files.

`terraform-config-manager` addresses this by providing programmatic control over Terraform configurations — generating, modifying, validating, and refactoring HCL files through Go code rather than manual editing.

## Current Scope

The initial implementation targets a deliberately narrow use case to keep things simple:

- **Provider:** Google Cloud
- **State:** Local (no remote backend)
- **Input:** GCP project ID and project owner (GitHub team or user)

This minimal scope lets us explore the `hclwrite` API without getting lost in infrastructure complexity.

The generated configuration is split into two files to support distinct ownership:

- **`main.tf`** — owned by the terraform-config-manager team (platform config: Terraform version, providers)
- **`project.tf`** — owned by the project owner (module calls for resources)

This separation is enforced via a generated `.github/CODEOWNERS` file and enables targeted validation of project-owner changes.

## Prerequisites

- **Go** — to build the application
- **tfenv** — must be installed and on PATH; the scaffolder checks for it and writes a `.terraform-version` file into each generated project

## Use Cases

- **[Scaffold new configurations](docs/scaffolding.md)** — Generate a complete Terraform project from a GCP project ID and project owner: provider setup, module references to a [shared Terraform modules repo](https://github.com/larkintuckerllc/terraform-modules), CODEOWNERS, `.terraform-version`, `.terraform-config-manager-version`, and `.gitignore`.
- **[Validate project configurations](docs/validation.md)** — Ensure `project.tf` contains only approved module references pinned to valid tags, suitable for CI review builds.
- **[Migrate configurations](docs/migrations.md)** — Apply versioned migrations to bring configs up to date: provider bumps, Terraform version changes, module tag updates, and breaking module changes. Inspired by database migrations.

## Technology

- **Language:** Go
- **Core library:** [`github.com/hashicorp/hcl/v2`](https://pkg.go.dev/github.com/hashicorp/hcl/v2) — the same library Terraform itself uses to parse and generate HCL.
  - `hclwrite` — create and surgically edit HCL files

## Project Structure

```
├── cmd/
│   └── terraform-config-manager/   # CLI entry point
├── internal/
│   ├── scaffold/                    # Scaffolding logic
│   └── validate/                    # Validation logic
├── docs/                            # Detailed documentation
├── go.mod                           # Go module definition
└── README.md
```

## Getting Started

Build the binary:

```sh
go build -o terraform-config-manager ./cmd/terraform-config-manager
```

Scaffold a new Terraform project:

```sh
./terraform-config-manager scaffold -project=my-project-id -owner=@acme-org/my-team
```

This creates a `my-project-id/` directory in the current folder with all the Terraform files ready to init and apply.

## Testing Approach

This project deliberately does not include unit tests. The scaffolder generates output with minimal logic — tests would just assert the generated HCL matches a hardcoded expected string, mirroring the implementation. Any change to the output would require updating both the code and the test in lockstep, with the test catching nothing a manual run wouldn't.

The validator follows the same reasoning. Its rules are simple conditionals that map directly to the [validation doc](docs/validation.md). There is no complex logic where a change could accidentally break an unrelated rule.

The right question before writing a test: *could someone accidentally break this code in a way the test would catch but running the tool wouldn't?* For both the scaffolder and validator, the answer is no. Validation is done by running the tool against real Terraform configurations.

## Documentation

Detailed documentation for each use case is available in the [docs/](docs/) directory.
