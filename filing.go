package edgar

import (
	"encoding/json"
	"errors"
	"fmt"
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
	fmt.Println("FILING")
	return string(data)
}

func (f *filing) Ticker() string {
	return f.Company
}

func (f *filing) FiledOn() Date {
	return getDate(f.Date)
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

func (f *filing) ShareCount() (int64, error) {
	if f.FinData != nil && f.FinData.Entity != nil {
		return f.FinData.Entity.ShareCount, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) Revenue() (int64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.Revenue, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) CostOfRevenue() (int64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.CostOfSales, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) GrossMargin() (int64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.GrossMargin, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) OperatingIncome() (int64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.OpIncome, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) OperatingExpense() (int64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.OpExpense, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) NetIncome() (int64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		return f.FinData.Ops.NetIncome, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) TotalEquity() (int64, error) {
	return 0, errors.New(filingErrorString)
}

func (f *filing) ShortTermDebt() (int64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		return f.FinData.Bs.SDebt, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) LongTermDebt() (int64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		return f.FinData.Bs.LDebt, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) CurrentLiabilities() (int64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		return f.FinData.Bs.CLiab, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) DeferredRevenue() (int64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		return f.FinData.Bs.Deferred, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) RetainedEarnings() (int64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		return f.FinData.Bs.Retained, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) OperatingCashFlow() (int64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		return f.FinData.Cf.OpCashFlow, nil
	}
	return 0, errors.New(filingErrorString)
}

func (f *filing) CapitalExpenditure() (int64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		return f.FinData.Cf.CapEx, nil
	}
	return 0, errors.New(filingErrorString)
}
