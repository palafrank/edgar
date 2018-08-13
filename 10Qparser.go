package main

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

var docs10Q = map[string]filingDocType{
	"CONDENSED CONSOLIDATED STATEMENTS OF OPERATIONS":           filingDocOps,
	"CONDENSED CONSOLIDATED STATEMENTS OF COMPREHENSIVE INCOME": filingDocInc,
	"CONDENSED CONSOLIDATED BALANCE SHEETS":                     filingDocBS,
	"CONDENSED CONSOLIDATED STATEMENTS OF CASH FLOWS":           filingDocCF,
	"Document and Entity Information":                           filingDocEN,
}

func getDocType(title string) filingDocType {
	strs := strings.Split(title, " (")
	docType, ok := docs10Q[strs[0]]
	if ok && !strings.Contains(title, "Parenthetical") {
		//Found a wanted document
		return docType
	}
	return filingDocIg
}

func map10QReports(page io.Reader, filingLinks []string) map[filingDocType]string {
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
					docType := getDocType(token.String())
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
