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
	"fmt"
	"os"
	"text/template"
)

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

// GenTF - Generate the terraform file.
func GenTF(fileName string) error {
	var err error
	var tfFile *os.File

	tfFile, err = os.Create(fileName)
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
	zoneTemplate.Execute(tfFile, DNSZone)

	// Create zone records

	zoneRecordTemplate := template.New("zoneRecordTemplate")
	_, err = zoneRecordTemplate.Parse(zoneRecordTemplateDef)
	if err != nil {
		return err
	}

	for _, v := range DNSZone.DNSRecords {
		zoneRecordTemplate.Execute(tfFile, v)
	}
	return nil
}
