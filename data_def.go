package main

import (
	"errors"
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

type entityData struct {
	shareCount uint64
}

type opsData struct {
	revenue     uint64
	costOfSales uint64
	grossMargin uint64
	opIncome    uint64
	opExpense   uint64
	netIncome   uint64
}

func (e *entityData) SetData(d string, t finDataType) error {
	switch t {
	case finDataSharesOutstanding:
		e.shareCount = normalizeNumber(d)
		if e.shareCount <= 0 {
			return errors.New("Not the share count data")
		}
	}
	return nil
}
