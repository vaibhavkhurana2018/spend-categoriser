package extractor

import (
	"testing"
)

func TestExtractTransactions_PrimaryDateFormat(t *testing.T) {
	input := "15/01/2025 10:30:00             ACME STORE               MUMBAI                                                                           100                       1,250.00"
	txns, _ := ExtractTransactions(input)
	if len(txns) != 1 {
		t.Fatalf("got %d transactions, want 1", len(txns))
	}
	if txns[0].Date != "15/01/2025" {
		t.Errorf("date = %q, want 15/01/2025", txns[0].Date)
	}
	if txns[0].Amount != 1250.0 {
		t.Errorf("amount = %f, want 1250.00", txns[0].Amount)
	}
}

func TestExtractTransactions_DateWithoutTime(t *testing.T) {
	input := "20/01/2025                      ACME STORE               MUMBAI                                                                           - 10                      500.00 Cr"
	txns, _ := ExtractTransactions(input)
	if len(txns) != 1 {
		t.Fatalf("got %d transactions, want 1", len(txns))
	}
	if txns[0].Amount >= 0 {
		t.Errorf("credit should be negative, got %f", txns[0].Amount)
	}
}

func TestExtractTransactions_DashDateFormat(t *testing.T) {
	input := "15-01-2025 ACME STORE MUMBAI 1,250.00"
	txns, _ := ExtractTransactions(input)
	if len(txns) != 1 {
		t.Fatalf("got %d transactions, want 1", len(txns))
	}
	if txns[0].Date != "15-01-2025" {
		t.Errorf("date = %q, want 15-01-2025", txns[0].Date)
	}
}

func TestExtractTransactions_CreditEntry(t *testing.T) {
	input := "05/02/2025 09:15:00         NEFT CR 000011112222 000011112222A (Ref# 00000000000000000000000)                                                              10,000.00 Cr"
	txns, _ := ExtractTransactions(input)
	if len(txns) != 1 {
		t.Fatalf("got %d transactions, want 1", len(txns))
	}
	if txns[0].Amount != -10000.0 {
		t.Errorf("amount = %f, want -10000.00", txns[0].Amount)
	}
}

func TestExtractTransactions_SkipsHeaders(t *testing.T) {
	input := `Domestic Transactions
          Date                  Transaction Description                                                                            Amount (in Rs.)
15/01/2025 10:30:00             ACME STORE               MUMBAI                                                                           100                       1,250.00
Page 1 of 4
Opening Balance 50,000.00`
	txns, _ := ExtractTransactions(input)
	if len(txns) != 1 {
		t.Fatalf("got %d transactions, want 1 (only ACME line)", len(txns))
	}
}

func TestExtractTransactions_EmptyInput(t *testing.T) {
	txns, period := ExtractTransactions("")
	if len(txns) != 0 {
		t.Errorf("got %d transactions for empty input, want 0", len(txns))
	}
	if period.From != "" || period.To != "" {
		t.Errorf("got non-empty period for empty input: %+v", period)
	}
}

func TestExtractTransactions_MultipleTransactions(t *testing.T) {
	input := `15/01/2025 10:30:00             ACME STORE               MUMBAI                                                                           100                       1,250.00
15/01/2025 14:00:00             GLOBEX MART DELHI                                                                                                                     300.00
20/01/2025                      ACME STORE               MUMBAI                                                                           - 10                      500.00 Cr`
	txns, _ := ExtractTransactions(input)
	if len(txns) != 3 {
		t.Fatalf("got %d transactions, want 3", len(txns))
	}
}

func TestExtractTransactions_BillPeriod(t *testing.T) {
	input := "Statement period: 01/01/2025 to 31/01/2025\n15/01/2025 ACME STORE MUMBAI 100.00"
	_, period := ExtractTransactions(input)
	if period.From != "01/01/2025" || period.To != "31/01/2025" {
		t.Errorf("period = %+v, want {01/01/2025, 31/01/2025}", period)
	}
}

func TestMatchTransaction_PrimaryPattern(t *testing.T) {
	line := "15/01/2025 10:30:00             ACME STORE               MUMBAI                                                                           100                       1,250.00"
	date, txnTime, desc, amount, isCredit := matchTransaction(line)
	if date != "15/01/2025" {
		t.Errorf("date = %q, want 15/01/2025", date)
	}
	if txnTime != "10:30:00" {
		t.Errorf("time = %q, want 10:30:00", txnTime)
	}
	if amount != 1250.0 {
		t.Errorf("amount = %f, want 1250.0", amount)
	}
	if isCredit {
		t.Error("should not be credit")
	}
	if desc == "" {
		t.Error("description should not be empty")
	}
}

func TestMatchTransaction_NoMatch(t *testing.T) {
	date, _, _, _, _ := matchTransaction("This is not a transaction line")
	if date != "" {
		t.Errorf("expected no match, got date=%q", date)
	}
}

func TestParseMatch_CommaRemoval(t *testing.T) {
	_, _, _, amount, _ := parseMatch("15/01/2025", "10:00:00", "DESC", "1,23,456.78", "")
	if amount != 123456.78 {
		t.Errorf("amount = %f, want 123456.78", amount)
	}
}

func TestParseMatch_CreditDetection(t *testing.T) {
	_, _, _, _, isCredit := parseMatch("15/01/2025", "", "DESC", "100.00", " Cr")
	if !isCredit {
		t.Error("expected credit detection")
	}
}

func TestParseMatch_NonCredit(t *testing.T) {
	_, _, _, _, isCredit := parseMatch("15/01/2025", "", "DESC", "100.00", "")
	if isCredit {
		t.Error("should not be credit")
	}
}

func TestShouldSkip_MatchingLines(t *testing.T) {
	lines := []string{
		"Page 1 of 4",
		"Opening Balance 50,000.00",
		"Total Dues 75,000.00",
		"Minimum Amount Due 3,750.00",
		"Credit Limit 5,00,000",
		"Domestic Transactions",
		"Reward Points Summary",
		"Amount (in Rs.)",
		"Important Information",
	}
	for _, l := range lines {
		if !shouldSkip(l) {
			t.Errorf("shouldSkip(%q) = false, want true", l)
		}
	}
}

func TestShouldSkip_TransactionLine(t *testing.T) {
	line := "15/01/2025 10:30:00             ACME STORE               MUMBAI                                                                           100                       1,250.00"
	if shouldSkip(line) {
		t.Error("transaction line should not be skipped")
	}
}

func TestCleanDescription_RemovesRewardPoints(t *testing.T) {
	input := "ACME STORE               MUMBAI                                                                           100"
	got := cleanDescription(input)
	if got != "ACME STORE MUMBAI" {
		t.Errorf("got %q, want %q", got, "ACME STORE MUMBAI")
	}
}

func TestCleanDescription_RemovesNegativeRewardPoints(t *testing.T) {
	input := "ACME STORE               MUMBAI                                                                           - 10"
	got := cleanDescription(input)
	if got != "ACME STORE MUMBAI" {
		t.Errorf("got %q, want %q", got, "ACME STORE MUMBAI")
	}
}

func TestCleanDescription_CollapsesWhitespace(t *testing.T) {
	input := "GLOBEX MART DELHI"
	got := cleanDescription(input)
	if got != "GLOBEX MART DELHI" {
		t.Errorf("got %q, want %q", got, "GLOBEX MART DELHI")
	}
}

func TestExtractPeriod_StandardFormat(t *testing.T) {
	text := "Statement period: 01/01/2025 to 31/01/2025"
	period := extractPeriod(text)
	if period.From != "01/01/2025" || period.To != "31/01/2025" {
		t.Errorf("got %+v, want {01/01/2025, 31/01/2025}", period)
	}
}

func TestExtractPeriod_FromToFormat(t *testing.T) {
	text := "From 01/01/2025 to 31/01/2025"
	period := extractPeriod(text)
	if period.From != "01/01/2025" || period.To != "31/01/2025" {
		t.Errorf("got %+v, want {01/01/2025, 31/01/2025}", period)
	}
}

func TestExtractPeriod_NoMatch(t *testing.T) {
	period := extractPeriod("No period information here")
	if period.From != "" || period.To != "" {
		t.Errorf("got %+v, want empty", period)
	}
}

func TestGenerateHash_Deterministic(t *testing.T) {
	h1 := generateHash("15/01/2025", "10:30:00", "ACME STORE", 1250.00)
	h2 := generateHash("15/01/2025", "10:30:00", "ACME STORE", 1250.00)
	if h1 != h2 {
		t.Errorf("hashes differ for identical inputs: %q vs %q", h1, h2)
	}
}

func TestGenerateHash_DifferentInputs(t *testing.T) {
	h1 := generateHash("15/01/2025", "10:30:00", "ACME STORE", 1250.00)
	h2 := generateHash("15/01/2025", "10:30:00", "GLOBEX MART", 1250.00)
	h3 := generateHash("16/01/2025", "10:30:00", "ACME STORE", 1250.00)
	h4 := generateHash("15/01/2025", "10:30:00", "ACME STORE", 999.00)
	if h1 == h2 || h1 == h3 || h1 == h4 {
		t.Error("different inputs should produce different hashes")
	}
}

func TestGenerateHash_DifferentTimes(t *testing.T) {
	h1 := generateHash("15/01/2025", "10:30:00", "ACME STORE", 1250.00)
	h2 := generateHash("15/01/2025", "14:00:00", "ACME STORE", 1250.00)
	if h1 == h2 {
		t.Error("same date/description/amount but different times should produce different hashes")
	}
}

func TestGenerateHash_CaseNormalization(t *testing.T) {
	h1 := generateHash("15/01/2025", "10:30:00", "ACME STORE", 100.0)
	h2 := generateHash("15/01/2025", "10:30:00", "acme store", 100.0)
	if h1 != h2 {
		t.Error("hash should be case-insensitive")
	}
}
