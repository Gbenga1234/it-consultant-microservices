package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Profile struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Tagline     string   `json:"tagline"`
	Summary     string   `json:"summary"`
	Technologies []string `json:"technologies"`
	Location    string   `json:"location"`
}

type Service struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
	Message string `json:"message"`
}

type ServicesResponse struct {
	Items []Service `json:"items"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/profile", handleProfile)
	mux.HandleFunc("/api/services", handleServices)
	mux.HandleFunc("/api/contact", handleContact)

	addr := ":8081"
	log.Printf("Starting API service on %s\n", addr)
	if err := http.ListenAndServe(addr, withCORS(mux)); err != nil {
		log.Fatalf("API server failed: %v", err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	profile := Profile{
		Name:  "Tosin Femi",
		Title: "Cloud & DevOps Consultant",
		Tagline: "Modern cloud solutions for resilient, secure and scalable systems.",
		Summary: "I help teams design, build, and operate cloud infrastructure on Azure and AWS with a strong focus on Kubernetes, Terraform, automation, and observability.",
		Technologies: []string{"Azure", "AWS", "Kubernetes (AKS/EKS)", "Terraform", "GitHub Actions", "Azure DevOps", "Docker", "Argo CD"},
		Location: "Calgary · Remote-friendly",
	}

	writeJSON(w, profile, http.StatusOK)
}

func handleServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	services := []Service{
		{
			Slug:        "cloud-architecture",
			Name:        "Cloud Architecture & Migration",
			Description: "Azure & AWS landing zones, workload migration plans, and security-focused network design.",
			Category:    "Cloud",
		},
		{
			Slug:        "kubernetes-platform",
			Name:        "Kubernetes & Platform Engineering",
			Description: "Production-ready AKS/EKS clusters, GitOps workflows, and multi-environment platform design.",
			Category:    "Kubernetes",
		},
		{
			Slug:        "devops-automation",
			Name:        "DevOps & Automation",
			Description: "Infrastructure as Code, CI/CD pipelines, and release automation for faster, safer delivery.",
			Category:    "DevOps",
		},
		{
			Slug:        "observability-reliability",
			Name:        "Observability & Reliability",
			Description: "Logging, metrics, alerting, and SLOs so your team can detect and fix issues quickly.",
			Category:    "Operations",
		},
	}

	writeJSON(w, ServicesResponse{Items: services}, http.StatusOK)
}

func handleContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("New contact request: %+v\n", payload)

	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(`{"status":"accepted","message":"Thanks, your request has been received."}`))
}

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Error writing JSON: %v", err)
	}
}
