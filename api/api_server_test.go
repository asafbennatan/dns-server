package api

import (
	"bytes"
	"context"
	"dnsServer/daos"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Start the DNS server and get the stop channel
	server := StartApiServer(":8080")
	defer server.Shutdown(context.Background())

	// Wait a bit to ensure the server is ready
	time.Sleep(time.Second)
	code := m.Run()

	os.Exit(code)

}

func TestZone(t *testing.T) {
	var createdZone daos.DNSZone
	t.Run("CreateZone", func(t *testing.T) {
		newZone := daos.DNSZoneCreate{Name: uuid.NewString() + ".com"}
		body, _ := json.Marshal(newZone)
		resp, err := http.Post("http://localhost:8080/api/zone", "application/json", bytes.NewReader(body))
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to create zone, err: %v, status code: %v", err, resp.StatusCode)
		}
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&createdZone)
		if createdZone.Name != newZone.Name {
			t.Errorf("Expected zone name %v, got %v", newZone.Name, createdZone.Name)
		}
	})

	t.Run("GetAllZones", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/api/zone")
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to get zones, err: %v, status code: %v", err, resp.StatusCode)
		}
		defer resp.Body.Close()
		var zones []daos.DNSZone
		json.NewDecoder(resp.Body).Decode(&zones)
		found := false
		for _, zone := range zones {
			if zone.ID == createdZone.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Created zone not found in list")
		}
	})

	t.Run("UpdateZone", func(t *testing.T) {
		client := &http.Client{}

		newZone := daos.DNSZoneUpdate{
			ID:            createdZone.ID,
			DNSZoneCreate: daos.DNSZoneCreate{Name: uuid.NewString() + ".com"},
		}
		body, _ := json.Marshal(newZone)

		req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/api/zone", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)

		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to create zone, err: %v, status code: %v", err, resp.StatusCode)
		}
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&createdZone)
		if createdZone.Name != newZone.Name {
			t.Errorf("Expected zone name %v, got %v", newZone.Name, createdZone.Name)
		}
	})
	t.Run("GetZone", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/api/zone/" + createdZone.ID)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to get zones, err: %v, status code: %v", err, resp.StatusCode)
		}
		defer resp.Body.Close()
		var fetched daos.DNSZone
		json.NewDecoder(resp.Body).Decode(&fetched)

		if fetched.ID != createdZone.ID {
			t.Errorf("Expected zone id %v, got %v", createdZone.ID, fetched.ID)
		}
	})

	//do the same for record , use the created zone to create a record
	var createdRecord daos.DNSRecord

	t.Run("CreateRecord", func(t *testing.T) {
		newRecord := daos.DNSRecordCreate{
			Name:  "www",
			Type:  "A",
			Value: "1.2.3.4",
			TTL:   300,
		}
		body, _ := json.Marshal(newRecord)
		resp, err := http.Post(fmt.Sprintf("http://localhost:8080/api/zone/%s/record", createdZone.ID), "application/json", bytes.NewReader(body))
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to create record, err: %v, status code: %v", err, resp.StatusCode)

		}
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&createdRecord)
		validateRecordCreate(t, createdRecord, newRecord, createdZone.ID)
	})

	t.Run("GetAllRecords", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/zone/%s/record", createdRecord.DNSZoneID))
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to get records, err: %v, status code: %v", err, resp.StatusCode)
		}
		defer resp.Body.Close()
		var records []daos.DNSRecord
		json.NewDecoder(resp.Body).Decode(&records)
		found := false
		for _, record := range records {
			if record.ID == createdRecord.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Created record not found in list")
		}
	})

	t.Run("UpdateRecord", func(t *testing.T) {
		client := &http.Client{}

		newRecord := daos.DNSRecordUpdate{
			ID: createdRecord.ID,
			DNSRecordCreate: daos.DNSRecordCreate{
				Name:  "www",
				Type:  "A",
				Value: "2.3.4.5",
				TTL:   300,
			},
		}
		body, _ := json.Marshal(newRecord)

		req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/api/record", bytes.NewReader(body))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)

		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to create record, err: %v, status code: %v", err, resp.StatusCode)
		}
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&createdRecord)
		validateRecordCreate(t, createdRecord, newRecord.DNSRecordCreate, createdRecord.DNSZoneID)

	})

	t.Run("GetRecord", func(t *testing.T) {
		url := fmt.Sprintf("http://localhost:8080/api/record/%s", createdRecord.ID)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to get record, err: %v, status code: %v", err, resp.StatusCode)
		}
		defer resp.Body.Close()
		var fetched daos.DNSRecord
		json.NewDecoder(resp.Body).Decode(&fetched)

		validateRecord(t, createdRecord, fetched)
	})

	t.Run("DeleteRecord", func(t *testing.T) {

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/record/"+createdRecord.ID, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to delete record, err: %v, status code: %v", err, resp.StatusCode)
		}
	})

	t.Run("GetRecordAfterDelete", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/api/record/" + createdRecord.ID)
		if err != nil || resp.StatusCode != http.StatusNotFound {
			t.Fatalf("got record,expected not found err: %v, status code: %v", err, resp.StatusCode)
		}
		defer resp.Body.Close()

	})

	t.Run("DeleteZone", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/zone/"+createdZone.ID, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("Failed to delete zone, err: %v, status code: %v", err, resp.StatusCode)
		}
	})

	t.Run("GetZoneAfterDelete", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/api/zone/" + createdZone.ID)
		if err != nil || resp.StatusCode != http.StatusNotFound {
			t.Fatalf("got zone,expected not found err: %v, status code: %v", err, resp.StatusCode)
		}
		defer resp.Body.Close()

	})

}

func validateRecordCreate(t *testing.T, createdRecord daos.DNSRecord, newRecord daos.DNSRecordCreate, expectedZoneId string) {
	if createdRecord.Name != newRecord.Name {
		t.Errorf("Expected record name %v, got %v", newRecord.Name, createdRecord.Name)
	}
	if createdRecord.TTL != newRecord.TTL {
		t.Errorf("Expected record TTL %v, got %v", newRecord.TTL, createdRecord.TTL)
	}
	if createdRecord.Type != newRecord.Type {
		t.Errorf("Expected record Type %v, got %v", newRecord.Type, createdRecord.Type)
	}
	if createdRecord.Value != newRecord.Value {
		t.Errorf("Expected record Value %v, got %v", newRecord.Value, createdRecord.Value)
	}
	if createdRecord.DNSZoneID != expectedZoneId {
		t.Errorf("Expected record ZoneId %v, got %v", expectedZoneId, createdRecord.DNSZoneID)
	}
}

func validateRecord(t *testing.T, createdRecord daos.DNSRecord, newRecord daos.DNSRecord) {
	if createdRecord.Name != newRecord.Name {
		t.Errorf("Expected record name %v, got %v", newRecord.Name, createdRecord.Name)
	}
	if createdRecord.TTL != newRecord.TTL {
		t.Errorf("Expected record TTL %v, got %v", newRecord.TTL, createdRecord.TTL)
	}
	if createdRecord.Type != newRecord.Type {
		t.Errorf("Expected record Type %v, got %v", newRecord.Type, createdRecord.Type)
	}
	if createdRecord.Value != newRecord.Value {
		t.Errorf("Expected record Value %v, got %v", newRecord.Value, createdRecord.Value)
	}
	if createdRecord.DNSZoneID != newRecord.DNSZoneID {
		t.Errorf("Expected record ZoneId %v, got %v", newRecord.DNSZoneID, createdRecord.DNSZoneID)
	}
}
