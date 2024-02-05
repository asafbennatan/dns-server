package daos

type DNSRecord struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	TTL       int    `json:"ttl"`
	DNSZoneID string `json:"dnsZoneID"`
}

type DNSZone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
