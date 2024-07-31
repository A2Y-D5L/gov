package gov

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func use(cmd *cobra.Command, args []string) {
	version := args[0]
	rc, err := cmd.Flags().GetBool("rc")
	if err != nil {
		log.Fatalf("Failed to parse 'rc' flag: %v", err)
	}

	if err := validateVersion(version, rc); err != nil {
		log.Fatalf("Invalid version: %v", err)
	}

	currentVersion, err := getCurrentGoVersion()
	if err == nil {
		if err := os.Setenv("GO_ROLLBACK_VERSION", currentVersion); err != nil {
			log.Fatalf("Failed to set rollback version: %v", err)
		}
	} else {
		fmt.Println("No current Go version installed.")
	}

	if err := removeExistingGo(); err != nil {
		log.Fatalf("Error removing existing Go installation: %v", err)
	}

	if err := installGoVersion(version, rc); err != nil {
		log.Fatalf("Error installing Go version: %v", err)
	}

	if err := updateEnvVariables(version); err != nil {
		log.Fatalf("Error updating environment variables: %v", err)
	}

	fmt.Println("Switched to Go version", version)
}

func rollback(cmd *cobra.Command, args []string) {
	rollbackVersion := os.Getenv("GO_ROLLBACK_VERSION")
	if rollbackVersion == "" {
		fmt.Println("No rollback version set. Use 'gov use latest' to install the latest stable version.")
		return
	}

	if err := installGoVersion(rollbackVersion, false); err != nil {
		log.Fatalf("Error installing rollback Go version: %v", err)
	}

	if err := updateEnvVariables(rollbackVersion); err != nil {
		log.Fatalf("Error updating environment variables: %v", err)
	}

	fmt.Println("Rolled back to Go version", rollbackVersion)
}

func validateVersion(version string, rc bool) error {
	// Basic validation logic for Go version
	if rc && !strings.HasPrefix(version, "1.") {
		return fmt.Errorf("release candidate versions must start with '1.'")
	}
	return nil
}

func getCurrentGoVersion() (string, error) {
	// Get the currently installed Go version
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	parts := strings.Fields(string(output))
	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected output format from 'go version' command")
	}

	return parts[2], nil
}

func removeExistingGo() error {
	// Logic to remove existing Go installation
	goRoot := os.Getenv("GOROOT")
	if goRoot == "" {
		goRoot = "/usr/local/go"
	}

	err := os.RemoveAll(goRoot)
	if err != nil {
		return fmt.Errorf("failed to remove existing Go installation: %w", err)
	}

	return nil
}

func installGoVersion(version string, rc bool) error {
	if rc {
		version += "rc"
	}

	url := fmt.Sprintf("https://golang.org/dl/go%s.linux-amd64.tar.gz", version)
	cmd := exec.Command("curl", "-LO", url)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to download Go version %s: %w", version, err)
	}

	cmd = exec.Command("tar", "-C", "/usr/local", "-xzf", fmt.Sprintf("go%s.linux-amd64.tar.gz", version))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract Go version %s: %w", version, err)
	}

	return nil
}

func updateEnvVariables(version string) error {
	goRoot := "/usr/local/go"
	goBin := filepath.Join(goRoot, "bin")

	err := os.Setenv("GOROOT", goRoot)
	if err != nil {
		return fmt.Errorf("failed to set GOROOT: %w", err)
	}

	err = os.Setenv("PATH", fmt.Sprintf("%s:%s", goBin, os.Getenv("PATH")))
	if err != nil {
		return fmt.Errorf("failed to update PATH: %w", err)
	}

	return nil
}

func main() {
	useCmd := &cobra.Command{
		Use:   "use [version]",
		Short: "Switches to a specified version of Go, setting it as the active version",
		Args:  cobra.ExactArgs(1),
		Run: use,
	}
	useCmd.Flags().Bool("rc", false, "Install a release candidate version if available")
	rollbackCmd := &cobra.Command{
		Use:   "rollback",
		Short: "Rolls back to the previously installed version of Go",
		Args:  cobra.ExactArgs(1),
		Run: rollback,
	}

	if err := &cobra.Command{
		Use:   "gov",
		Short: "gov is a CLI utility to manage Go installations",
		Long:  `A CLI utility to switch between different versions of Go and manage installations.`,
		Run:   func(cmd *cobra.Command, args []string) { cmd.AddCommand(useCmd, rollbackCmd) },
	}.Execute(); err != nil {
		log.Fatalf("Error executing root command: %v", err)
	}
}
