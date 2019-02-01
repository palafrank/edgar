package edgar

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"
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

func (f *filing) filingErrorString() string {
	return "Filing information has not been collected for " + f.FiledOn() + " "
}

func (f *filing) Type() (FilingType, error) {
	if f.FinData != nil {
		return f.FinData.DocType, nil
	}
	return "", errors.New(f.filingErrorString())
}

func (f *filing) ShareCount() (float64, error) {
	if f.FinData != nil && f.FinData.Entity != nil {
		if isCollectedDataSet(f.FinData.Entity, "ShareCount") {
			return f.FinData.Entity.ShareCount, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Share Count")
}

func (f *filing) Revenue() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "Revenue") {
			return f.FinData.Ops.Revenue, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Revenue")
}

func (f *filing) CostOfRevenue() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "CostOfSales") {
			return f.FinData.Ops.CostOfSales, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Cost of Revenue")
}

func (f *filing) GrossMargin() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "GrossMargin") {
			return f.FinData.Ops.GrossMargin, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Gross Margin")
}

func (f *filing) OperatingIncome() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "OpIncome") {
			return f.FinData.Ops.OpIncome, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Operating Income")
}

func (f *filing) OperatingExpense() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "OpExpense") {
			return f.FinData.Ops.OpExpense, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Operating Expense")
}

func (f *filing) NetIncome() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "NetIncome") {
			return f.FinData.Ops.NetIncome, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Net Income")
}

func (f *filing) TotalEquity() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "Equity") {
			return f.FinData.Bs.Equity, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Total Equity")
}

func (f *filing) ShortTermDebt() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "SDebt") {
			return f.FinData.Bs.SDebt, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Short Term Debt")
}

func (f *filing) LongTermDebt() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "LDebt") {
			return f.FinData.Bs.LDebt, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Long Term Debt")
}

func (f *filing) CurrentLiabilities() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "CLiab") {
			return f.FinData.Bs.CLiab, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Current Liabilities")
}

func (f *filing) DeferredRevenue() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "Deferred") {
			return f.FinData.Bs.Deferred, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Deferred Revenue")
}

func (f *filing) RetainedEarnings() (float64, error) {
	if f.FinData != nil && f.FinData.Bs != nil {
		if isCollectedDataSet(f.FinData.Bs, "Retained") {
			return f.FinData.Bs.Retained, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Retained Earnings")
}

func (f *filing) OperatingCashFlow() (float64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		if isCollectedDataSet(f.FinData.Cf, "OpCashFlow") {
			return f.FinData.Cf.OpCashFlow, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Operating Cash Flow")
}

func (f *filing) CapitalExpenditure() (float64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		if isCollectedDataSet(f.FinData.Cf, "CapEx") {
			return f.FinData.Cf.CapEx, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Capital Expenditur")
}

func (f *filing) Dividend() (float64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		if isCollectedDataSet(f.FinData.Cf, "Dividends") {
			// Dividend is recorded as an expense and is -ve. Hence reversing sign
			return f.FinData.Cf.Dividends * -1, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Dividend")
}

func (f *filing) WAShares() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "WAShares") {
			return f.FinData.Ops.WAShares, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Weighted Average Shares")
}

func (f *filing) DividendPerShare() (float64, error) {
	if f.FinData != nil && f.FinData.Ops != nil {
		if isCollectedDataSet(f.FinData.Ops, "Dps") {
			return f.FinData.Ops.Dps, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Dividend Per Share")
}

func (f *filing) Interest() (float64, error) {
	if f.FinData != nil && f.FinData.Cf != nil {
		if isCollectedDataSet(f.FinData.Cf, "Interest") {
			return f.FinData.Cf.Interest, nil
		}
	}
	return 0, errors.New(f.filingErrorString() + "Interest paid")
}

func (f *filing) CollectedData() []string {

	eval := func(data interface{}) []string {
		var ret []string
		if data != nil {
			t := reflect.TypeOf(data)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			for i := 0; i < t.NumField(); i++ {
				if isCollectedDataSet(data, t.Field(i).Name) {
					ret = append(ret, t.Field(i).Name)
				}
			}
		}
		return ret
	}
	ret := eval(f.FinData.Entity)
	ret = append(ret, eval(f.FinData.Bs)...)
	ret = append(ret, eval(f.FinData.Cf)...)
	ret = append(ret, eval(f.FinData.Ops)...)

	return ret
}
