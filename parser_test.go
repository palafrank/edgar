package main

import (
	"os"
	"testing"
)

func TestFilingQuery(t *testing.T) {
	f, _ := os.Open("./sample_query.html")
	links := queryPageParser(f, filingType10Q)
	if len(links) != 10 {
		t.Error("Incorrect number of filing links found")
	}
}

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

func TestCfParser(t *testing.T) {
	f, _ := os.Open("./sample_cf.html")
	cf, err := getCfData(f)
	if err != nil {
		t.Error(err.Error())
	} else if cf == nil {
		t.Error("Cash flow object was not returned")
	} else {
		if cf.OpCashFlow != 57911 {
			t.Error("Incorrect cash flow from operations value parsed")
		}
		if cf.CapEx != int64(-10272) {
			t.Error("Incorrect capital expenditure value parsed")
		}
	}
}

func TestBSParser(t *testing.T) {
	f, _ := os.Open("./sample_bs.html")
	bs, err := getBSData(f)
	if err != nil {
		t.Error(err.Error())
	} else if bs == nil {
		t.Error("Balance sheet object was not returned")
	} else {
		if bs.CLiab != 88548 {
			t.Error("Incorrect current liabilities from balance sheet value parsed")
		}
		if bs.LDebt != 97128 {
			t.Error("Incorrect long term debt from balance sheet value parsed")
		}
		if bs.SDebt != 5498 {
			t.Error("Incorrect short term debt from balance sheet value parsed")
		}
		if bs.Retained != 79436 {
			t.Error("Incorrect retained earningd from balance sheet value parsed")
		}
	}
}
