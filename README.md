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
- **Modify existing configurations** — Surgically add, update, or remove resources, variables, and outputs in existing `.tf` files while preserving formatting and comments.
- **Enforce standards** — Validate configurations against organizational conventions (naming, tagging, required provider versions) and optionally auto-fix violations.
- **Compose modules** — Programmatically wire together module calls with the correct variable bindings, producing ready-to-plan configurations.
- **Bulk refactoring** — Rename resources, move blocks between files, update provider versions, or migrate patterns across many `.tf` files at once.

## Technology

- **Language:** Go
- **Core library:** [`github.com/hashicorp/hcl/v2`](https://pkg.go.dev/github.com/hashicorp/hcl/v2) — the same library Terraform itself uses to parse and generate HCL.
  - `hclwrite` — create and surgically edit HCL files
  - `hclparse` — parse HCL into an AST
  - `gohcl` — marshal/unmarshal Go structs to/from HCL

## Project Structure

```
├── cmd/
│   └── terraform-config-manager/   # CLI entry point
├── internal/
│   └── scaffold/                    # Scaffolding logic
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

## Documentation

Detailed documentation for each use case is available in the [docs/](docs/) directory.
