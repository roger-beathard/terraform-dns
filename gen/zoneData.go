/*
Copyright © 2020 Roger V. Beathard

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gen

import (
	"sort"
	"strings"
)

// DNZZone is the common zone definition
var DNSZone DNSZoneDef

// DNSRecord represents a single DNS Record
type DNSRecord struct {
	// read from csv file
	HostName    string
	RRDatas     string
	TTL         string
	Type        string
	Description string

	// computed
	ManagedZone  string // obtained from zone information
	ProjectID    string // obtained from zone information
	Name         string // FQDNS Name, calcuated HostName + zone DNSName
	ResourceName string // calculated by modified FQDNSName. dash replace dots
}

// DNSZoneDef maintains the attributes of a zone
type DNSZoneDef struct {
	// read from csv file
	DNSName     string
	ProjectID   string
	Description string

	// computed
	ResourceName string // (Also used in name attribute) calculated by modified DNSName. dash replace dots

	DNSRecords []DNSRecord
}

// SortRecords - zone records by Type and HostName
func (dnsZone *DNSZoneDef) SortRecords() {
	sort.SliceStable(dnsZone.DNSRecords, func(i, j int) bool {
		if dnsZone.DNSRecords[i].Type != dnsZone.DNSRecords[j].Type {
			return dnsZone.DNSRecords[i].Type < dnsZone.DNSRecords[j].Type
		}
		return dnsZone.DNSRecords[i].HostName < dnsZone.DNSRecords[j].HostName
	})
}

// ComputeValues - Populate computed values
func (dnsZone *DNSZoneDef) ComputeValues() {

	// Compute additional zone information
	// May or maynot have a trailing '.' but triming anyway.
	DNSZone.DNSName = strings.TrimSuffix(DNSZone.DNSName, ".")
	DNSZone.ResourceName = strings.ReplaceAll(DNSZone.DNSName, ".", "-")

	// Compute additional zone record values.
	computedRecords := []DNSRecord{}
	for _, dnsRecord := range dnsZone.DNSRecords {

		// Should not have a trailing '.' but trim anyway.
		dnsRecord.HostName = strings.TrimPrefix(dnsRecord.HostName, ".")

		// Escape rrdatas if contain spaces
		if strings.Contains(dnsRecord.RRDatas, " ") {
			dnsRecord.RRDatas = "\\\"" + dnsRecord.RRDatas + "\\\""
		}
		dnsRecord.ProjectID = DNSZone.ProjectID
		dnsRecord.ManagedZone = DNSZone.ResourceName
		dnsRecord.Name = dnsRecord.HostName + "." + DNSZone.DNSName
		dnsRecord.ResourceName = strings.ReplaceAll(dnsRecord.Name, ".", "-")
		computedRecords = append(computedRecords, dnsRecord)
	}
	dnsZone.DNSRecords = computedRecords
}
