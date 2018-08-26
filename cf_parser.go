package edgar

import (
	"io"

	"golang.org/x/net/html"
)

func getCfData(page io.Reader) (*cfData, error) {
	retData := new(cfData)

	z := html.NewTokenizer(page)

	data, err := parseTableRow(z, false)
	for err == nil {
		if len(data) > 0 {
			finType := getFinDataType(data[0], filingDocCF)
			if finType != finDataUnknown {
				for _, str := range data[1:] {
					if len(str) > 0 {
						if setData(retData, finType, str) == nil {
							break
						}
					}
				}
			}
		}
		//Early break out if all required data is collected
		if validate(retData) == nil {
			break
		}
		data, err = parseTableRow(z, false)
	}
	return retData, validate(retData)
}
