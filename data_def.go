package edgar

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

type finDataType string
type filingDocType string

type finDataSearchInfo struct {
	finDataName finDataType
	finDataStr  []string
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

type company struct {
	Company     string                            `json:"Company"`
	FilingLinks map[FilingType]map[string]string  `json:"-"`
	Reports     map[FilingType]map[string]*filing `json:"Financial Reports"`
}

func newCompany(ticker string) *company {
	return &company{
		Company:     ticker,
		FilingLinks: make(map[FilingType]map[string]string),
		Reports:     make(map[FilingType]map[string]*filing),
	}
}

func (c *company) Ticker() string {
	return c.Company
}

func (c *company) AvailableFilings(filingType FilingType) []Date {
	var d []Date
	links := c.FilingLinks[filingType]
	for key, _ := range links {
		d = append(d, getDate(key))
	}
	sort.Slice(d, func(i, j int) bool {
		return d[i].String() > d[j].String()
	})
	return d
}

func (c *company) Filing(fileType FilingType, key Date) (Filing, error) {
	file, ok := c.Reports[fileType][key.String()]
	if !ok {
		link, ok1 := c.FilingLinks[fileType][key.String()]
		if !ok1 {
			fmt.Println(c.FilingLinks[fileType])
			return nil, errors.New("No filing available for given date " + key.String())
		}
		file := new(filing)
		file.FinData = getFinancialData(link, fileType)
		file.Date = key.String()
		file.Company = c.Ticker()
		c.AddReport(file)
		return file, nil
	}
	return file, nil
}

func (c *company) AddReport(file *filing) {
	t, err := file.Type()
	if err != nil {
		log.Fatal("Adding invalid report")
		return
	}
	if c.Reports[t] == nil {
		c.Reports[t] = make(map[string]*filing)
	}
	c.Reports[t][file.Date] = file
}

type filing struct {
	Company string           `json:"Company"`
	Date    string           `json:"Report date"`
	FinData *financialReport `json:"Financial Data"`
}

func (f *filing) Ticker() string {
	return f.Company
}

func (f *filing) Year() int {
	return getYear(f.Date)
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

func (c company) String() string {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Company data")
	}
	return string(data)
}

func (f filing) String() string {
	data, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Filing data")
	}
	fmt.Println("FILING")
	return string(data)
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

func generateData(data interface{}, name string) int64 {
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

func setData(data interface{}, finType finDataType, val string) error {

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
