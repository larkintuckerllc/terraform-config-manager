package hclutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func ExtractStringValue(attr *hclwrite.Attribute) string {
	tokens := attr.Expr().BuildTokens(nil)
	var sb strings.Builder
	for _, token := range tokens {
		if token.Type == hclsyntax.TokenQuotedLit {
			sb.Write(token.Bytes)
		}
	}
	return sb.String()
}

func UpdateModuleTag(dir, oldTag, newTag string) error {
	path := filepath.Join(dir, "project.tf")
	src, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading project.tf: %w", err)
	}

	f, diags := hclwrite.ParseConfig(src, path, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return fmt.Errorf("parsing project.tf: %s", diags.Error())
	}

	for _, block := range f.Body().Blocks() {
		if block.Type() != "module" {
			continue
		}
		sourceAttr := block.Body().GetAttribute("source")
		if sourceAttr == nil {
			continue
		}
		source := ExtractStringValue(sourceAttr)
		if !strings.Contains(source, "?ref=") {
			continue
		}
		if !strings.Contains(source, "?ref="+oldTag) {
			labels := block.Labels()
			name := "unknown"
			if len(labels) > 0 {
				name = labels[0]
			}
			return fmt.Errorf("module %q has unexpected tag in source: %s (expected %s)", name, source, oldTag)
		}
		updated := strings.Replace(source, "?ref="+oldTag, "?ref="+newTag, 1)
		block.Body().SetAttributeValue("source", cty.StringVal(updated))
	}

	return os.WriteFile(path, f.Bytes(), 0644)
}
