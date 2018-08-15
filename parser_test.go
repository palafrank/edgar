package main

import (
	"os"
	"testing"
)

func TestFilingParser(t *testing.T) {
	f, _ := os.Open("./sample_10Q.html")
	docs := filingPageParser(f, filingType10Q)
	if len(docs) != 5 {
		t.Error("Did not get the expected number of filing document in the 10Q")
	}
}

func TestEntityParser(t *testing.T) {
	//sample := `<tr class="ro"><td class="pl " style="border-bottom: 0px;" valign="top"><a class="a" href="javascript:void(0);" onclick="top.Show.showAR( this, 'defref_dei_DocumentType', window );">Document Type</a></td><td class="text">10-Q<span></span></td><td class="text">&#160;<span></span></td></tr>`
	//f := strings.NewReader(sample)
	f, _ := os.Open("./sample_entity.html")
	entity := getEntityData(f)
	if entity == nil || entity.shareCount != 4829926 {
		t.Error("Unexpected share count in enity parsing")
	}

}
