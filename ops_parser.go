package main

import (
	"io"
)

func getOpsData(page io.Reader) *opsData {
	retData := new(opsData)
	/*
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
	*/
	return retData
}
