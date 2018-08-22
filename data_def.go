package edgar_parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
)

type finDataType string
type filingType string
type filingDocType string

type finDataSearchInfo struct {
	finDataName finDataType
	finDataStr  []string
}

var (
	// Threshold year is the earliest year for which we will collect data
	thresholdYear int = 2011
	//Filing types
	filingType10Q filingType = "10-Q"
	filingType10K filingType = "10-K"

	//Document types
	filingDocOps filingDocType = "Operations"
	filingDocInc filingDocType = "Income"
	filingDocBS  filingDocType = "Assets"
	filingDocCF  filingDocType = "Cash Flow"
	filingDocEN  filingDocType = "Entity Info"
	filingDocIg  filingDocType = "Ignore"

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
	finDataUnknown           finDataType = "Unknown"

	//Keys to search for financial data in the filings
	finDataSearchKeys = map[filingDocType][]finDataSearchInfo{
		filingDocOps: {
			{finDataRevenue, []string{"net(?s)(.*)revenue"}},
			{finDataRevenue, []string{"net(?s)(.*)sales"}},
			{finDataRevenue, []string{"total(?s)(.*)revenue"}},
			{finDataRevenue, []string{"total(?s)(.*)sales"}},
			{finDataCostOfRevenue, []string{"cost(?s)(.*)sales"}},
			{finDataCostOfRevenue, []string{"cost(?s)(.*)revenue"}},
			{finDataGrossMargin, []string{"gross(?s)(.*)margin"}},
			{finDataOpsExpense, []string{"operating(?s)(.*)expenses"}},
			{finDataOpsIncome, []string{"operating(?s)(.*)income"}},
			{finDataOpsIncome, []string{"operating(?s)(.*)loss"}},
			{finDataNetIncome, []string{"net(?s)(.*)income"}},
		},
		filingDocCF: {
			{finDataOpCashFlow, []string{"operating(?s)(.*)activities"}},
			{finDataCapEx, []string{"plant(?s)(.*)equipment"}},
			{finDataCapEx, []string{"capital(?s)(.*)expense"}},
		},
		filingDocBS: {
			{finDataSDebt, []string{"current portion(?s)(.*)long-term debt"}},
			{finDataLDebt, []string{"long-term debt"}},
			{finDataCLiab, []string{"total(?s)(.*)current(?s)(.*)liabilities"}},
			{finDataDeferred, []string{"deferred(?s)(.*)revenue"}},
			{finDataRetained, []string{"retained(?s)(.*)earnings"}},
		},
		filingDocEN: {
			{finDataSharesOutstanding, []string{"^(?s)(.*)shares outstanding"}},
		},
	}

	//Required Documents list
	requiredDocTypes = map[filingDocType]bool{
		filingDocOps: true,
		filingDocInc: true,
		filingDocBS:  true,
		filingDocCF:  true,
		filingDocEN:  true,
	}
)

func lookupDocType(data string) filingDocType {

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
		match, _ := regexp.MatchString("^(?s)(.*)CONSOLIDATED.*$", data)
		if !match {
			//fmt.Println("PASSING ON :", data)
			return filingDocIg
		}
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

func getFinDataType(key string, docType filingDocType) finDataType {
	db, ok := finDataSearchKeys[docType]
	if !ok {
		return finDataUnknown
	}
	key = strings.ToLower(key)
	for _, val := range db {
		for _, str := range val.finDataStr {
			match, _ := regexp.MatchString(str, key)
			if match {
				return val.finDataName
			}
		}
	}
	return finDataUnknown
}

type Company struct {
	Ticker  string    `json:"Company"`
	Reports []*Filing `json:"Financial Reports"`
}

type Filing struct {
	Date    string           `json:"Report date"`
	FinData *FinancialReport `json:"Financial Data"`
}

type FinancialReport struct {
	DocType filingType  `json:"Filing Type"`
	Entity  *EntityData `json:"Entity Information"`
	Ops     *OpsData    `json:"Operational Information"`
	Bs      *BSData     `json:"Balance Sheet Information"`
	Cf      *CfData     `json:"Cash Flow Information"`
}

type EntityData struct {
	ShareCount int64 `json:"Shares Outstanding" required:"true"`
}

type OpsData struct {
	Revenue     int64 `json:"Revenue" required:"true"`
	CostOfSales int64 `json:"Cost Of Revenue" required:"true"`
	GrossMargin int64 `json:"Gross Margin" required:"true" generate:"true"`
	OpIncome    int64 `json:"Operational Income" required:"true"`
	OpExpense   int64 `json:"Operational Expense" required:"true"`
	NetIncome   int64 `json:"Net Income" required:"true"`
}

type CfData struct {
	OpCashFlow int64 `json:"Operating Cash Flow" required:"true"`
	CapEx      int64 `json:"Capital Expenditure" required:"true"`
}

type BSData struct {
	LDebt    int64 `json:"Long-Term debt" required:"false"`
	SDebt    int64 `json:"Short-Term debt" required:"false"`
	CLiab    int64 `json:"Current Liabilities" required:"true"`
	Deferred int64 `json:"Deferred revenue" required:"false"`
	Retained int64 `json:"Retained Earnings" required:"true"`
}

func (c Company) String() string {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Company data")
	}
	return string(data)
}

func (f Filing) String() string {
	data, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Filing data")
	}
	fmt.Println("FILING")
	return string(data)
}

func (f FinancialReport) String() string {
	data, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling financial data")
	}
	return string(data)
}

func (bs BSData) String() string {
	data, err := json.MarshalIndent(bs, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling balance sheet data")
	}
	return string(data)
}

func (cf CfData) String() string {
	data, err := json.MarshalIndent(cf, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling cash flow data")
	}
	return string(data)
}

func (ops OpsData) String() string {
	data, err := json.MarshalIndent(ops, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Operational information data")
	}
	return string(data)
}

func generateData(data interface{}, name string) int64 {
	switch name {
	case "GrossMargin":
		val, ok := data.(*OpsData)
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
func Validate(data interface{}) error {
	var err string
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		tag, ok := t.Field(i).Tag.Lookup("required")
		val := v.Field(i).Int()
		if val == 0 && (ok && tag == "true") {
			tag, ok = t.Field(i).Tag.Lookup("generate")
			if ok && tag == "true" {
				num := generateData(data, t.Field(i).Name)
				if num == 0 {
					err += t.Field(i).Name + ","
				} else {
					v.Field(i).SetInt(num)
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

func SetData(data interface{}, finType finDataType, val string) error {

	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		tag, ok := t.Field(i).Tag.Lookup("json")
		if ok && string(finType) == tag {
			if v.Field(i).Int() == 0 {
				v.Field(i).SetInt(normalizeNumber(val))
			}
			return nil
		}
	}
	return errors.New("Could not find the field to set: " + string(finType))
}
