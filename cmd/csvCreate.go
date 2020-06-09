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
package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"terraform-dns/gen"

	"github.com/spf13/cobra"
)

var (
	zoneInfoFile string
	recordsFile  string
)

// tfCreateCmd represents the tfCreate command
var tfCreateCmd = &cobra.Command{
	Use:   "csvCreate",
	Short: "Create DNS terraform files.",
	Long:  `Used to create DNS terraform files from a csv encoded zone and records file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return genFromCSV()
	},
}

func init() {
	rootCmd.AddCommand(tfCreateCmd)

	tfCreateCmd.PersistentFlags().StringVar(&zoneInfoFile, "zone-info", "zone.csv", "zone info file (zone.csv)")
	tfCreateCmd.PersistentFlags().StringVar(&recordsFile, "records", "records.csv", "zone records file (records.csv)")
}

func genFromCSV() error {
	var err error

	err = processZoneFile()
	if err != nil {
		return err
	}

	err = processRecordsFile()
	if err != nil {
		return err
	}

	// Generate computed values && sort
	gen.DNSZone.ComputeValues()
	gen.DNSZone.SortRecords()

	err = gen.GenTF(tfFileName)
	if err != nil {
		return err
	}

	err = gen.GenMD(mdFileName)
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
	gen.DNSZone.DNSName = records[0][DNSNameIDX]
	gen.DNSZone.ProjectID = records[0][ProjectIDIDX]
	gen.DNSZone.Description = records[0][DescriptionIDX]
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
		dnsRecord := gen.DNSRecord{
			HostName:    record[HostNameIDX],
			Type:        record[TypeIDX],
			TTL:         record[TTLIDX],
			RRDatas:     record[RRDatasIDX],
			Description: record[DescriptionIDX],
		}
		gen.DNSZone.DNSRecords = append(gen.DNSZone.DNSRecords, dnsRecord)
	}

	return nil
}
