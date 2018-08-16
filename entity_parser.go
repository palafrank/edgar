package main

import (
	"io"

	"golang.org/x/net/html"
)

var reqEntityData = map[string]finDataType{
	"Entity Common Stock, Shares Outstanding": finDataSharesOutstanding,
}

func getEntityData(page io.Reader) (*EntityData, error) {

	retData := new(EntityData)
	z := html.NewTokenizer(page)

	data, err := parseTableRow(z)
	for err == nil {
		if len(data) > 0 {
			finType := getFinDataType(data[0])
			if finType != finDataUnknown {
				for _, str := range data[1:] {
					if normalizeNumber(str) > 0 {
						err := SetData(retData, finType, str)
						if err != nil {
							return nil, err
						}
						break
					}
				}
			}
		}
		data, err = parseTableRow(z)
	}
	return retData, Validate(retData)
}
