package edgar_parser

import (
	"io"

	"golang.org/x/net/html"
)

func getBSData(page io.Reader) (*BSData, error) {
	retData := new(BSData)

	z := html.NewTokenizer(page)

	data, err := parseTableRow(z, false)
	for err == nil {
		if len(data) > 0 {
			finType := getFinDataType(data[0], filingDocBS)
			if finType != finDataUnknown {
				for _, str := range data[1:] {
					if len(str) > 0 {
						if SetData(retData, finType, str) == nil {
							break
						}
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
