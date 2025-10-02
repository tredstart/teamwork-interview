package customerimporter

import (
	"bytes"
	"container/heap"
	"importer/logger"
	"strings"
	"testing"
)

func TestImportData(t *testing.T) {
	path := "./test_data.csv"
	importer := NewCustomerImporter(&path)

	_, err := importer.ImportDomainData()
	if err != nil {
		t.Error(err)
	}
}

func TestImportDataSort(t *testing.T) {
	domains := []string{
		"360.cn",
		"acquirethisname.com",
		"blogtalkradio.com",
		"chicagotribune.com",
		"cnet.com",
		"cyberchimps.com",
		"github.io",
		"hubpages.com",
		"rediff.com",
		"statcounter.com",
	}
	path := "./test_data.csv"
	importer := NewCustomerImporter(&path)
	data, err := importer.ImportDomainData()
	if err != nil {
		t.Error(err)
	}
	if len(data) != len(domains) {
		t.Errorf("sorted data has a wrong number of domains \n %v \n\n vs %v\n", domains, data)
	}
	output := make([]uint64, len(data))
	for len(data) > 0 {
		output = append(output, heap.Pop(&data).(*DomainData).CustomerQuantity)
	}
	for i := 0; i < len(data)-1; i++ {
		if data[i].CustomerQuantity > data[i+1].CustomerQuantity {
			t.Errorf("data is not sorted: %v", data)
		}
	}
}

func TestImportInvalidPath(t *testing.T) {
	path := ""
	importer := NewCustomerImporter(&path)

	_, err := importer.ImportDomainData()
	if err == nil {
		t.Error("invalid path error not caught")
	}
}

func TestImportInvalidData(t *testing.T) {
	path := "./test_invalid_data.csv"
	importer := NewCustomerImporter(&path)

	var buf bytes.Buffer

	logger.SetupLogger(&buf, "warn")

	_, err := importer.ImportDomainData()
	if err != nil {
		t.Errorf("unexpected fail: %s", err)
	}

	if !strings.Contains(buf.String(), "invalid email") {
		t.Error("invalid data isn't caught")
	}

}

func BenchmarkImportDomainData(b *testing.B) {
	b.StopTimer()
	path := "./benchmark10k.csv"
	importer := NewCustomerImporter(&path)

	var buf bytes.Buffer
	logger.SetupLogger(&buf, "warn")

	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := importer.ImportDomainData(); err != nil {
			b.Error(err)
		}
	}

}
