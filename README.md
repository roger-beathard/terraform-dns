
## terraform-dns
Generate GCP DNS terraform files from csv file using golang.

------

**Purpose**

terraform code for managing DNS records can be daunting due to the fact that each zone record is a separate module definition. Managing records in a simple text file makes the task easier.

One of the other advantages is that we can now generate a documentation set that describes the DNS zone.

**How To Run**
- go build
- navigate to the example directory which contains sample records and zone csv files.
- ../terraform-dns tf-create

By default terraform-dns tf-create file will read two csv files zone.csv & records.csv and output the main.tf (terraform file) and zone markdown documentation set (zone_info.md). Please refer to "terraform-dns tf-create --help" on how to modify these parameters.

**Example deployment flow for new DNS records**

- New DNS record is inserted/updated in records.csv file
- terraform and documentation file is generated via terraform-dns tf-create.
- "terraform apply main.tf" is used to push changes to cloud

**zone.csv file**

The format of the zone.csv file is a standard csv file that has no header and is formated as follows:

"dns name of zone", "GCP ProjectID", "Description" 

**records.csv file**

The format of the records.csv file is a standard csv file that has no header and is formated as follows:

"hostname", "record type", "ttl", "rrdatas", "description"

*Note:* The rrdatas do not need to be escapsed as per terraform requirements. The rrdatas string is analyized for spaces and will escape the string if spaces are found.

**Caveats**

Currently there is a timing issue on initial creation of a zone which will result in an error when the terraform script is first applied. This is due to the fact that the zone creation will return before the zone is created and the subsequent application of zone records will result in error.




