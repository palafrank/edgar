package edgar

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
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
	filingDocOps filingDocType = "Operations"
	filingDocInc filingDocType = "Income"
	filingDocBS  filingDocType = "Assets"
	filingDocCF  filingDocType = "Cash Flow"
	filingDocEN  filingDocType = "Entity Info"
	filingDocIg  filingDocType = "Ignore"

	//Scale of the money in the filing
	scaleNone     scaleFactor = 1
	scaleThousand scaleFactor = 1000 * scaleNone
	scaleMillion  scaleFactor = 1000 * scaleThousand
	scaleBillion  scaleFactor = 1000 * scaleMillion

	// Scaling entities in filings
	scaleEntityShares scaleEntity = "Shares"
	scaleEntityMoney  scaleEntity = "Money"

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
	finDataUnknown           finDataType = "Unknown"

	//Required Documents list
	requiredDocTypes = map[filingDocType]bool{
		filingDocOps: true,
		filingDocInc: true,
		filingDocBS:  true,
		filingDocCF:  true,
		filingDocEN:  true,
	}

	filingScales = map[string]scaleInfo{
		"shares in thousand": scaleInfo{scale: scaleThousand, entity: scaleEntityShares},
		"shares in million":  scaleInfo{scale: scaleMillion, entity: scaleEntityShares},
		"$ in million":       scaleInfo{scale: scaleMillion, entity: scaleEntityMoney},
		"$ in billion":       scaleInfo{scale: scaleBillion, entity: scaleEntityMoney},
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

func lookupDocType(data string) filingDocType {

	//fmt.Println("Looking up: ", data)
	data = strings.ToUpper(data)

	if strings.Contains(data, "PARENTHETICAL") {
		//skip this doc
		return filingDocIg
	}

	if strings.Contains(data, "DOCUMENT") && strings.Contains(data, "ENTITY") {
		//Entity document
		return filingDocEN
		//} else if strings.Contains(data, "CONSOLIDATED") {
	} else {
		/*
			match, _ := regexp.MatchString("^(?s)(.*)CONSOLIDATED.*$", data)
			if !match {
				//fmt.Println("PASSING ON :", data)
				return filingDocIg
			}
		*/
		if strings.Contains(data, "BALANCE SHEETS") {
			//Balance sheet
			return filingDocBS
		} else if strings.Contains(data, "OPERATIONS") {
			//Operations statement
			return filingDocOps
		} else if strings.Contains(data, "INCOME") {
			//Income statement
			return filingDocInc
		} else if strings.Contains(data, "CASH FLOWS") {
			//Cash flow statement
			return filingDocCF
		}
	}
	return filingDocIg
}

func getMissingDocs(data map[filingDocType]string) string {
	var ret string
	ret = "[ "
	for key, val := range requiredDocTypes {
		if val == true {
			_, ok := data[key]
			if !ok {
				ret = ret + " " + string(key)
			}
		}
	}
	ret += " ]"
	return ret
}

func generateData(data interface{}, name string) float64 {
	switch name {
	case "GrossMargin":
		val, ok := data.(*opsData)
		if ok {
			//Do this only when the parsing is complete for required fields
			if val.Revenue != 0 && val.CostOfSales != 0 {
				return val.Revenue - val.CostOfSales
			}
		}
	}
	return 0
}

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
		tag, ok := t.Field(i).Tag.Lookup("required")
		val := v.Field(i).Float()
		if val == 0 && (ok && tag == "true") {
			tag, ok = t.Field(i).Tag.Lookup("generate")
			if ok && tag == "true" {
				num := generateData(data, t.Field(i).Name)
				if num == 0 {
					err += t.Field(i).Name + ","
				} else {
					v.Field(i).SetFloat(num)
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
				num := normalizeNumber(val)
				tag, ok := t.Field(i).Tag.Lookup("entity")
				if ok {
					factor, o := scale[scaleEntity(tag)]
					if o {
						num *= float64(factor)
					}
				}
				v.Field(i).SetFloat(num)
			}
			return nil
		}
	}
	return errors.New("Could not find the field to set: " + string(finType))
}
