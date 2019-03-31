package edgar

import (
	"errors"
	"log"
	"reflect"
)

type finDataType string
type filingDocType string
type scaleFactor int64
type scaleEntity string

type finDataSearchInfo struct {
	finDataName finDataType
	finDataStr  []string
}

type scaleInfo struct {
	scale  scaleFactor
	entity scaleEntity
}

var (

	// Threshold year is the earliest year for which we will collect data
	thresholdYear int = 2012

	//Document types
	filingDocOps      filingDocType = "Operations"
	filingDocInc      filingDocType = "Income"
	filingDocBS       filingDocType = "Assets"
	filingDocCF       filingDocType = "Cash Flow"
	filingDocEN       filingDocType = "Entity Info"
	filingDocEPSNotes filingDocType = "Notes on EPS"
	filingDocEquity   filingDocType = "Notes on Equity"
	filingDocDebt     filingDocType = "Notes on Debt"
	filingDocIg       filingDocType = "Ignore"

	//Scale of the money in the filing
	scaleNone     scaleFactor = 1
	scaleThousand scaleFactor = 1000 * scaleNone
	scaleMillion  scaleFactor = 1000 * scaleThousand
	scaleBillion  scaleFactor = 1000 * scaleMillion

	// Scaling entities in filings
	scaleEntityShares   scaleEntity = "Shares"
	scaleEntityMoney    scaleEntity = "Money"
	scaleEntityPerShare scaleEntity = "PerShare"

	//Types of financial data collected
	finDataSharesOutstanding finDataType = "Shares Outstanding"
	finDataRevenue           finDataType = "Revenue"
	finDataCostOfRevenue     finDataType = "Cost Of Revenue"
	finDataGrossMargin       finDataType = "Gross Margin"
	finDataOpsIncome         finDataType = "Operational Income"
	finDataOpsExpense        finDataType = "Operational Expense"
	finDataNetIncome         finDataType = "Net Income"
	finDataOpCashFlow        finDataType = "Operating Cash Flow"
	finDataCapEx             finDataType = "Capital Expenditure"
	finDataLDebt             finDataType = "Long-Term debt"
	finDataSDebt             finDataType = "Short-Term debt"
	finDataCLiab             finDataType = "Current Liabilities"
	finDataCAssets           finDataType = "Current Assets"
	finDataCash              finDataType = "Cash"
	finDataGoodwill          finDataType = "Goodwill"
	finDataIntangible        finDataType = "Intangibles"
	finDataSecurities        finDataType = "Securities"
	finDataDeferred          finDataType = "Deferred revenue"
	finDataRetained          finDataType = "Retained Earnings"
	finDataTotalEquity       finDataType = "Total Shareholder Equity"
	finDataDividend          finDataType = "Dividends paid"
	finDataWAShares          finDataType = "Weighted Average Share Count"
	finDataDps               finDataType = "Dividend Per Share"
	finDataInterest          finDataType = "Interest paid"
	finDataUnknown           finDataType = "Unknown"

	//Required Documents list
	requiredDocTypes = map[filingDocType]bool{
		filingDocOps: true,
		filingDocInc: true,
		filingDocBS:  true,
		filingDocCF:  true,
		filingDocEN:  true,
	}
)

func generateData(fin *financialReport, name string) float64 {
	log.Println("Generating data: ", name)
	switch name {
	case "GrossMargin":
		//Do this only when the parsing is complete for required fields
		if isCollectedDataSet(fin.Ops, "Revenue") && isCollectedDataSet(fin.Ops, "CostOfSales") {
			log.Println("Generating Gross Margin")
			return fin.Ops.Revenue - fin.Ops.CostOfSales
		}

	case "Dps":
		if isCollectedDataSet(fin.Cf, "Dividends") {
			if isCollectedDataSet(fin.Ops, "WAShares") {
				return round(fin.Cf.Dividends * -1 / fin.Ops.WAShares)
			} else if isCollectedDataSet(fin.Entity, "ShareCount") {
				return round(fin.Cf.Dividends * -1 / fin.Entity.ShareCount)
			}
		}
	case "OpExpense":
		if isCollectedDataSet(fin.Ops, "Revenue") &&
			isCollectedDataSet(fin.Ops, "CostOfSales") &&
			isCollectedDataSet(fin.Ops, "OpIncome") {
			return round(fin.Ops.Revenue - fin.Ops.CostOfSales - fin.Ops.OpIncome)
		}
	}
	return 0
}

func validateFinancialReport(fin *financialReport) error {

	validate := func(data interface{}) error {
		var err string
		t := reflect.TypeOf(data)
		v := reflect.ValueOf(data)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
			v = v.Elem()
		}
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).Type.Kind() != reflect.Float64 {
				continue
			}
			tag, ok := t.Field(i).Tag.Lookup("required")
			if !isCollectedDataSet(data, t.Field(i).Name) && (ok && tag == "true") {
				tag, ok = t.Field(i).Tag.Lookup("generate")
				if ok && tag == "true" {
					num := generateData(fin, t.Field(i).Name)
					if num == 0 {
						err += t.Field(i).Name + ","
					} else {
						v.Field(i).SetFloat(num)
						setCollectedData(data, i)
					}
				} else {
					err += t.Field(i).Name + ","
				}
			}
		}
		if len(err) > 0 {
			return errors.New("[" + err + "]")
		}
		return nil
	}

	// Now make sure the scale factors make sense
	if !isSameScale(fin.Entity.ShareCount, fin.Ops.WAShares) {
		//somethings wrong. Override with Share count
		fin.Ops.WAShares = fin.Entity.ShareCount
	}

	var ret string
	if err := validate(fin.Bs); err != nil {
		ret = ret + "Missing fields in " + string(filingDocBS) + err.Error() + "\n"
	}
	if err := validate(fin.Entity); err != nil {
		ret = ret + "Missing fields in " + string(filingDocEN) + err.Error() + "\n"
	}
	if err := validate(fin.Cf); err != nil {
		ret = ret + "Missing fields in " + string(filingDocCF) + err.Error() + "\n"
	}
	if err := validate(fin.Ops); err != nil {
		ret = ret + "Missing fields in " + string(filingDocOps) + err.Error() + "\n"
	}

	if len(ret) > 0 {
		return errors.New(ret)
	}
	return nil
}

func setData(fr *financialReport,
	finType finDataType,
	val string,
	scale map[scaleEntity]scaleFactor) error {

	setter := func(data interface{},
		finType finDataType,
		val string,
		scale map[scaleEntity]scaleFactor) error {

		t := reflect.TypeOf(data)
		v := reflect.ValueOf(data)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
			v = v.Elem()
		}
		for i := 0; i < t.NumField(); i++ {
			tag, ok := t.Field(i).Tag.Lookup("json")
			if ok && string(finType) == tag {
				if v.Field(i).Float() == 0 {
					num, err := normalizeNumber(val)
					if err != nil {
						return err
					}
					tag, ok := t.Field(i).Tag.Lookup("entity")
					if ok {
						factor, o := scale[scaleEntity(tag)]
						if o {
							num *= float64(factor)
						}
					}
					v.Field(i).SetFloat(num)
					setCollectedData(data, i)
				}
				return nil
			}
		}
		return errors.New("Could not find the field to set: " + string(finType))
	}
	var err error
	if err = setter(fr.Entity, finType, val, scale); err == nil {
		return nil
	}
	if err = setter(fr.Bs, finType, val, scale); err == nil {
		return nil
	}
	if err = setter(fr.Cf, finType, val, scale); err == nil {
		return nil
	}
	if err = setter(fr.Ops, finType, val, scale); err == nil {
		return nil
	}
	return err
}
