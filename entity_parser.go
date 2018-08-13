package main

import (
	"fmt"
	"io"

	"golang.org/x/net/html"
)

var reqEntityData = map[string]finDataType{
	"Entity Common Stock, Shares Outstanding": finDataSharesOutstanding,
}

func getEntityData(page io.Reader) *entityData {

	retData := new(entityData)
	z := html.NewTokenizer(page)
	tt := z.Next()
	for tt != html.ErrorToken {
		if tt == html.StartTagToken {
			token := z.Token()
			if token.Data != "tr" {
				tt = z.Next()
				continue
			}

			data := parseTableRow(z)
			if len(data) > 0 {
				finType, ok := reqEntityData[data[0]]
				if ok {
					fmt.Println("Found the share count row", finType, data)
					for _, str := range data {
						if len(str) > 0 {
							retData.SetData(str, finType)
							break
						}
					}
				}
			}

		}
		tt = z.Next()
	}
	return retData
}
