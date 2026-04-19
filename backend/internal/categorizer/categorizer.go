// Package categorizer assigns a spend category and subcategory to a
// transaction description using keyword-based matching. Keywords are
// checked in priority order defined in mcc.go; the first match wins.
package categorizer

import (
	"strings"
)

// Categorizer performs keyword-based categorisation of transaction descriptions.
type Categorizer struct{}

// New returns a ready-to-use Categorizer.
func New() *Categorizer {
	return &Categorizer{}
}

// Categorize returns the (category, subCategory) for a transaction description.
// Matching is case-insensitive. Returns ("Other", "Uncategorized") when no rule matches.
func (c *Categorizer) Categorize(description string) (string, string) {
	lower := strings.ToLower(description)

	for _, rule := range keywordRules {
		for _, kw := range rule.keywords {
			if strings.Contains(lower, kw) {
				return rule.category, rule.subCategory
			}
		}
	}

	return "Other", "Uncategorized"
}

// GetCategories returns all known categories with their subcategories,
// aggregated from both MCC codes and keyword rules.
func GetCategories() []map[string]interface{} {
	seen := make(map[string]map[string]bool)

	for _, m := range mccCodes {
		if seen[m.Category] == nil {
			seen[m.Category] = make(map[string]bool)
		}
		seen[m.Category][m.SubCategory] = true
	}

	for _, rule := range keywordRules {
		if seen[rule.category] == nil {
			seen[rule.category] = make(map[string]bool)
		}
		seen[rule.category][rule.subCategory] = true
	}

	var result []map[string]interface{}
	for cat, subs := range seen {
		subList := make([]string, 0, len(subs))
		for sub := range subs {
			subList = append(subList, sub)
		}
		result = append(result, map[string]interface{}{
			"name":          cat,
			"subCategories": subList,
		})
	}
	return result
}
