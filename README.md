# 🏥 Altheia Backend

**Altheia** es un sistema backend completo para una aplicación de Historia Clínica Electrónica (EHR) desarrollado en Go. Este sistema proporciona una API REST segura, escalable y moderna diseñada para facilitar la gestión integral de servicios de salud.

[![Go Version](https://img.shields.io/badge/Go-1.24-blue)](https://golang.org/)
[![Fiber Framework](https://img.shields.io/badge/Fiber-v2.52.6-00ADD8)](https://github.com/gofiber/fiber)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

## 📋 Descripción

Altheia Backend es una solución integral que permite la gestión completa de servicios médicos, incluyendo:
- Gestión de pacientes, médicos, recepcionistas y propietarios de clínicas
- Historiales médicos digitales y consultas
- Sistema de citas médicas
- Administración de clínicas y servicios
- Autenticación y autorización robusta con JWT

## ✨ Características Principales

### 🔐 Autenticación y Seguridad
- Autenticación JWT con tokens de acceso y refresh
- Middleware de protección de rutas
- Gestión de actividades de login
- Encriptación de contraseñas
- Control de acceso basado en roles

### 👥 Gestión de Usuarios
- **Pacientes**: Registro, actualización, eliminación lógica
- **Médicos**: Gestión completa con paginación
- **Recepcionistas**: Administración del personal de recepción
- **Propietarios de clínicas**: Gestión de dueños de centros médicos
- **Técnicos de laboratorio**: Soporte para personal técnico

### 🏥 Gestión Clínica
- Registro y administración de clínicas
- Información detallada de cada centro médico
- Horarios y servicios ofrecidos
- Asignación de personal a clínicas
- Integración con EPS (Entidades Promotoras de Salud)

### 📋 Historiales Médicos
- Creación y actualización de historiales
- Consultas médicas detalladas
- Prescripciones médicas
- Seguimiento de tratamientos

### 📅 Sistema de Citas
- Programación de citas médicas
- Gestión de estados de citas
- Filtrado por médico y fecha
- Notificaciones y recordatorios

## 🛠️ Tecnologías

- **[Go](https://golang.org/)** 1.24 - Lenguaje de programación
- **[Fiber](https://github.com/gofiber/fiber)** v2 - Framework web ultrarrápido
- **[GORM](https://gorm.io/)** - ORM para Go
- **[PostgreSQL](https://www.postgresql.org/)** 16 - Base de datos relacional
- **[JWT](https://github.com/golang-jwt/jwt)** - Autenticación con tokens
- **[Docker](https://www.docker.com/)** - Containerización
- **[bcrypt](https://golang.org/x/crypto)** - Encriptación de contraseñas



## 🚀 Instalación y Configuración

### Prerrequisitos

- Go 1.24 o superior
- PostgreSQL 16
- Docker y Docker Compose (opcional)

### 1. Clonar el Repositorio

```bash
git clone https://github.com/Kahyberth/Altheia-Backend.git
cd Altheia-Backend
```


```
### 2. Instalación con Docker (Recomendado)

```bash
# Iniciar base de datos PostgreSQL
docker-compose up -d

# Instalar dependencias
make deps

# Ejecutar la aplicación
make build
./bin/altheia
```

### 4. Instalación Manual

```bash
# Instalar dependencias
go mod download

# Construir la aplicación
go build -o bin/altheia ./cmd/main.go


```

## 📖 API Endpoints

### Autenticación
- `POST /auth/login` - Iniciar sesión
- `POST /auth/logout` - Cerrar sesión
- `GET /auth/verify-token` - Verificar token
- `POST /auth/refresh-token/:refresh_token` - Renovar token

### Pacientes
- `POST /patient/register` - Registrar paciente
- `GET /patient/getAll` - Obtener todos los pacientes
- `GET /patient/getAllPaginated` - Obtener pacientes paginados
- `PATCH /patient/update/:id` - Actualizar paciente

### Médicos
- `POST /physician/register` - Registrar médico
- `GET /physician/getAll` - Obtener todos los médicos
- `GET /physician/:id` - Obtener médico por ID

### Clínicas
- `POST /clinic/register` - Registrar clínica
- `GET /clinic/:clinicId` - Obtener clínica por ID
- `GET /clinic/by-owner/:ownerId` - Obtener clínicas por propietario

### Citas Médicas
- `POST /appointments/create` - Crear cita
- `GET /appointments/getAll` - Obtener todas las citas
- `PATCH /appointments/updateStatus/:id` - Actualizar estado de cita

### Historiales Médicos
- `POST /medical-history/create` - Crear historial médico
- `GET /medical-history/patient/:patientId` - Obtener historial por paciente

## 🛡️ Seguridad

- **Encriptación de contraseñas** con bcrypt
- **Tokens JWT** para autenticación stateless
- **Middleware de autenticación** en rutas protegidas
- **Validación de datos** de entrada
- **CORS** configurado adecuadamente
- **Variables de entorno** para datos sensibles
