/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// tfCreateCmd represents the tfCreate command
var tfCreateCmd = &cobra.Command{
	Use:   "tf-create",
	Short: "Create DNS terraform files.",
	Long:  `Used to create DNS terraform files from a csv encoded zone and records file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generate()
	},
}

func init() {
	rootCmd.AddCommand(tfCreateCmd)

	tfCreateCmd.PersistentFlags().StringVar(&zoneInfoFile, "zone-info", "zone.csv", "zone info file (zone.csv)")
	tfCreateCmd.PersistentFlags().StringVar(&recordsFile, "records", "records.csv", "zone records file (records.csv)")
	tfCreateCmd.PersistentFlags().StringVar(&tfFileName, "tfFile", "main.tf", "generated terraform file. (main.tf)")
	tfCreateCmd.PersistentFlags().StringVar(&mdFileName, "mdFile", "zone_info.md", "generated zone info markdown file. (zone_info.md)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tfCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tfCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

//
var zoneInfoFile string
var recordsFile string
var mdFileName string
var tfFileName string

// Template definitions
const tfHeader = `
provider "google" {
  version = "~> 3.9.0"
}
`
const zoneTemplateDef = `
# Zone Info for {{.DNSName}}
resource "google_dns_managed_zone" "{{.ResourceName}}" {
  name          = "{{.ResourceName}}"
  dns_name      = "{{.DNSName}}."
  description   = "{{.Description}}"
  project       = "{{.ProjectID}}"
}
`
const zoneRecordTemplateDef = `
# Zone Record for {{.Name}}
resource "google_dns_record_set" "{{.ResourceName}}" {
  managed_zone  = "{{.ManagedZone}}"
  project       = "{{.ProjectID}}"

  name          = "{{.Name}}."
  type          = "{{.Type}}"
  rrdatas       = ["{{.RRDatas}}"]
  ttl           = "{{.TTL}}"
}

`
const mdZoneHeaderTemplateDef = `
# {{.DNSName}}
ProjectID: {{.ProjectID}}  
Description: {{.Description}}  

---
## Zone Records

host|type|ttl|rrdatas|description
----|----|---|-------|-----------`

const mdZoneRecordTemplateDef = `
{{.HostName}}|{{.Type}}|{{.TTL}}|{{.RRDatas}}|{{.Description}}`

type ZoneInfo struct {
	// read from csv file
	DNSName      string
	Descriptions string
	ProjectID    string
	Description  string

	// computed
	ResourceName string // (Also used in name attribute) calculated by modified DNSName. dash replace dots
}

var zoneInfo ZoneInfo

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

var dnsRecords []DNSRecord

func generate() error {
	var err error

	err = processZoneFile()
	if err != nil {
		return err
	}

	err = processRecordsFile()
	if err != nil {
		return err
	}
	err = genTF()
	if err != nil {
		return err
	}

	// sort recods for documentation
	sort.SliceStable(dnsRecords, func(i, j int) bool {
		if dnsRecords[i].Type != dnsRecords[j].Type {
			return dnsRecords[i].Type < dnsRecords[j].Type
		}
		return dnsRecords[i].HostName < dnsRecords[j].HostName
	})

	err = genMD()
	if err != nil {
		return err
	}
	fmt.Println("  Generation Complete")
	return nil
}

func readCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileReader := csv.NewReader(file)
	records, err := fileReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func processZoneFile() error {
	const DNSNameIDX = 0
	const ProjectIDIDX = 1
	const DescriptionIDX = 2

	// read zone info file
	records, err := readCSV(zoneInfoFile)
	if err != nil {
		return err
	}
	if len(records) > 1 {
		return errors.New("Should only be 1 records in zone file.")
	}
	zoneInfo.DNSName = records[0][DNSNameIDX]

	// May or maynot have a trailing '.' but triming anyway.
	zoneInfo.DNSName = strings.TrimSuffix(zoneInfo.DNSName, ".")
	zoneInfo.ResourceName = strings.ReplaceAll(zoneInfo.DNSName, ".", "-")
	zoneInfo.ProjectID = records[0][ProjectIDIDX]
	zoneInfo.Description = records[0][DescriptionIDX]
	return nil
}

func processRecordsFile() error {
	const HostNameIDX = 0
	const TypeIDX = 1
	const TTLIDX = 2
	const RRDatasIDX = 3
	const DescriptionIDX = 4

	// read zone info file
	records, err := readCSV(recordsFile)
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record) != 5 {
			return errors.New("Column error in zone record FileName: " + recordsFile)
		}
		dnsRecord := DNSRecord{
			HostName:    record[HostNameIDX],
			Type:        record[TypeIDX],
			TTL:         record[TTLIDX],
			RRDatas:     record[RRDatasIDX],
			Description: record[DescriptionIDX],
		}

		// Should not have a trailing '.' but trim anyway.
		dnsRecord.HostName = strings.TrimPrefix(dnsRecord.HostName, ".")

		// Escape rrdatas if contain spaces
		if strings.Contains(dnsRecord.RRDatas, " ") {
			dnsRecord.RRDatas = "\\\"" + dnsRecord.RRDatas + "\\\""
		}
		dnsRecord.ProjectID = zoneInfo.ProjectID
		dnsRecord.ManagedZone = zoneInfo.ResourceName
		dnsRecord.Name = dnsRecord.HostName + "." + zoneInfo.DNSName
		dnsRecord.ResourceName = strings.ReplaceAll(dnsRecord.Name, ".", "-")
		dnsRecords = append(dnsRecords, dnsRecord)
	}

	return nil
}

func genTF() error {
	var err error
	var tfFile *os.File

	tfFile, err = os.Create(tfFileName)
	if err != nil {
		return err
	}
	defer tfFile.Close()
	// Dump out header

	fmt.Fprintf(tfFile, "%s\n", tfHeader)

	// Create Zone definition

	zoneTemplate := template.New("zoneTemplate")
	_, err = zoneTemplate.Parse(zoneTemplateDef)
	if err != nil {
		return err
	}
	zoneTemplate.Execute(tfFile, zoneInfo)

	// Create zone records

	zoneRecordTemplate := template.New("zoneRecordTemplate")
	_, err = zoneRecordTemplate.Parse(zoneRecordTemplateDef)
	if err != nil {
		return err
	}

	for _, v := range dnsRecords {
		zoneRecordTemplate.Execute(tfFile, v)
	}
	return nil
}

func genMD() error {
	var err error
	var mdFile *os.File

	mdFile, err = os.Create(mdFileName)
	if err != nil {
		return err
	}
	defer mdFile.Close()

	// dump out zone information header
	mdZoneHeaderTemplate := template.New("mdZoneHeader")
	_, err = mdZoneHeaderTemplate.Parse(mdZoneHeaderTemplateDef)
	if err != nil {
		return err
	}
	mdZoneHeaderTemplate.Execute(mdFile, zoneInfo)

	// dump out zone record information
	mdZoneRecordTemplate := template.New("mdZoneRecordTemplate")
	_, err = mdZoneRecordTemplate.Parse(mdZoneRecordTemplateDef)
	if err != nil {
		return err
	}

	for _, v := range dnsRecords {
		mdZoneRecordTemplate.Execute(mdFile, v)
	}
	return nil
}
