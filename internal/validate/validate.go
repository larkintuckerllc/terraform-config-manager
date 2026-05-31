package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

const approvedSourcePrefix = "git::https://github.com/larkintuckerllc/terraform-modules.git//"

func Run(dir string) error {
	path := filepath.Join(dir, "project.tf")
	src, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading project.tf: %w", err)
	}

	f, diags := hclwrite.ParseConfig(src, path, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return fmt.Errorf("parsing project.tf: %s", diags.Error())
	}

	var errors []string

	for _, block := range f.Body().Blocks() {
		if block.Type() != "module" {
			errors = append(errors, fmt.Sprintf("block type %q is not allowed, only module blocks are permitted", block.Type()))
			continue
		}

		sourceAttr := block.Body().GetAttribute("source")
		if sourceAttr == nil {
			labels := block.Labels()
			name := "unknown"
			if len(labels) > 0 {
				name = labels[0]
			}
			errors = append(errors, fmt.Sprintf("module %q is missing a source attribute", name))
			continue
		}

		source := extractStringValue(sourceAttr)
		labels := block.Labels()
		name := "unknown"
		if len(labels) > 0 {
			name = labels[0]
		}

		if !strings.HasPrefix(source, approvedSourcePrefix) {
			errors = append(errors, fmt.Sprintf("module %q has an unapproved source: %s", name, source))
			continue
		}

		if !strings.Contains(source, "?ref=") {
			errors = append(errors, fmt.Sprintf("module %q source is not pinned to a tag", name))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation failed:\n  - %s", strings.Join(errors, "\n  - "))
	}

	fmt.Println("Validation passed")
	return nil
}

func extractStringValue(attr *hclwrite.Attribute) string {
	tokens := attr.Expr().BuildTokens(nil)
	var sb strings.Builder
	for _, token := range tokens {
		if token.Type == hclsyntax.TokenQuotedLit {
			sb.Write(token.Bytes)
		}
	}
	return sb.String()
}
