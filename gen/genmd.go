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
	"os"
	"text/template"
)

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

// GenMD - Generate the markdowdn documentation
func GenMD(fileName string) error {
	var err error
	var mdFile *os.File

	mdFile, err = os.Create(fileName)
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
	mdZoneHeaderTemplate.Execute(mdFile, DNSZone)

	// dump out zone record information
	mdZoneRecordTemplate := template.New("mdZoneRecordTemplate")
	_, err = mdZoneRecordTemplate.Parse(mdZoneRecordTemplateDef)
	if err != nil {
		return err
	}

	for _, v := range DNSZone.DNSRecords {
		mdZoneRecordTemplate.Execute(mdFile, v)
	}
	return nil
}
