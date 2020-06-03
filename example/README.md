# terraform-dns example

Contents of this directory:
- zone.csv - example zone configuration
- records.csv - example zone records configuration 
- main.tf - example generated terraform configuration
- zone-info.md - [example generated documentation set for zone](zone_info.md)

---

**Example Configuration**

The zone configuration is specified in the zone.csv. For this example, the zone.csv will generate terraform module definition for "roger.beathard.com" in the mydns project.

The records for the zone is specified in the records.csv file. For this example there are 8 entries: "www", "azure", "home-red", "home-green", "aweb", "apple", "home-lb", & "home2".

**Re-generation of main.tf and zone-info.md**

In order to re-generate these files run ../terraform-dns tf-create.

**Example documenation set**
