package main

import (
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
				finType := getFinDataType(data[0])
				if finType != finDataUnknown {
					for _, str := range data[1:] {
						if len(str) > 0 {
							if retData.SetData(str, finType) == nil {
								break
							}
						}
					}
				}
			}

		}
		tt = z.Next()
	}
	return retData
}
