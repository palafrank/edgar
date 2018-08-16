package main

import (
	"os"
	"testing"
)

func TestFilingParser(t *testing.T) {
	f, _ := os.Open("./sample_10Q.html")
	docs := filingPageParser(f, filingType10Q)
	if len(docs) != 5 {
		t.Error("Did not get the expected number of filing document in the 10Q")
	}
}

func TestEntityParser(t *testing.T) {
	f, _ := os.Open("./sample_entity.html")
	entity, err := getEntityData(f)
	if err != nil {
		t.Error(err.Error())
	} else if entity == nil {
		t.Error("Entity object was not returned")
	} else if entity.ShareCount != 4829926 {
		t.Error("Incorrect sharecount value parsed")
	}
}

func TestOpsParser(t *testing.T) {
	f, _ := os.Open("./sample_ops.html")
	ops, err := getOpsData(f)
	if err != nil {
		t.Error(err.Error())
	} else if ops == nil {
		t.Error("Operations object was not returned")
	} else {
		if ops.Revenue != 53265 {
			t.Error("Revenue amount did not match")
		}
		if ops.CostOfSales != 32844 {
			t.Error("Cost of Sales amount did not match")
		}
		if ops.GrossMargin != 20421 {
			t.Error("Gross margin amount did not match")
		}
		if ops.OpExpense != 7809 {
			t.Error("Operational Expense amount did not match")
		}
		if ops.OpIncome != 12612 {
			t.Error("Operational Income amount did not match")
		}
		if ops.NetIncome != 11519 {
			t.Error("Net income amount did not match")
		}
	}
}
