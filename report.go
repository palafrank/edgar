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
	ShareCount int64 `json:"Shares Outstanding" required:"true"`
}

type opsData struct {
	Revenue     int64 `json:"Revenue" required:"true"`
	CostOfSales int64 `json:"Cost Of Revenue" required:"true"`
	GrossMargin int64 `json:"Gross Margin" required:"true" generate:"true"`
	OpIncome    int64 `json:"Operational Income" required:"true"`
	OpExpense   int64 `json:"Operational Expense" required:"true"`
	NetIncome   int64 `json:"Net Income" required:"true"`
}

type cfData struct {
	OpCashFlow int64 `json:"Operating Cash Flow" required:"true"`
	CapEx      int64 `json:"Capital Expenditure" required:"true"`
}

type bsData struct {
	LDebt    int64 `json:"Long-Term debt" required:"false"`
	SDebt    int64 `json:"Short-Term debt" required:"false"`
	CLiab    int64 `json:"Current Liabilities" required:"true"`
	Deferred int64 `json:"Deferred revenue" required:"false"`
	Retained int64 `json:"Retained Earnings" required:"true"`
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
