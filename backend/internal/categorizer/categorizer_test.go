package categorizer

import (
	"testing"
)

func TestCategorize_KnownKeywords(t *testing.T) {
	c := New()

	cases := []struct {
		desc       string
		wantCat    string
		wantSubCat string
	}{
		{"SWIGGY BANGALORE", "Food & Dining", "Delivery"},
		{"RAZ*SWIGGY TECHNOLOGIE http://www.s", "Food & Dining", "Delivery"},
		{"BUNDL TECHNOLOGIES PVT Bengaluru", "Food & Dining", "Delivery"},
		{"ZOMATO NEW DELHI", "Food & Dining", "Delivery"},
		{"AMAZON PAY INDIA PRIVA Bangalore", "Shopping", "Online Shopping"},
		{"FLIPKART PAYMENTS BANGALORE", "Shopping", "Online Shopping"},
		{"PAY*MYNTRA VIA SMARTBU GURGAON", "Shopping", "Online Shopping"},
		{"UBER INDIA BANGALORE", "Travel", "Ride Sharing"},
		{"Booking.com40552749257 ONLINE", "Travel", "Hotels"},
		{"AGODA.COM PG USM-HKT Internet", "Travel", "Hotels"},
		{"NETFLIX.COM", "Entertainment", "Streaming"},
		{"APOLLO PHARMACY BANGALORE", "Healthcare", "Pharmacy"},
		{"POLICY BAZAAR CYBS S GURGAON", "Insurance", "General"},
		{"CHAAYOS VEGA CITY MALL BANGALORE", "Food & Dining", "Cafes"},
		{"ADITYA BIRLA FASHION A BANGALORE", "Shopping", "Clothing"},
		{"LIFE STYLE INTERNATIONABANGALORE", "Shopping", "Clothing"},
		{"IMPS PMT 999900001234", "Financial", "Transfer"},
		{"ANAND SWEETS AND SAVOURBANGALORE", "Grocery", "Sweets & Bakery"},
		{"DEFMACRO SOFTWARE", "Shopping", "Software & Subscriptions"},
		{"GYFTR VIA SMARTBUY GURGAON", "Shopping", "Online Shopping"},
	}

	for _, tc := range cases {
		cat, sub := c.Categorize(tc.desc)
		if cat != tc.wantCat || sub != tc.wantSubCat {
			t.Errorf("Categorize(%q) = (%q, %q), want (%q, %q)", tc.desc, cat, sub, tc.wantCat, tc.wantSubCat)
		}
	}
}

func TestCategorize_UnknownDescription(t *testing.T) {
	c := New()
	cat, sub := c.Categorize("XYZZY UNKNOWN MERCHANT 12345")
	if cat != "Other" || sub != "Uncategorized" {
		t.Errorf("got (%q, %q), want (Other, Uncategorized)", cat, sub)
	}
}

func TestCategorize_CaseInsensitivity(t *testing.T) {
	c := New()

	variants := []string{"SWIGGY", "swiggy", "Swiggy", "sWiGgY"}
	for _, v := range variants {
		cat, sub := c.Categorize(v)
		if cat != "Food & Dining" || sub != "Delivery" {
			t.Errorf("Categorize(%q) = (%q, %q), want (Food & Dining, Delivery)", v, cat, sub)
		}
	}
}

func TestGetCategories_NonEmpty(t *testing.T) {
	cats := GetCategories()
	if len(cats) == 0 {
		t.Fatal("GetCategories returned empty slice")
	}
}

func TestGetCategories_ContainsExpectedNames(t *testing.T) {
	cats := GetCategories()
	expected := map[string]bool{
		"Grocery": false, "Food & Dining": false, "Travel": false,
		"Shopping": false, "Healthcare": false, "Financial": false,
	}

	for _, c := range cats {
		name, ok := c["name"].(string)
		if ok {
			if _, exists := expected[name]; exists {
				expected[name] = true
			}
		}
	}

	for name, found := range expected {
		if !found {
			t.Errorf("expected category %q not found in GetCategories()", name)
		}
	}
}

func TestGetCategories_SubCategoriesPresent(t *testing.T) {
	cats := GetCategories()
	for _, c := range cats {
		subs, ok := c["subCategories"].([]string)
		if !ok || len(subs) == 0 {
			t.Errorf("category %v has no subcategories", c["name"])
		}
	}
}
