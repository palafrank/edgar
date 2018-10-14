package edgar

import (
	"io"

	"golang.org/x/net/html"
)

var reqEntityData = map[string]finDataType{
	"Entity Common Stock, Shares Outstanding": finDataSharesOutstanding,
}

func getEntityData(page io.Reader) (*entityData, error) {

	retData := new(entityData)
	z := html.NewTokenizer(page)

	scales := parseFilingScale(z)
	data, err := parseTableRow(z, true)
	for err == nil {
		if len(data) > 0 {
			finType := getFinDataType(data[0], filingDocEN)
			if finType != finDataUnknown {
				for _, str := range data[1:] {
					if normalizeNumber(str) > 0 {
						err := setData(retData, finType, str, scales)
						if err != nil {
							return nil, err
						}
						break
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
