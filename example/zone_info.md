
# roger-beathard.com
ProjectID: my-dns-279215  
Description: zone file for roger-beathard.com  

---
## Zone Records

host|type|ttl|rrdatas|description
----|----|---|-------|-----------
apple|A|300|10.1.5.6|Main home web service
aweb|A|300|10.1.4.3|Main home web service
home-green|A|300|10.1.4.3|home-green web service
home-red|A|300|10.1.4.3|home-red web service
www|A|300|10.1.2.1|Default web entry point
home-lb|CNAME|300|home-red.rvbtech.com.|Main home web service
home2|CNAME|300|home-red.vbtech.com.|Main home web service
azure|TXT|300|\"v=spf1 ip4:111.111.111.111 include:backoff.example.com -all\"|Main home web service