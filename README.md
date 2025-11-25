# IT Consultant Website - Go Microservices

This project splits the IT consultant website into **two Go microservices**:

- `api-service` – JSON API providing profile, services, and a demo contact endpoint
- `web-service` – Frontend website in Go (server-side templates) that calls the API service

Everything is containerised with Docker and wired together via `docker-compose`.

## Structure

- `api-service/`
  - `main.go`
  - `go.mod`
  - `Dockerfile`
- `web-service/`
  - `main.go`
  - `go.mod`
  - `Dockerfile`
  - `templates/`
  - `static/css/style.css`
- `docker-compose.yml`

## Running with Docker Compose

```bash
docker-compose up --build
```

- Web frontend: http://localhost:8080
- API service: http://localhost:8081

The web service calls the API over the Docker network using the hostname `api` (the service name in `docker-compose.yml`).

## Endpoints

### API Service

- `GET /api/profile` – basic consultant profile
- `GET /api/services` – list of consulting services
- `POST /api/contact` – demo contact endpoint (logs payload, returns 202 Accepted)

### Web Service

- `GET /` – Home page (profile + highlighted services)
- `GET /services` – Services page
- `GET /about` – About page
- `GET /contact` – Contact page (demo form that POSTs to the API on `http://localhost:8081/api/contact`)

You can customise the content, styling, and behaviour as needed for your own consulting site.
