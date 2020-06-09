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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"terraform-dns/gen"

	"github.com/spf13/cobra"
)

var (
	jsonInfoFile string
)

// tfCreateCmd represents the tfCreate command
var jsonCreateCmd = &cobra.Command{
	Use:   "jsonCreate",
	Short: "Create DNS terraform files from json file..",
	Long:  `Used to create DNS terraform file for a zone from a json file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return genFromJSON()
	},
}

func init() {
	rootCmd.AddCommand(jsonCreateCmd)

	jsonCreateCmd.PersistentFlags().StringVar(&jsonInfoFile, "zone-info", "zone.json", "DNS configuration file (zone.json)")
}

func genFromJSON() error {
	var err error

	err = readJSON()
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

func readJSON() error {
	var err error
	var rawJSON []byte

	rawJSON, err = ioutil.ReadFile(jsonInfoFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(rawJSON, &gen.DNSZone)
	if err != nil {
		return err
	}

	return nil
}
