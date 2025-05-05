<div align="center">
  <img src="ruta/a/tu/imagen.png" alt="Altheia Logo" width="300"/>
</div>

# Altheia - Backend

**Altheia** es el sistema backend de una aplicación de Historia Clínica Electrónica (EHR, por sus siglas en inglés) diseñada para mejorar la gestión, almacenamiento y acceso a los datos médicos de los pacientes. Este sistema proporciona una API segura, escalable y moderna para facilitar la integración con distintos servicios de salud.

---

## 🧠 Características Principales

- Gestión de pacientes, médicos, historiales médicos, recetas, diagnósticos, y más.
- Autenticación y autorización con JWT.
- API RESTful bien estructurada.
- Base de datos relacional optimizada para entidades médicas.
- Buenas prácticas de desarrollo y seguridad.

---

## 📦 Tecnologías Utilizadas

- **Go** (o el lenguaje que estés usando)
- **Fiber** / **Gin** / **Echo** (según el framework)
- **GORM** como ORM
- **PostgreSQL**
- **Docker**
- **JWT** para autenticación
- **Swagger/OpenAPI** para documentación de endpoints

---

## 🚀 Instalación y Ejecución

```bash
# Clona el repositorio
git clone git@github.com:Kahyberth/Altheia-Backend.git
cd altheia-backend

# Configura las variables de entorno (ver .env.example)
cp .env.example .env

# Corre el docker-compose

docker-compose up -d

# Ejecuta la aplicación
go run main.go
