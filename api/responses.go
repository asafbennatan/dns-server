package api

type DNSRecord struct {
	ID      uint    `json:"id"`
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Value   string  `json:"value"`
	TTL     int     `json:"ttl"`
	DNSZone DNSZone `json:"dnsZone"`
}

type DNSZone struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
