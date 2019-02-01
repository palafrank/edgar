package edgar

import (
	"encoding/json"
	"log"
)

type financialReport struct {
	DocType FilingType  `json:"Filing Type"`
	Entity  *entityData `json:"Entity Information"`
	Ops     *opsData    `json:"Operational Information"`
	Bs      *bsData     `json:"Balance Sheet Information"`
	Cf      *cfData     `json:"Cash Flow Information"`
}

type entityData struct {
	CollectedData uint64  `json:"Collected Data"`
	ShareCount    float64 `json:"Shares Outstanding" required:"true" entity:"Shares" bit:"0"`
}

type opsData struct {
	CollectedData uint64  `json:"Collected Data"`
	Revenue       float64 `json:"Revenue" required:"true" entity:"Money" bit:"0"`
	CostOfSales   float64 `json:"Cost Of Revenue" required:"true" entity:"Money" bit:"1"`
	GrossMargin   float64 `json:"Gross Margin" required:"true" generate:"true" entity:"Money" bit:"2"`
	OpIncome      float64 `json:"Operational Income" required:"true" entity:"Money" bit:"3"`
	OpExpense     float64 `json:"Operational Expense" required:"true" generate:"true" entity:"Money" bit:"4"`
	NetIncome     float64 `json:"Net Income" required:"true" entity:"Money" bit:"5"`
	WAShares      float64 `json:"Weighted Average Share Count" required:"true" entity:"Shares" bit:"6"`
	Dps           float64 `json:"Dividend Per Share" required:"true" generate:"true" entity:"PerShare" bit:"7"`
}

type cfData struct {
	CollectedData uint64  `json:"Collected Data"`
	OpCashFlow    float64 `json:"Operating Cash Flow" required:"true" entity:"Money" bit:"0"`
	CapEx         float64 `json:"Capital Expenditure" required:"true" entity:"Money" bit:"1"`
	Dividends     float64 `json:"Dividends paid" required:"false" entity:"Money" bit:"2"`
	Interest      float64 `json:"Interest paid" required:"false" entity:"Money" bit:"3"`
}

type bsData struct {
	CollectedData uint64  `json:"Collected Data"`
	LDebt         float64 `json:"Long-Term debt" required:"false" entity:"Money" bit:"0"`
	SDebt         float64 `json:"Short-Term debt" required:"false" entity:"Money" bit:"1"`
	CLiab         float64 `json:"Current Liabilities" required:"true" entity:"Money" bit:"2"`
	Deferred      float64 `json:"Deferred revenue" required:"false" entity:"Money" bit:"3"`
	Retained      float64 `json:"Retained Earnings" required:"true" entity:"Money" bit:"4"`
	Equity        float64 `json:"Total Shareholder Equity" required:"true" entity:"Money" bit:"5"`
}

func newFinancialReport(docType FilingType) *financialReport {
	fr := new(financialReport)
	fr.DocType = docType
	fr.Bs = new(bsData)
	fr.Cf = new(cfData)
	fr.Entity = new(entityData)
	fr.Ops = new(opsData)
	return fr
}

func (f financialReport) String() string {
	data, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling financial data")
	}
	return string(data)
}

func (bs bsData) String() string {
	data, err := json.MarshalIndent(bs, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling balance sheet data")
	}
	return string(data)
}

func (cf cfData) String() string {
	data, err := json.MarshalIndent(cf, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling cash flow data")
	}
	return string(data)
}

func (ops opsData) String() string {
	data, err := json.MarshalIndent(ops, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Operational information data")
	}
	return string(data)
}
