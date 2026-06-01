package scaffold

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

const version = "0.2.0"

const terraformVersion = "1.15.5"

const platformTeam = "@acme-org/platform-team"

const gitignoreContent = `.terraform/
*.tfstate
*.tfstate.backup
*.tfplan
.terraform.lock.hcl
`

func Run(projectID string, projectOwner string, outputDir string) error {
	if _, err := exec.LookPath("tfenv"); err != nil {
		return fmt.Errorf("tfenv is not installed or not in PATH")
	}

	dir := filepath.Join(outputDir, projectID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dir, ".terraform-version"), []byte(terraformVersion+"\n"), 0644); err != nil {
		return fmt.Errorf("writing .terraform-version: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("writing .gitignore: %w", err)
	}

	if err := os.WriteFile(filepath.Join(dir, ".terraform-config-manager-version"), []byte(version+"\n"), 0644); err != nil {
		return fmt.Errorf("writing .terraform-config-manager-version: %w", err)
	}

	githubDir := filepath.Join(dir, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		return fmt.Errorf("creating .github directory: %w", err)
	}

	codeowners := generateCodeowners(projectOwner)
	if err := os.WriteFile(filepath.Join(githubDir, "CODEOWNERS"), codeowners, 0644); err != nil {
		return fmt.Errorf("writing CODEOWNERS: %w", err)
	}

	mainTF := generateMain(projectID)
	if err := os.WriteFile(filepath.Join(dir, "main.tf"), mainTF, 0644); err != nil {
		return fmt.Errorf("writing main.tf: %w", err)
	}

	bucketName := generateBucketName(projectID)

	projectTF := generateProject(bucketName)
	if err := os.WriteFile(filepath.Join(dir, "project.tf"), projectTF, 0644); err != nil {
		return fmt.Errorf("writing project.tf: %w", err)
	}

	fmt.Printf("Scaffolded Terraform configuration in %s\n", dir)
	return nil
}

func generateMain(projectID string) []byte {
	f := hclwrite.NewEmptyFile()
	body := f.Body()

	tfBlock := body.AppendNewBlock("terraform", nil)
	tfBody := tfBlock.Body()
	tfBody.SetAttributeValue("required_version", cty.StringVal("~> 1.15"))

	tfBody.AppendNewline()

	rpBlock := tfBody.AppendNewBlock("required_providers", nil)
	rpBody := rpBlock.Body()
	rpBody.SetAttributeValue("google", cty.ObjectVal(map[string]cty.Value{
		"source":  cty.StringVal("hashicorp/google"),
		"version": cty.StringVal("~> 7.0"),
	}))

	body.AppendNewline()

	providerBlock := body.AppendNewBlock("provider", []string{"google"})
	providerBody := providerBlock.Body()
	providerBody.SetAttributeValue("project", cty.StringVal(projectID))

	return f.Bytes()
}

func generateProject(bucketName string) []byte {
	f := hclwrite.NewEmptyFile()
	body := f.Body()

	moduleBlock := body.AppendNewBlock("module", []string{"my_bucket"})
	moduleBody := moduleBlock.Body()
	moduleBody.SetAttributeValue("source", cty.StringVal("git::https://github.com/larkintuckerllc/terraform-modules.git//my-bucket?ref=v0.2.0"))
	moduleBody.SetAttributeValue("bucket_name", cty.StringVal(bucketName))

	return f.Bytes()
}

func generateCodeowners(projectOwner string) []byte {
	return []byte(fmt.Sprintf("*          %s\nproject.tf %s\n", platformTeam, projectOwner))
}

func generateBucketName(projectID string) string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%s-%x", projectID, b)
}
