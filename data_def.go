package edgar

import (
	"errors"
	"fmt"
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
	filingErrorString string = "Filing information has not been collected"
	// Threshold year is the earliest year for which we will collect data
	thresholdYear int = 2011

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

type date struct {
	day   int
	month int
	year  int
}

func (d date) Day() int {
	return d.day
}

func (d date) Month() int {
	return d.month
}

func (d date) Year() int {
	return d.year
}

func (d date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day)
}

func generateData(fin *financialReport, name string) float64 {
	switch name {
	case "GrossMargin":
		//Do this only when the parsing is complete for required fields
		if isCollectedDataSet(fin.Ops, "Revenue") && isCollectedDataSet(fin.Ops, "CostOfSales") {
			log.Println("Generating Gross Margin")
			return fin.Ops.Revenue - fin.Ops.CostOfSales
		}

	case "Dps":
		if isCollectedDataSet(fin.Cf, "Dividends") && isCollectedDataSet(fin.Ops, "WAShares") {
			log.Println("Generating Dividends per Share")
			return round(fin.Cf.Dividends * -1 / fin.Ops.WAShares)
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
			return errors.New("[" + err + "] " + "attributes did not parse")
		}
		return nil
	}
	if err := validate(fin.Bs); err != nil {
		return err
	}
	if err := validate(fin.Entity); err != nil {
		return err
	}
	if err := validate(fin.Cf); err != nil {
		return err
	}

	if err := validate(fin.Ops); err != nil {
		return err
	}
	return nil
}

/*
//Validate is a function to check that no field is set to 0 after parsing
func validate(data interface{}) error {
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
				fmt.Println("Generate time ", t.Field(i).Name)
				num := generateData(data, t.Field(i).Name)
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
		return errors.New("[" + err + "] " + "attributes did not parse")
	}
	return nil
}
*/

func setData(data interface{},
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
