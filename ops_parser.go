package edgar

import (
	"io"

	"golang.org/x/net/html"
)

func getOpsData(page io.Reader) (*opsData, error) {
	retData := new(opsData)

	z := html.NewTokenizer(page)

	scales := parseFilingScale(z)
	data, err := parseTableRow(z, true)
	for err == nil {
		if len(data) > 0 {
			finType := getFinDataType(data[0], filingDocOps)
			if finType != finDataUnknown {
				for _, str := range data[1:] {
					if len(str) > 0 {
						if setData(retData, finType, str, scales) == nil {
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
		data, err = parseTableRow(z, true)
	}
	return retData, validate(retData)
}
