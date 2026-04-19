package models

import (
	"encoding/json"
	"testing"
)

// TestTransaction_JSONFieldNames guards against accidental json tag changes
// that would break the frontend contract.
func TestTransaction_JSONFieldNames(t *testing.T) {
	tx := Transaction{
		Date: "22/07/2024", Description: "TEST", Amount: 100.5,
		Category: "Shopping", SubCategory: "Online", Hash: "abc",
	}

	data, err := json.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	expected := []string{"date", "description", "amount", "category", "subCategory", "hash"}
	for _, key := range expected {
		if _, ok := m[key]; !ok {
			t.Errorf("JSON missing expected field %q", key)
		}
	}
}

func TestParseResponse_JSONFieldNames(t *testing.T) {
	pr := ParseResponse{
		Transactions: []Transaction{},
		BillPeriod:   BillPeriod{From: "a", To: "b"},
		CardName:     "test",
		TotalAmount:  99.9,
		FileName:     "file.pdf",
	}

	data, err := json.Marshal(pr)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	expected := []string{"transactions", "billPeriod", "cardName", "totalAmount", "fileName"}
	for _, key := range expected {
		if _, ok := m[key]; !ok {
			t.Errorf("JSON missing expected field %q", key)
		}
	}
}

func TestBillPeriod_JSONFieldNames(t *testing.T) {
	bp := BillPeriod{From: "01/01/2024", To: "31/01/2024"}

	data, err := json.Marshal(bp)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}

	if _, ok := m["from"]; !ok {
		t.Error("JSON missing field 'from'")
	}
	if _, ok := m["to"]; !ok {
		t.Error("JSON missing field 'to'")
	}
}
