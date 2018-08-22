package edgar_parser

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

	data, err := parseTableRow(z, false)
	for err == nil {
		if len(data) > 0 {
			finType := getFinDataType(data[0], filingDocEN)
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
		//Early break out if all required data is collected
		if Validate(retData) == nil {
			break
		}
		data, err = parseTableRow(z, false)
	}
	return retData, Validate(retData)
}
