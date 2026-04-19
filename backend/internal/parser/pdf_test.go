package parser

import (
	"os/exec"
	"testing"
)

func requirePdftotext(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("pdftotext"); err != nil {
		t.Skip("pdftotext not available, skipping")
	}
}

func TestExtractText_InvalidPDF(t *testing.T) {
	requirePdftotext(t)

	_, err := ExtractText([]byte("this is not a pdf"))
	if err == nil {
		t.Error("expected error for invalid PDF input")
	}
}

func TestExtractText_EmptyInput(t *testing.T) {
	requirePdftotext(t)

	_, err := ExtractText([]byte{})
	if err == nil {
		t.Error("expected error for empty input")
	}
}
