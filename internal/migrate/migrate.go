package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const currentVersion = "0.1.0"

type migration struct {
	fromVersion string
	toVersion   string
	apply       func(dir string) error
}

var migrations = []migration{
	// To add a new migration:
	// {fromVersion: "0.1.0", toVersion: "0.2.0", apply: migrate010to020},
}

func Run(dir string, targetVersion string) error {
	if targetVersion == "" {
		targetVersion = currentVersion
	}

	configVersion, err := readVersion(dir)
	if err != nil {
		return err
	}

	if configVersion == targetVersion {
		fmt.Printf("Configuration is already at version %s\n", targetVersion)
		return nil
	}

	if !isKnownVersion(configVersion) {
		return fmt.Errorf("unknown configuration version: %s", configVersion)
	}

	if !isReachable(configVersion, targetVersion) {
		return fmt.Errorf("no migration path from %s to %s", configVersion, targetVersion)
	}

	for _, m := range migrations {
		if m.fromVersion < configVersion {
			continue
		}
		if m.fromVersion >= targetVersion {
			break
		}

		fmt.Printf("Migrating %s → %s\n", m.fromVersion, m.toVersion)
		if err := m.apply(dir); err != nil {
			return fmt.Errorf("migration %s → %s failed: %w", m.fromVersion, m.toVersion, err)
		}

		if err := writeVersion(dir, m.toVersion); err != nil {
			return fmt.Errorf("updating version after %s → %s: %w", m.fromVersion, m.toVersion, err)
		}
	}

	fmt.Printf("Configuration migrated to version %s\n", targetVersion)
	return nil
}

func readVersion(dir string) (string, error) {
	path := filepath.Join(dir, ".terraform-config-manager-version")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading .terraform-config-manager-version: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

func writeVersion(dir string, version string) error {
	path := filepath.Join(dir, ".terraform-config-manager-version")
	return os.WriteFile(path, []byte(version+"\n"), 0644)
}

func isKnownVersion(v string) bool {
	if v == currentVersion {
		return true
	}
	for _, m := range migrations {
		if m.fromVersion == v || m.toVersion == v {
			return true
		}
	}
	return false
}

func isReachable(from, to string) bool {
	if from == to {
		return true
	}
	v := from
	for _, m := range migrations {
		if m.fromVersion == v {
			v = m.toVersion
			if v == to {
				return true
			}
		}
	}
	return false
}
