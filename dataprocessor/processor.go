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
		slog.Error("error importing customer data: ", err.Error(), err)
		return
	}
	ExportData(data, outFile)
}

func ExportData(data []customerimporter.DomainData, outFile string) {
	if outFile == "" {
		printData(data)
	} else {
		exporter := exporter.NewCustomerExporter(&outFile)
		if saveErr := exporter.ExportData(data); saveErr != nil {
			slog.Error("error saving domain data: ", saveErr.Error(), saveErr)
		}
	}
}

func printData(data []customerimporter.DomainData) {
	fmt.Println("domain,number_of_customers")
	for _, v := range data {
		fmt.Printf("%s,%v\n", v.Domain, v.CustomerQuantity)
	}
}
