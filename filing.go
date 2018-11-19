package edgar

import (
	"encoding/json"
	"errors"
	"log"
)

type filing struct {
	Company string           `json:"Company"`
	Date    string           `json:"Report date"`
	FinData *financialReport `json:"Financial Data"`
}

func (f filing) String() string {
	data, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Filing data")
	}
	return string(data)
}

func (f *filing) Ticker() string {
	return f.Company
}

func (f *filing) FiledOn() string {
	return getDate(f.Date).String()
}

func (f *filing) Month() int {
	return getMonth(f.Date)
}

func (f *filing) Type() (FilingType, error) {
	if f.FinData != nil {
		return f.FinData.DocType, nil
	}
	return "", errors.New(filingErrorString)
}

func (f *filing) ShareCount() (float64, error) {
	if f.FinData != nil && f.FinData.Entity != nil {
		return f.FinData.Entity.ShareCount, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) Revenue() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.Revenue, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) CostOfRevenue() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.CostOfSales, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) GrossMargin() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.GrossMargin, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) OperatingIncome() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.OpIncome, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) OperatingExpense() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.OpExpense, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) NetIncome() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.NetIncome, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) TotalEquity() (float64, error) {
	return f.FinData.Bs.Equity, nil
}

func (f *filing) ShortTermDebt() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		return f.FinData.Bs.SDebt, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) LongTermDebt() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		return f.FinData.Bs.LDebt, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) CurrentLiabilities() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		return f.FinData.Bs.CLiab, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) DeferredRevenue() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		return f.FinData.Bs.Deferred, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) RetainedEarnings() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		return f.FinData.Bs.Retained, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) OperatingCashFlow() (float64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		return f.FinData.Cf.OpCashFlow, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) CapitalExpenditure() (float64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		return f.FinData.Cf.CapEx, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) Dividend() (float64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		return f.FinData.Cf.Dividends, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) WAShares() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.WAShares, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) DividendPerShare() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.Dps, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) Interest() (float64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		return f.FinData.Cf.Interest, nil
	}
	return 0, errors.New(filingErrorString)
}
