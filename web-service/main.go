package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type Profile struct {
	Name         string   `json:"name"`
	Title        string   `json:"title"`
	Tagline      string   `json:"tagline"`
	Summary      string   `json:"summary"`
	Technologies []string `json:"technologies"`
	Location     string   `json:"location"`
}

type Service struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type ServicesResponse struct {
	Items []Service `json:"items"`
}

type PageData struct {
	Title       string
	Tagline     string
	Description string
	Profile     *Profile
	Services    []Service
	Error       string
}

var templates *template.Template
var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func main() {
	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/services", servicesHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/contact", contactHandler)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	addr := ":8080"
	log.Printf("Starting Web (frontend) service on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Web server failed: %v", err)
	}
}

func renderTemplate(w http.ResponseWriter, name string, data PageData) {
	tmpl := templates.Lookup(filepath.Base(name))
	if tmpl == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template %s: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	profile, services, err := loadCoreData()
	data := PageData{
		Title:       "Cloud & DevOps Consulting",
		Tagline:     "Modern cloud solutions for resilient, secure and scalable systems.",
		Description: "I help teams design, build, and operate cloud infrastructure on Azure and AWS with a strong focus on Kubernetes, Terraform, automation, and observability.",
		Profile:     profile,
		Services:    services,
	}

	if err != nil {
		data.Error = "Some data could not be loaded from the API. Please try again later."
	}

	renderTemplate(w, "home.html", data)
}

func servicesHandler(w http.ResponseWriter, r *http.Request) {
	_, services, err := loadCoreData()
	data := PageData{
		Title:       "Services",
		Tagline:     "End-to-end consulting across Cloud, DevOps, and Platform Engineering.",
		Description: "From greenfield builds to rescuing existing platforms, I partner with teams to design and implement pragmatic cloud solutions.",
		Services:    services,
	}

	if err != nil {
		data.Error = "Unable to load services from the API."
	}

	renderTemplate(w, "services.html", data)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	profile, _, err := loadCoreData()
	data := PageData{
		Title:       "About",
		Tagline:     "Hands-on Cloud Engineer & IT Consultant.",
		Description: "I work with organizations to modernize infrastructure, improve reliability, and ship faster using proven DevOps practices.",
		Profile:     profile,
	}

	if err != nil {
		data.Error = "Unable to load profile information from the API."
	}

	renderTemplate(w, "about.html", data)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:       "Contact",
		Tagline:     "Let’s talk about your infrastructure and delivery challenges.",
		Description: "Use the form below to describe your current environment and what you’d like to improve. I’ll respond with practical next steps.",
	}
	renderTemplate(w, "contact.html", data)
}

// loadCoreData fetches profile and services from the API service over the Docker network.
func loadCoreData() (*Profile, []Service, error) {
	profile, err := fetchProfile()
	if err != nil {
		log.Printf("Error fetching profile: %v", err)
	}
	services, err2 := fetchServices()
	if err2 != nil {
		log.Printf("Error fetching services: %v", err2)
	}
	if err != nil || err2 != nil {
		if err != nil {
			return profile, services, err
		}
		return profile, services, err2
	}
	return profile, services, nil
}

func fetchProfile() (*Profile, error) {
	req, err := http.NewRequest(http.MethodGet, "http://api:8081/api/profile", nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func fetchServices() ([]Service, error) {
	req, err := http.NewRequest(http.MethodGet, "http://api:8081/api/services", nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sr ServicesResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, err
	}

	return sr.Items, nil
}
