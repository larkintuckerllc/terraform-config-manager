package migrate

import "terraform-config-manager/internal/hclutil"

func migrate010to020(dir string) error {
	return hclutil.UpdateModuleTag(dir, "v0.1.0", "v0.2.0")
}
