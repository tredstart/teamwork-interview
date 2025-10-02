package exporter

import (
	"bytes"
	"container/heap"
	"fmt"
	"importer/customerimporter"
	"importer/logger"
	"testing"
)

func TestExportData(t *testing.T) {
	path := "./test_output.csv"
	data := customerimporter.PriorityQueue{
		&customerimporter.DomainData{
			Domain:           "livejournal.com",
			CustomerQuantity: 12,
		},
		&customerimporter.DomainData{
			Domain:           "microsoft.com",
			CustomerQuantity: 22,
		},
		&customerimporter.DomainData{
			Domain:           "newsvine.com",
			CustomerQuantity: 15,
		},
		&customerimporter.DomainData{
			Domain:           "pinteres.uk",
			CustomerQuantity: 10,
		},
		&customerimporter.DomainData{
			Domain:           "yandex.ru",
			CustomerQuantity: 43,
		},
	}
	heap.Init(&data)
	exporter := NewCustomerExporter(&path)

	err := exporter.ExportData(data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExportInvalidPath(t *testing.T) {
	path := ""
	exporter := NewCustomerExporter(&path)

	err := exporter.ExportData(customerimporter.PriorityQueue{})
	if err == nil {
		t.Fatal(err)
	}
	t.Log(err)
}

func TestExportEmptyData(t *testing.T) {
	path := "./test_output.csv"
	exporter := NewCustomerExporter(&path)

	err := exporter.ExportData(customerimporter.PriorityQueue(nil))
	if err == nil {
		t.Fatal(err)
	}
	t.Log(err)
}

func BenchmarkExportDomainData(b *testing.B) {
	b.StopTimer()
	dir := b.TempDir()
	path := fmt.Sprintf("%s/test_output.csv", dir)
	dataPath := "../customerimporter/benchmark10k.csv"
	importer := customerimporter.NewCustomerImporter(&dataPath)
	data, err := importer.ImportDomainData()
	if err != nil {
		b.Error(err)
	}
	exporter := NewCustomerExporter(&path)

	var buf bytes.Buffer
	logger.SetupLogger(&buf, "warn")

	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := exporter.ExportData(data); err != nil {
			b.Fatal(err)
		}
	}
}
