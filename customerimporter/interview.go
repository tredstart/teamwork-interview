package customerimporter

import (
	"container/heap"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

type PriorityQueue []*DomainData

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].CustomerQuantity <= pq[j].CustomerQuantity
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*DomainData)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *DomainData, domain string, quantity uint64) {
	item.Domain = domain
	item.CustomerQuantity = quantity
	heap.Fix(pq, item.Index)
}

type DomainData struct {
	Domain           string
	CustomerQuantity uint64
	Index            int
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
func (ci CustomerImporter) ImportDomainData() (PriorityQueue, error) {
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
	domainData := make(PriorityQueue, len(data))
	i := 0
	for k, v := range data {
		item := &DomainData{
			Domain:           k,
			CustomerQuantity: v,
			Index:            i,
		}
		domainData[i] = item
		i++
	}
	heap.Init(&domainData)

	slog.Info("Import successful")
	return domainData, nil
}
