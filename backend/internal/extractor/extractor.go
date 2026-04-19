// Package extractor parses plain-text output from pdftotext into structured
// Transaction records. It handles multiple date formats and the column-based
// layout typical of Indian bank credit card statements (e.g. HDFC).
package extractor

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/spend-categoriser/backend/internal/categorizer"
	"github.com/spend-categoriser/backend/internal/models"
)

// txnPattern matches the primary HDFC format: DD/MM/YYYY [HH:MM:SS] <description> <amount> [Cr]
var txnPattern = regexp.MustCompile(
	`^\s*(\d{2}/\d{2}/\d{4})` +
		`(?:\s+(\d{2}:\d{2}:\d{2}))?` +
		`\s+(.+?)` +
		`\s+([\d,]+\.\d{2})` +
		`(\s+Cr)?\s*$`,
)

// altDatePatterns cover DD-MM-YYYY, DD Mon YYYY, and DD/MM/YY formats.
var altDatePatterns = []*regexp.Regexp{
	regexp.MustCompile(`^\s*(\d{2}-\d{2}-\d{4})(?:\s+(\d{2}:\d{2}:\d{2}))?\s+(.+?)\s+([\d,]+\.\d{2})(\s+Cr)?\s*$`),
	regexp.MustCompile(`^\s*(\d{2}\s+\w{3}\s+\d{4})(?:\s+(\d{2}:\d{2}:\d{2}))?\s+(.+?)\s+([\d,]+\.\d{2})(\s+Cr)?\s*$`),
	regexp.MustCompile(`^\s*(\d{2}/\d{2}/\d{2})(?:\s+(\d{2}:\d{2}:\d{2}))?\s+(.+?)\s+([\d,]+\.\d{2})(\s+Cr)?\s*$`),
}

// skipPatterns identify header, footer, and summary lines that should not be parsed as transactions.
var skipPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)page\s+\d+\s+of\s+\d+`),
	regexp.MustCompile(`(?i)statement\s+(for|date|period)`),
	regexp.MustCompile(`(?i)^date\s+transaction`),
	regexp.MustCompile(`(?i)opening\s+balance`),
	regexp.MustCompile(`(?i)closing\s+balance`),
	regexp.MustCompile(`(?i)total\s+dues`),
	regexp.MustCompile(`(?i)minimum\s+amount`),
	regexp.MustCompile(`(?i)credit\s+limit`),
	regexp.MustCompile(`(?i)reward\s+points?\s+summary`),
	regexp.MustCompile(`(?i)feature\s+reward`),
	regexp.MustCompile(`(?i)amount\s+\(in\s+rs`),
	regexp.MustCompile(`(?i)domestic\s+transactions`),
	regexp.MustCompile(`(?i)international\s+transactions`),
	regexp.MustCompile(`(?i)payment\s+due\s+date`),
	regexp.MustCompile(`(?i)account\s+summary`),
	regexp.MustCompile(`(?i)past\s+dues`),
	regexp.MustCompile(`(?i)important\s+information`),
	regexp.MustCompile(`(?i)overlimit`),
}

// periodPatterns extract the billing date range from statement header text.
var periodPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)statement\s+period[:\s]+(\d{2}/\d{2}/\d{4})\s*(?:to|-)\s*(\d{2}/\d{2}/\d{4})`),
	regexp.MustCompile(`(?i)from\s+(\d{2}/\d{2}/\d{4})\s*to\s*(\d{2}/\d{2}/\d{4})`),
	regexp.MustCompile(`(?i)billing\s+period[:\s]+(\d{2}\s+\w{3}\s+\d{4})\s*(?:to|-)\s*(\d{2}\s+\w{3}\s+\d{4})`),
	regexp.MustCompile(`(?i)statement\s+date[:\s]+(\d{2}/\d{2}/\d{4})`),
}

// ExtractTransactions scans each line of the pdftotext output, identifies
// transaction rows, categorises them, and returns the results along with
// the billing period found in the statement header.
func ExtractTransactions(text string) ([]models.Transaction, models.BillPeriod) {
	lines := strings.Split(text, "\n")
	var transactions []models.Transaction
	cat := categorizer.New()

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		if shouldSkip(line) {
			continue
		}

		date, txnTime, desc, amount, isCredit := matchTransaction(line)
		if date == "" {
			continue
		}

		desc = cleanDescription(desc)
		if desc == "" || amount == 0 {
			continue
		}

		if isCredit {
			amount = -amount
		}

		category, subCategory := cat.Categorize(desc)
		hash := generateHash(date, txnTime, desc, amount)

		transactions = append(transactions, models.Transaction{
			Date:        date,
			Description: desc,
			Amount:      amount,
			Category:    category,
			SubCategory: subCategory,
			Hash:        hash,
		})
	}

	period := extractPeriod(text)
	return transactions, period
}

// matchTransaction tries the primary pattern first, then falls back to
// alternative date formats. Returns zero values if the line is not a transaction.
func matchTransaction(line string) (date, txnTime, desc string, amount float64, isCredit bool) {
	matches := txnPattern.FindStringSubmatch(line)
	if len(matches) >= 6 {
		return parseMatch(matches[1], matches[2], matches[3], matches[4], matches[5])
	}

	for _, pat := range altDatePatterns {
		matches = pat.FindStringSubmatch(line)
		if len(matches) >= 6 {
			return parseMatch(matches[1], matches[2], matches[3], matches[4], matches[5])
		}
	}

	return "", "", "", 0, false
}

// parseMatch converts the raw regex capture groups into typed values.
// Commas are stripped from amounts and the "Cr" suffix signals a credit.
func parseMatch(dateStr, timeStr, descStr, amountStr, crStr string) (string, string, string, float64, bool) {
	amountStr = strings.ReplaceAll(amountStr, ",", "")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return "", "", "", 0, false
	}
	isCredit := strings.Contains(crStr, "Cr")
	return dateStr, timeStr, descStr, amount, isCredit
}

// shouldSkip returns true for header, footer, and summary lines that look
// like transactions (they may contain dates and amounts) but are not.
func shouldSkip(line string) bool {
	trimmed := strings.TrimSpace(line)
	for _, pat := range skipPatterns {
		if pat.MatchString(trimmed) {
			return true
		}
	}
	return false
}

// rewardPointsInDesc strips trailing reward-point numbers (e.g. "  136") that
// appear between the description and amount in HDFC layout-mode output.
var rewardPointsInDesc = regexp.MustCompile(`\s{2,}-?\s*\d+\s*$`)
var multiSpaces = regexp.MustCompile(`\s{2,}`)

// cleanDescription removes trailing reward-point columns and collapses
// runs of whitespace that result from the PDF's fixed-width layout.
func cleanDescription(desc string) string {
	desc = rewardPointsInDesc.ReplaceAllString(desc, "")
	desc = multiSpaces.ReplaceAllString(desc, " ")
	desc = strings.TrimSpace(desc)
	return desc
}

// extractPeriod searches the full text for a billing-period header and
// returns the date range. Returns an empty BillPeriod if none found.
func extractPeriod(text string) models.BillPeriod {
	for _, pattern := range periodPatterns {
		matches := pattern.FindStringSubmatch(text)
		if len(matches) >= 3 {
			return models.BillPeriod{From: matches[1], To: matches[2]}
		}
		if len(matches) >= 2 {
			return models.BillPeriod{From: matches[1], To: matches[1]}
		}
	}
	return models.BillPeriod{}
}

// generateHash produces a deterministic SHA-256 hex string from the transaction's
// date, time, description, and amount. The time component distinguishes
// transactions that share the same date, description, and amount but occurred
// at different times.
func generateHash(date, txnTime, description string, amount float64) string {
	data := fmt.Sprintf("%s|%s|%s|%.2f", date, txnTime, strings.ToLower(strings.TrimSpace(description)), amount)
	h := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", h)
}
