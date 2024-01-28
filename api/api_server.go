package api

import (
	"context"
	"dnsServer/data"
	"dnsServer/service"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func main() {
	r := mux.NewRouter()
	db := data.InitDB()
	// Create service instances
	zoneService := service.NewZoneService(db)
	recordService := service.NewRecordService(db) //
	r.Use(
		injectService("zoneService", zoneService),
		injectService("recordService", recordService),
	)

	// Set up routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/zone", createZone).Methods(http.MethodPost)
	api.HandleFunc("/zone", getZones).Methods(http.MethodGet)
	api.HandleFunc("/zone/{id}", getZone).Methods(http.MethodGet)
	api.HandleFunc("/zone/{id}", updateZone).Methods(http.MethodPut)
	api.HandleFunc("/zone/{id}", deleteZone).Methods(http.MethodDelete)
	api.HandleFunc("/zone/{zone_id}/record", createRecord).Methods(http.MethodPost)
	api.HandleFunc("/zone/{zone_id}/record", getRecords).Methods(http.MethodGet)
	api.HandleFunc("/record/{id}", getRecord).Methods(http.MethodGet)
	api.HandleFunc("/record/{id}", updateRecord).Methods(http.MethodPut)
	api.HandleFunc("/record/{id}", deleteRecord).Methods(http.MethodDelete)

	// Start the HTTP server
	http.ListenAndServe(":8080", r)
}

func injectService(key string, service any) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set the DB connection in the context
			ctx := context.WithValue(r.Context(), key, service)
			// Call the next handler with the new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func createZone(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var data DNSZoneCreate
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
		return
	}
	fmt.Printf("Received data: %+v\n", data)

	zoneService, ok := r.Context().Value("zoneService").(*service.ZoneService)
	if !ok {
		http.Error(w, "Could not get database connection", http.StatusInternalServerError)
		return
	}
	zone := zoneService.CreateZone(data)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zone)
}

func updateZone(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var data DNSZoneUpdate
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
		return
	}
	zoneService, ok := r.Context().Value("zoneService").(*service.ZoneService)
	if !ok {
		http.Error(w, "Could not get database connection", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received data: %+v\n", data)
	zone := zoneService.UpdateZone(data)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zone)
}

func getZone(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	zoneService, ok := r.Context().Value("zoneService").(*service.ZoneService)
	if !ok {
		http.Error(w, "Could not get database connection", http.StatusInternalServerError)
		return
	}
	zone := zoneService.GetZone(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zone)
}

func deleteZone(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	zoneService, ok := r.Context().Value("zoneService").(*service.ZoneService)
	if !ok {
		http.Error(w, "Could not get database connection", http.StatusInternalServerError)
		return
	}
	zoneService.DeleteZone(id)

}

func getZones(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	zoneService, ok := r.Context().Value("zoneService").(*service.ZoneService)
	if !ok {
		http.Error(w, "Could not get database connection", http.StatusInternalServerError)
		return
	}
	zone := zoneService.GetZones("%%")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zone)
}

func createRecord(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	vars := mux.Vars(r)
	zoneId := vars["zone_id"]
	var data DNSRecordCreate
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received data: %+v\n", data)
	recordService, ok := r.Context().Value("recordService").(*service.RecordService)
	if !ok {
		http.Error(w, "Could not get database connection", http.StatusInternalServerError)
		return
	}
	record := recordService.CreateRecord(uint(strconv.ParseInt(zoneId, 10, 32)), data)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

func updateRecord(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var data DNSRecordUpdate
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received data: %+v\n", data)
	recordService, ok := r.Context().Value("recordService").(*service.RecordService)
	if !ok {
		http.Error(w, "Could not get database connection", http.StatusInternalServerError)
		return
	}
	record := recordService.UpdateRecord(data)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

func getRecord(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	recordService, ok := r.Context().Value("recordService").(*service.RecordService)
	if !ok {
		http.Error(w, "Could not get database connection", http.StatusInternalServerError)
		return
	}
	record := recordService.GetRecord(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	recordService, ok := r.Context().Value("recordService").(*service.RecordService)
	if !ok {
		http.Error(w, "Could not get database connection", http.StatusInternalServerError)
		return
	}
	recordService.DeleteRecord(id)

}

func getRecords(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	vars := mux.Vars(r)
	zoneId := vars["zone_id"]
	recordService, ok := r.Context().Value("recordService").(*service.RecordService)
	if !ok {
		http.Error(w, "Could not get database connection", http.StatusInternalServerError)
		return
	}
	record := recordService.GetRecords(zoneId, "%%")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(record)
}
