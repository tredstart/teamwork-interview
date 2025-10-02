package dataprocessor

import (
	"fmt"
	"importer/customerimporter"
	"importer/exporter"
	"log/slog"
)

func ProcessDomainData(path, outFile string) {
	importer := customerimporter.NewCustomerImporter(&path)
	data, err := importer.ImportDomainData()
	if err != nil {
		slog.Error("error importing customer data", "importer", err)
		return
	}
	ExportData(data.Slice(), outFile)
}

func ExportData(data []customerimporter.DomainData, outFile string) {
	slog.Debug("Trying to export data...")
	if outFile == "" {
		printData(data)
	} else {
		exporter := exporter.NewCustomerExporter(&outFile)
		if saveErr := exporter.ExportData(data); saveErr != nil {
			slog.Error("error saving domain data", "exporter", saveErr)
		}
	}
}

func printData(data []customerimporter.DomainData) {
	slog.Info("Printing data to the stdout.")
	fmt.Println("domain,number_of_customers")
	for _, v := range data {
		fmt.Printf("%s,%v\n", v.Domain, v.CustomerQuantity)
	}
}
