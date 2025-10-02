// Package customerimporter reads from a CSV file and returns a sorted (data
// structure of your choice) of email domains along with the number of customers
// with e-mail addresses for each domain. This should be able to be ran from the
// CLI and output the sorted domains to the terminal or to a file. Any errors
// should be logged (or handled). Performance matters (this is only ~3k lines,
// but could be 1m lines or run on a small machine).
package customerimporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

type DomainData struct {
	Domain           string
	CustomerQuantity uint64
}

type DomainDataNode struct {
	DomainData
	ParentNode *DomainDataNode
	LeftNode   *DomainDataNode
	RightNode  *DomainDataNode
}

type DomainDataTree struct {
	Root *DomainDataNode
}

func appendNode(node, newNode *DomainDataNode) {
	if newNode.CustomerQuantity <= node.CustomerQuantity {
		if node.LeftNode == nil {
			newNode.ParentNode = node
			node.LeftNode = newNode
			return
		}
		appendNode(node.LeftNode, newNode)
	} else {
		if node.RightNode == nil {
			newNode.ParentNode = node
			node.RightNode = newNode
			return
		}
		appendNode(node.RightNode, newNode)
	}
}

func (dt *DomainDataTree) Append(node *DomainDataNode) {
	if dt.Root == nil {
		dt.Root = node
		return
	}
	appendNode(dt.Root, node)
}

func walk(node *DomainDataNode, data *[]DomainData) {
	if node == nil {
		return
	}

	walk(node.LeftNode, data)
	*data = append(*data, node.DomainData)
	walk(node.RightNode, data)
}

func (dt DomainDataTree) Slice() []DomainData {
	data := []DomainData{}
	walk(dt.Root, &data)
	return data
}

type CustomerImporter struct {
	path *string
}

// NewCustomerImporter returns a new CustomerImporter that reads from file at specified path.
func NewCustomerImporter(filePath *string) *CustomerImporter {
	return &CustomerImporter{
		path: filePath,
	}
}

// ImportDomainData reads and returns sorted customer domain data from CSV file.
func (ci CustomerImporter) ImportDomainData() (*DomainDataTree, error) {
	slog.Info(fmt.Sprintf("starting import of %s", *ci.path))
	file, err := os.Open(*ci.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	data := make(map[string]uint64)

	// skip first line with headers
	line, readErr := csvReader.Read()
	if readErr != nil {
		slog.Warn("cannot read csv file", "importer", fmt.Sprint(line, readErr))
		return nil, readErr
	}
	for line, readErr := csvReader.Read(); readErr != io.EOF; line, readErr = csvReader.Read() {
		if readErr != nil {
			return nil, readErr
		}
		email, domain, found := strings.Cut(line[2], "@")
		if email == "" || !found {
			slog.Warn("skipping line", "importer", fmt.Sprintf("invalid email address: %s", line[2]))
			continue
		}
		data[domain] += 1
	}
	domainData := new(DomainDataTree)
	for k, v := range data {
		dataNode := new(DomainDataNode)
		dataNode.Domain = k
		dataNode.CustomerQuantity = v
		domainData.Append(dataNode)
	}

	slog.Info("Import successful")
	return domainData, nil
}
