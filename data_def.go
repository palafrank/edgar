package main

import (
	"errors"
	"reflect"
	"strings"
)

type finDataType string

type finDataSearchInfo struct {
	finDataName finDataType
	finDataStr  string
}

var (
	finDataSharesOutstanding finDataType = "Shares Outstanding"
	finDataRevenue           finDataType = "Revenue"
	finDataCostOfRevenue     finDataType = "Cost Of Revenue"
	finDataGrossMargin       finDataType = "Gross Margin"
	finDataOpsIncome         finDataType = "Operational Income"
	finDataOpsExpense        finDataType = "Operational Expense"
	finDataNetIncome         finDataType = "Net Income"
	finDataOpCashFlow        finDataType = "Operating Cash Flow"
	finDataCapEx             finDataType = "Capital Expenditure"
	finDataUnknown           finDataType = "Unknown"

	finDataSearchKeys = []finDataSearchInfo{
		{finDataRevenue, "net revenue"},
		{finDataRevenue, "net sales"},
		{finDataRevenue, "total revenue"},
		{finDataRevenue, "total sales"},
		{finDataCostOfRevenue, "cost of sales"},
		{finDataCostOfRevenue, "cost of revenue"},
		{finDataGrossMargin, "gross margin"},
		{finDataSharesOutstanding, "shares outstanding"},
		{finDataOpsExpense, "operating expenses"},
		{finDataOpsIncome, "operating income"},
		{finDataNetIncome, "net income"},
		{finDataOpCashFlow, "operating activities"},
		{finDataCapEx, "plant and equipment"},
		{finDataCapEx, "capital expen"},
	}
)

func getFinDataType(key string) finDataType {
	key = strings.ToLower(key)
	for _, val := range finDataSearchKeys {
		lup := strings.ToLower(val.finDataStr)
		if strings.Contains(key, lup) {
			return val.finDataName
		}
	}
	return finDataUnknown
}

type EntityData struct {
	ShareCount int64 `finDataType:"Shares Outstanding"`
}

type OpsData struct {
	Revenue     int64 `finDataType:"Revenue"`
	CostOfSales int64 `finDataType:"Cost Of Revenue"`
	GrossMargin int64 `finDataType:"Gross Margin"`
	OpIncome    int64 `finDataType:"Operational Income"`
	OpExpense   int64 `finDataType:"Operational Expense"`
	NetIncome   int64 `finDataType:"Net Income"`
}

type CfData struct {
	OpCashFlow int64 `finDataType:"Operating Cash Flow"`
	CapEx      int64 `finDataType:"Capital Expenditure"`
}

func (e *EntityData) SetData(d string, t finDataType) error {
	switch t {
	case finDataSharesOutstanding:
		e.ShareCount = normalizeNumber(d)
		if e.ShareCount <= 0 {
			return errors.New("Not the share count data")
		}
	}
	return nil
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
		val := v.Field(i).Int()
		if val == 0 {
			err += t.Field(i).Name + ","
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
		tag, ok := t.Field(i).Tag.Lookup("finDataType")
		if ok && string(finType) == tag {
			v.Field(i).SetInt(normalizeNumber(val))
			return nil
		}
	}
	return errors.New("Could not find the field to set: " + string(finType))
}
