package edgar

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
)

type company struct {
	Company     string                            `json:"Company"`
	CIK         string                            `json:"-"`
	FilingLinks map[FilingType]map[string]string  `json:"-"`
	Reports     map[FilingType]map[string]*filing `json:"Financial Reports"`
}

func (c company) String() string {
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Fatal("Error marshaling Company data")
	}
	return string(data)
}

func newCompany(ticker string) *company {
	return &company{
		Company:     ticker,
		CIK:         getCompanyCIK(ticker),
		FilingLinks: make(map[FilingType]map[string]string),
		Reports:     make(map[FilingType]map[string]*filing),
	}
}

func (c *company) Ticker() string {
	return c.Company
}

func (c *company) AvailableFilings(filingType FilingType) []Date {
	var d []Date
	links := c.FilingLinks[filingType]
	for key, _ := range links {
		d = append(d, getDate(key))
	}
	sort.Slice(d, func(i, j int) bool {
		return d[i].String() > d[j].String()
	})
	return d
}

func (c *company) Filing(fileType FilingType, key Date) (Filing, error) {
	file, ok := c.Reports[fileType][key.String()]
	if !ok {
		link, ok1 := c.FilingLinks[fileType][key.String()]
		if !ok1 {
			fmt.Println(c.FilingLinks[fileType])
			return nil, errors.New("No filing available for given date " + key.String())
		}
		file := new(filing)
		var err error
		file.FinData, err = getFinancialData(link, fileType)
		if file.FinData != nil {
			file.Date = key.String()
			file.Company = c.Ticker()
			c.AddReport(file)
			if err != nil {
				log.Println(file.Company + "-Filed on: " + key.String() + ":" + err.Error())
			}
			return file, nil
		} else {
			return nil, err
		}
	}
	return file, nil
}

func (c *company) AddReport(file *filing) {
	t, err := file.Type()
	if err != nil {
		log.Fatal("Adding invalid report")
		return
	}
	if c.Reports[t] == nil {
		c.Reports[t] = make(map[string]*filing)
	}
	c.Reports[t][file.Date] = file
}

// Save the Company folder into the writer in JSON format
func (c *company) SaveFolder(w io.Writer) error {
	_, err := w.Write([]byte(c.String()))
	if err != nil {
		log.Println("Failed to save data")
		return err
	}
	return nil
}
