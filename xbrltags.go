package edgar

import (
	"strings"
)

var (
	// A Map of XBRL tags to financial data type
	// This map contains the corresponding GAAP tag and a version of the tag
	// without the GAAP keyword in case the company has only file non-gaap
	xbrlTags = map[string]finDataType{
		//Balance Sheet info
		"defref_us-gaap_StockholdersEquity":                            finDataTotalEquity,
		"StockholdersEquity":                                           finDataTotalEquity,
		"defref_us-gaap_RetainedEarningsAccumulatedDeficit":            finDataRetained,
		"RetainedEarningsAccumulatedDeficit":                           finDataRetained,
		"defref_us-gaap_LiabilitiesCurrent":                            finDataCLiab,
		"LiabilitiesCurrent":                                           finDataCLiab,
		"defref_us-gaap_AssetsCurrent":                                 finDataCAssets,
		"AssetsCurrent":                                                finDataCAssets,
		"defref_us-gaap_Assets":                                        finDataAssets,
		"Assets":                                                       finDataAssets,
		"defref_us-gaap_Liabilities":                                   finDataLiab,
		"Liabilities":                                                  finDataLiab,
		"defref_us-gaap_CashAndCashEquivalentsAtCarryingValue":         finDataCash,
		"CashAndCashEquivalentsAtCarryingValue":                        finDataCash,
		"defref_us-gaap_Goodwill":                                      finDataGoodwill,
		"Goodwill":                                                     finDataGoodwill,
		"defref_us-gaap_IntangibleAssetsNetExcludingGoodwill":          finDataIntangible,
		"IntangibleAssetsNetExcludingGoodwill":                         finDataIntangible,
		"defref_us-gaap_LongTermDebtNoncurrent":                        finDataLDebt,
		"LongTermDebtNoncurrent":                                       finDataLDebt,
		"defref_us-gaap_LongTermDebtAndCapitalLeaseObligations":        finDataLDebt,
		"LongTermDebtAndCapitalLeaseObligations":                       finDataLDebt,
		"defref_us-gaap_ShortTermBorrowings":                           finDataSDebt,
		"ShortTermBorrowings":                                          finDataSDebt,
		"defref_us-gaap_DebtCurrent":                                   finDataSDebt,
		"DebtCurrent":                                                  finDataSDebt,
		"defref_us-gaap_LongTermDebtAndCapitalLeaseObligationsCurrent": finDataSDebt,
		"LongTermDebtAndCapitalLeaseObligationsCurrent":                finDataSDebt,
		"defref_us-gaap_DeferredRevenueCurrent":                        finDataDeferred,
		"DeferredRevenueCurrent":                                       finDataDeferred,
		"defref_us-gaap_RetainedEarningsAccumulatedDeficitAndAccumulatedOtherComprehensiveIncomeLossNetOfTax": finDataRetained,
		"RetainedEarningsAccumulatedDeficitAndAccumulatedOtherComprehensiveIncomeLossNetOfTax":                finDataRetained,

		//Operations Sheet info
		"defref_us-gaap_SalesRevenueNet": finDataRevenue,
		"SalesRevenueNet":                finDataRevenue,
		"defref_us-gaap_Revenues":        finDataRevenue,
		"Revenues":                       finDataRevenue,
		"defref_us-gaap_RevenueFromContractWithCustomerExcludingAssessedTax":            finDataRevenue,
		"RevenueFromContractWithCustomerExcludingAssessedTax":                           finDataRevenue,
		"defref_us-gaap_CostOfRevenue":                                                  finDataCostOfRevenue,
		"defref_us-gaap_CostOfGoodsAndServicesSold":                                     finDataCostOfRevenue,
		"CostOfGoodsAndServicesSold":                                                    finDataCostOfRevenue,
		"defref_us-gaap_CostOfPurchasedOilAndGas":                                       finDataCostOfRevenue,
		"CostOfPurchasedOilAndGas":                                                      finDataCostOfRevenue,
		"defref_us-gaap_CostOfGoodsSold":                                                finDataCostOfRevenue,
		"CostOfGoodsSold":                                                               finDataCostOfRevenue,
		"defref_us-gaap_CostOfGoodsSoldExcludingAmortizationOfAcquiredIntangibleAssets": finDataCostOfRevenue,
		"CostOfGoodsSoldExcludingAmortizationOfAcquiredIntangibleAssets":                finDataCostOfRevenue,
		"defref_us-gaap_GrossProfit":                                                    finDataGrossMargin,
		"GrossProfit":                                                                   finDataGrossMargin,
		"defref_us-gaap_OperatingExpenses":                                              finDataOpsExpense,
		"OperatingExpenses":                                                             finDataOpsExpense,
		"defref_us-gaap_CostsAndExpenses":                                               finDataOpsExpense,
		"CostsAndExpenses":                                                              finDataOpsExpense,
		"defref_us-gaap_OtherCostAndExpenseOperating":                                   finDataOpsExpense,
		"OtherCostAndExpenseOperating":                                                  finDataOpsExpense,
		"defref_us-gaap_OperatingIncomeLoss":                                            finDataOpsIncome,
		"OperatingIncomeLoss":                                                           finDataOpsIncome,
		"defref_us-gaap_IncomeLossIncludingPortionAttributableToNoncontrollingInterest": finDataOpsIncome,
		"defref_us-gaap_IncomeLossFromContinuingOperationsIncludingPortionAttributableToNoncontrollingInterest":                      finDataOpsIncome,
		"IncomeLossFromContinuingOperationsIncludingPortionAttributableToNoncontrollingInterest":                                     finDataOpsIncome,
		"IncomeLossIncludingPortionAttributableToNoncontrollingInterest":                                                             finDataOpsIncome,
		"defref_us-gaap_IncomeLossFromContinuingOperationsBeforeIncomeTaxesMinorityInterestAndIncomeLossFromEquityMethodInvestments": finDataOpsIncome,
		"IncomeLossFromContinuingOperationsBeforeIncomeTaxesMinorityInterestAndIncomeLossFromEquityMethodInvestments":                finDataOpsIncome,
		"defref_us-gaap_NetIncomeLoss": finDataNetIncome,
		"NetIncomeLoss":                finDataNetIncome,
		"defref_us-gaap_ProfitLoss":    finDataNetIncome,
		"ProfitLoss":                   finDataNetIncome,
		"defref_us-gaap_NetIncomeLossAvailableToCommonStockholdersBasic": finDataNetIncome,
		"NetIncomeLossAvailableToCommonStockholdersBasic":                finDataNetIncome,
		"defref_us-gaap_WeightedAverageNumberOfDilutedSharesOutstanding": finDataWAShares,
		"WeightedAverageNumberOfDilutedSharesOutstanding":                finDataWAShares,
		"defref_us-gaap_CommonStockDividendsPerShareDeclared":            finDataDps,
		"CommonStockDividendsPerShareDeclared":                           finDataDps,

		"defref_us-gaap_IncomeLossFromContinuingOperationsBeforeIncomeTaxesExtraordinaryItemsNoncontrollingInterest": finDataOpsIncome,
		"IncomeLossFromContinuingOperationsBeforeIncomeTaxesExtraordinaryItemsNoncontrollingInterest":                finDataOpsIncome,

		//Cash Flow Sheet info
		"defref_us-gaap_NetCashProvidedByUsedInOperatingActivities":                     finDataOpCashFlow,
		"NetCashProvidedByUsedInOperatingActivities":                                    finDataOpCashFlow,
		"defref_us-gaap_NetCashProvidedByUsedInOperatingActivitiesContinuingOperations": finDataOpCashFlow,
		"NetCashProvidedByUsedInOperatingActivitiesContinuingOperations":                finDataOpCashFlow,
		"defref_us-gaap_PaymentsToAcquirePropertyPlantAndEquipment":                     finDataCapEx,
		"PaymentsToAcquirePropertyPlantAndEquipment":                                    finDataCapEx,
		"defref_us-gaap_PaymentsToAcquireProductiveAssets":                              finDataCapEx,
		"PaymentsToAcquireProductiveAssets":                                             finDataCapEx,
		"defref_us-gaap_CapitalExpendituresAndInvestments":                              finDataCapEx,
		"CapitalExpendituresAndInvestments":                                             finDataCapEx,
		"defref_us-gaap_PaymentsOfDividends":                                            finDataDividend,
		"PaymentsOfDividends":                                                           finDataDividend,
		"defref_us-gaap_PaymentsOfDividendsCommonStock":                                 finDataDividend,
		"PaymentsOfDividendsCommonStock":                                                finDataDividend,
		"defref_us-gaap_InterestPaidNet":                                                finDataInterest,
		"InterestPaidNet":                                                               finDataInterest,
		"defref_us-gaap_InterestAndDebtExpense":                                         finDataInterest,
		"InterestAndDebtExpense":                                                        finDataInterest,
		"defref_us-gaap_InterestIncomeExpenseNet":                                       finDataInterest,
		"InterestIncomeExpenseNet":                                                      finDataInterest,
		//Entity sheet information
		"defref_dei_EntityCommonStockSharesOutstanding": finDataSharesOutstanding,
		"EntityCommonStockSharesOutstanding":            finDataSharesOutstanding,
	}
)

func getFinDataTypeFromXBRLTag(key string) finDataType {

	data, ok := xbrlTags[key]
	if !ok {

		// Now look for non-gaap filing
		// defref_us-gaap_XXX could be filed company specific
		// as defref_msft_XXX
		splits := strings.Split(key, "_")
		if len(splits) == 3 {
			data, ok = xbrlTags[splits[2]]
			if ok {
				return data
			}
		}
		return finDataUnknown
	}
	return data
}
