package wallet

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/0x3639/znn_cli_go/cmd"
	"github.com/0x3639/znn_cli_go/pkg/format"
	"github.com/spf13/cobra"
)

// exportCmd exports a wallet to a file
var exportCmd = &cobra.Command{
	Use:   "export <filePath>",
	Short: "Export a wallet to a file",
	Long: `Export a wallet keyStore file to the specified path.
This creates a copy of the encrypted keyStore file that can be backed up
or transferred to another system.

Example:
  znn-cli wallet export ./my-wallet-backup.json

Requires --keyStore flag to specify which wallet to export.`,
	Args: cobra.ExactArgs(1),
	RunE: runExport,
}

func init() {
	walletCmd.AddCommand(exportCmd)
}

func runExport(c *cobra.Command, args []string) error {
	destPath := args[0]

	cfg := cmd.GetConfig()
	keystoreName := cmd.GetKeyStore()

	if keystoreName == "" {
		return fmt.Errorf("--keyStore flag is required to specify which wallet to export")
	}

	// Get the source keyStore file path
	sourcePath := filepath.Join(cfg.Wallet.WalletDir, keystoreName)

	// Check if source exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("keyStore %s not found", keystoreName)
	}

	// Open source file
	// #nosec G304 - Source path is constrained to wallet directory
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() { _ = sourceFile.Close() }()

	// Create destination file
	// #nosec G304 - Destination path is user-specified (expected CLI behavior)
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() { _ = destFile.Close() }()

	// Copy file
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	format.Success(fmt.Sprintf("keyStore exported to: %s", destPath))
	fmt.Println(format.Yellow("⚠️  Keep this file safe! It contains your encrypted wallet."))

	return nil
}
