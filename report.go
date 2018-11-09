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
	ShareCount float64 `json:"Shares Outstanding" required:"true" entity:"Shares"`
}

type opsData struct {
	Revenue     float64 `json:"Revenue" required:"true" entity:"Money"`
	CostOfSales float64 `json:"Cost Of Revenue" required:"true" entity:"Money"`
	GrossMargin float64 `json:"Gross Margin" required:"true" generate:"true" entity:"Money"`
	OpIncome    float64 `json:"Operational Income" required:"true" entity:"Money"`
	OpExpense   float64 `json:"Operational Expense" required:"true" entity:"Money"`
	NetIncome   float64 `json:"Net Income" required:"true" entity:"Money"`
}

type cfData struct {
	OpCashFlow float64 `json:"Operating Cash Flow" required:"true" entity:"Money"`
	CapEx      float64 `json:"Capital Expenditure" required:"true" entity:"Money"`
}

type bsData struct {
	LDebt    float64 `json:"Long-Term debt" required:"false" entity:"Money"`
	SDebt    float64 `json:"Short-Term debt" required:"false" entity:"Money"`
	CLiab    float64 `json:"Current Liabilities" required:"true" entity:"Money"`
	Deferred float64 `json:"Deferred revenue" required:"false" entity:"Money"`
	Retained float64 `json:"Retained Earnings" required:"true" entity:"Money"`
	Equity   float64 `json:"Total Shareholder Equity" required:"true" entity:"Money"`
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
