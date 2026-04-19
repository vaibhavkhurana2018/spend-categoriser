// Package parser converts raw PDF bytes into plain text using the external
// pdftotext binary (part of poppler-utils). The binary must be available on PATH.
package parser

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// ExtractText writes the PDF bytes to a temp file, invokes pdftotext with
// -layout mode to preserve column alignment, and returns the extracted text.
func ExtractText(data []byte) (string, error) {
	tmpFile, err := os.CreateTemp("", "spend-*.pdf")
	if err != nil {
		return "", fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("writing temp file: %w", err)
	}
	tmpFile.Close()

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("pdftotext", "-layout", tmpFile.Name(), "-")
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pdftotext: %s: %w", stderr.String(), err)
	}

	return out.String(), nil
}
