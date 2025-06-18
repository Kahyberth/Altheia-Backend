# 🏥 Altheia Backend

**Altheia** is a complete backend system for an Electronic Health Record (EHR) application developed in Go. This system provides a secure, scalable, and modern REST API designed to facilitate comprehensive healthcare services management.

[![Go Version](https://img.shields.io/badge/Go-1.24-blue)](https://golang.org/)
[![Fiber Framework](https://img.shields.io/badge/Fiber-v2.52.6-00ADD8)](https://github.com/gofiber/fiber)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

## 📋 Description

Altheia Backend is a comprehensive solution that enables complete medical services management, including:
- Management of patients, physicians, receptionists, and clinic owners
- Digital medical records and consultations
- Medical appointment system
- Clinic and services administration
- Robust authentication and authorization with JWT

## ✨ Key Features

### 🔐 Authentication and Security
- JWT authentication with access and refresh tokens
- Route protection middleware
- Login activity management
- Password encryption
- Role-based access control

### 👥 User Management
- **Patients**: Registration, updates, soft deletion
- **Physicians**: Complete management with pagination
- **Receptionists**: Reception staff administration
- **Clinic Owners**: Medical center owner management
- **Laboratory Technicians**: Technical staff support

### 🏥 Clinical Management
- Clinic registration and administration
- Detailed information for each medical center
- Schedules and offered services
- Staff assignment to clinics
- Integration with EPS (Health Promotion Entities)

### 📋 Medical Records
- Creation and updating of medical records
- Detailed medical consultations
- Medical prescriptions
- Treatment tracking

### 📅 Appointment System
- Medical appointment scheduling
- Appointment status management
- Filtering by physician and date
- Notifications and reminders

## 🛠️ Technologies

- **[Go](https://golang.org/)** 1.24 - Programming language
- **[Fiber](https://github.com/gofiber/fiber)** v2 - Ultra-fast web framework
- **[GORM](https://gorm.io/)** - ORM for Go
- **[PostgreSQL](https://www.postgresql.org/)** 16 - Relational database
- **[JWT](https://github.com/golang-jwt/jwt)** - Token-based authentication
- **[Docker](https://www.docker.com/)** - Containerization
- **[bcrypt](https://golang.org/x/crypto)** - Password encryption

## 🚀 Installation and Setup

### Prerequisites

- Go 1.24 or higher
- PostgreSQL 16
- Docker and Docker Compose (optional)

### 1. Clone the Repository

```bash
git clone https://github.com/Kahyberth/Altheia-Backend.git
cd Altheia-Backend
```

### 2. Docker Installation (Recommended)

```bash
# Start PostgreSQL database
docker-compose up -d

# Install dependencies
make deps

# Run the application
make build
./bin/altheia
```

### 4. Manual Installation

```bash
# Install dependencies
go mod download

# Build the application
go build -o bin/altheia ./cmd/main.go
```

## 📖 API Endpoints

### Authentication
- `POST /auth/login` - Login
- `POST /auth/logout` - Logout
- `GET /auth/verify-token` - Verify token
- `POST /auth/refresh-token/:refresh_token` - Refresh token

### Patients
- `POST /patient/register` - Register patient
- `GET /patient/getAll` - Get all patients
- `GET /patient/getAllPaginated` - Get paginated patients
- `PATCH /patient/update/:id` - Update patient

### Physicians
- `POST /physician/register` - Register physician
- `GET /physician/getAll` - Get all physicians
- `GET /physician/:id` - Get physician by ID

### Clinics
- `POST /clinic/register` - Register clinic
- `GET /clinic/:clinicId` - Get clinic by ID
- `GET /clinic/by-owner/:ownerId` - Get clinics by owner

### Medical Appointments
- `POST /appointments/create` - Create appointment
- `GET /appointments/getAll` - Get all appointments
- `PATCH /appointments/updateStatus/:id` - Update appointment status

### Medical Records
- `POST /medical-history/create` - Create medical record
- `GET /medical-history/patient/:patientId` - Get record by patient

## 🛡️ Security

- **Password encryption** with bcrypt
- **JWT tokens** for stateless authentication
- **Authentication middleware** on protected routes
- **Input data validation**
- **CORS** properly configured
- **Environment variables** for sensitive data
