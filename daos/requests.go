package daos

type DNSRecordCreate struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

type DNSRecordUpdate struct {
	DNSRecordCreate
	ID string `json:"id"`
}

type DNSZoneCreate struct {
	Name string `json:"name"`
}

type DNSZoneUpdate struct {
	DNSZoneCreate
	ID string `json:"id"`
}
