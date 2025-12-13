# IT Consultant Microservices

A modern microservices architecture for an IT consultant website, built with **Go** and containerized with **Docker**. Features a clean separation between a JSON API backend and a dynamic web frontend.

## Overview

This project demonstrates microservices best practices by splitting functionality into two independent services:

| Service | Purpose | Technology |
|---------|---------|-----------|
| **API Service** | Provides JSON endpoints for profile, services, and contact management | Go, net/http |
| **Web Service** | Server-side rendered frontend with template-based HTML | Go, html/template |

Both services communicate over a Docker network and are orchestrated together using `docker-compose`.

## Architecture

```
┌─────────────┐
│  Web Client │
└──────┬──────┘
       │ HTTP
       ▼
┌──────────────────────┐
│  Web Service (8080)  │
│  - HTML Templates    │────────┐
│  - Static CSS        │        │ HTTP (over Docker network)
│  - API Consumer      │        │
└──────────────────────┘        │
                                 ▼
                        ┌──────────────────────┐
                        │  API Service (8081)  │
                        │  - Profile Data      │
                        │  - Services List     │
                        │  - Contact Handler   │
                        └──────────────────────┘
```

## Project Structure

```
.
├── api-service/                 # JSON API microservice
│   ├── main.go                  # API handlers and routes
│   ├── go.mod                   # Go module definition
│   └── Dockerfile               # Container image definition
├── web-service/                 # Web frontend microservice
│   ├── main.go                  # Web handlers and routing
│   ├── go.mod                   # Go module definition
│   ├── Dockerfile               # Container image definition
│   ├── templates/               # Server-side HTML templates
│   │   ├── base.html
│   │   ├── home.html
│   │   ├── about.html
│   │   ├── services.html
│   │   └── contact.html
│   └── static/                  # Static assets
│       └── css/
│           └── style.css        # Main stylesheet
├── K8S/                         # Kubernetes configuration (optional)
│   ├── base/                    # Base resources
│   │   ├── api/
│   │   └── web/
│   └── overlays/                # Environment-specific overrides
│       ├── dev/
│       └── prod/
└── docker-compose.yml           # Local development orchestration
```

## Prerequisites

- **Docker** and **Docker Compose** (for containerized deployment)
- **Go 1.19+** (optional, for local development)

## Quick Start

### Using Docker Compose (Recommended)

```bash
docker-compose up --build
```

The application will be available at:
- **Web Frontend**: http://localhost:8080
- **API Service**: http://localhost:8081

### Local Development (without Docker)

1. **API Service**:
```bash
cd api-service
go run main.go
# API will run on http://localhost:8081
```

2. **Web Service** (in a new terminal):
```bash
cd web-service
go run main.go
# Web will run on http://localhost:8080
```

## API Reference

### API Service Endpoints

#### GET /api/profile
Returns the consultant's profile information.

**Response**:
```json
{
  "name": "Your Name",
  "title": "IT Consultant",
  "tagline": "Your professional tagline",
  "summary": "Brief bio...",
  "technologies": ["Go", "Docker", "Kubernetes"],
  "location": "City, Country"
}
```

#### GET /api/services
Returns a list of consulting services offered.

**Response**:
```json
{
  "items": [
    {
      "slug": "service-id",
      "name": "Service Name",
      "description": "Service description...",
      "category": "Category"
    }
  ]
}
```

#### POST /api/contact
Accepts contact form submissions (demo endpoint).

**Request Body**:
```json
{
  "name": "Client Name",
  "email": "client@example.com",
  "company": "Company Name",
  "message": "Message content..."
}
```

**Response**: `202 Accepted` (logs the payload)

### Web Service Endpoints

| Path | Description |
|------|-------------|
| `GET /` | Home page with profile and featured services |
| `GET /services` | Complete services catalog |
| `GET /about` | About page |
| `GET /contact` | Contact form (submits to `/api/contact`) |

## Configuration

### Docker Compose Environment

The `docker-compose.yml` defines two services with the following configuration:

- **API Service**: Runs on port `8081`, accessible at hostname `api` within the Docker network
- **Web Service**: Runs on port `8080`, depends on API service being available
- Both services restart automatically unless manually stopped

### Updating Content

- **Profile Data**: Edit the profile data in `api-service/main.go`
- **Services**: Modify service definitions in `api-service/main.go`
- **Styling**: Update `web-service/static/css/style.css`
- **Templates**: Modify HTML templates in `web-service/templates/`

## Kubernetes Deployment (Optional)

A Kustomize-based Kubernetes configuration is included for production deployments:

```bash
# Dev environments
kubectl apply -k K8S/overlays/dev

# Production environment
kubectl apply -k K8S/overlays/prod
```

## Development Tips

### Service Communication

The web service communicates with the API service using the hostname `api` (the service name in `docker-compose.yml`). This works due to Docker's built-in DNS resolution within the network.

### Logs

View service logs:
```bash
docker-compose logs -f api    # API service logs
docker-compose logs -f web    # Web service logs
docker-compose logs           # All logs
```

### Stopping Services

```bash
docker-compose down           # Stop and remove containers
docker-compose down -v        # Also remove volumes
```

## Customization

This is a template project designed to be easily customized:

1. **Brand & Content**: Update profile info, services, and copy in the code
2. **Styling**: Modify `static/css/style.css` for custom branding
3. **Templates**: Edit HTML templates in `templates/` to change layout
4. **Functionality**: Extend API endpoints and web handlers as needed

## License

This project is open source. Customize and deploy as needed for your consulting business.

## Support

For questions or issues, refer to the service-specific code in `api-service/` and `web-service/` directories.
