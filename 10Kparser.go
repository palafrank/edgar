package main

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var docs10K = map[string]filingDocType{
	"CONSOLIDATED STATEMENTS OF OPERATIONS":           filingDocOps,
	"CONSOLIDATED STATEMENTS OF COMPREHENSIVE INCOME": filingDocInc,
	"CONSOLIDATED BALANCE SHEETS":                     filingDocBS,
	"CONSOLIDATED STATEMENTS OF CASH FLOWS":           filingDocCF,
	"DOCUMENT AND ENTITY INFORMATION":                 filingDocEN,
}

func map10KReports(page io.Reader, filingLinks []string) map[filingDocType]string {
	retData := make(map[filingDocType]string)

	z := html.NewTokenizer(page)
	tt := z.Next()
	for tt != html.ErrorToken {
		token := z.Token()
		if token.Data == "var" {
			fmt.Println("Found the var")
		}
		if token.Data == "a" {
			for _, a := range token.Attr {
				if a.Key == "href" && strings.Contains(a.Val, "loadReport") {
					strs := strings.Split(a.Val, "loadReport")
					strs[1] = strings.Trim(strs[1], ";")
					reportNum, _ := strconv.Atoi(strings.Trim(strs[1], "()"))
					tt = z.Next() //Contains the text that describes the report
					if tt != html.TextToken {
						break
					}
					token = z.Token()
					docType := getDocType(token.String(), filingType10K)
					if docType != filingDocIg {
						//Get the report number
						//fmt.Println("Found a wanted doc ", docType, token.String(), reportNum)
						retData[docType] = filingLinks[reportNum-1]
					}
				}
			}
		}
		tt = z.Next()
	}
	if len(retData) <= 0 {
		log.Fatal("Did not find any documents for the filing requested")
	}
	return retData
}
