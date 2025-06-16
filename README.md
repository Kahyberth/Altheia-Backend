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

## 📁 Estructura del Proyecto

```
Altheia-Backend/
├── cmd/                    # Punto de entrada de la aplicación
│   └── main.go
├── internal/               # Código interno de la aplicación
│   ├── auth/              # Módulo de autenticación
│   ├── clinical/          # Módulo clínico
│   │   └── appointments/  # Gestión de citas
│   ├── users/             # Gestión de usuarios
│   │   ├── patient/       # Pacientes
│   │   ├── physician/     # Médicos
│   │   ├── receptionist/  # Recepcionistas
│   │   └── clinicOwner/   # Propietarios de clínicas
│   ├── middleware/        # Middlewares de la aplicación
│   ├── mail/              # Servicio de email
│   └── db/                # Configuración de base de datos
├── pkg/                   # Paquetes públicos
│   └── utils/             # Utilidades
├── config/                # Configuraciones
├── docker-compose.yaml    # Configuración de Docker
├── Makefile              # Comandos de automatización
├── go.mod                # Dependencias de Go
└── README.md             # Documentación
```

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

### 2. Configurar Variables de Entorno

Crea un archivo `.env` basado en las siguientes variables:

```env
# Base de datos
DB_HOST=localhost
DB_PORT=5432
DB_USER=tu_usuario
DB_PASSWORD=tu_contraseña
DB_NAME=altheia_db

# Servidor
PORT=8080
CLIENT=http://localhost:3000

# JWT
JWT_SECRET=tu_clave_secreta_muy_segura

# Email (opcional)
SMTP_HOST=tu_servidor_smtp
SMTP_PORT=587
SMTP_USER=tu_email
SMTP_PASSWORD=tu_contraseña_email
```

### 3. Instalación con Docker (Recomendado)

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

# Ejecutar
./bin/altheia
```

## 🧪 Testing

El proyecto incluye un conjunto completo de pruebas automatizadas:

```bash
# Ejecutar todas las pruebas
make test

# Pruebas con output detallado
make test-verbose

# Generar reporte de cobertura
make test-coverage

# Pruebas específicas por módulo
make test-auth
make test-clinical
make test-utils

# Pruebas de rendimiento
make benchmark
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

## 🔧 Comandos Útiles

```bash
# Desarrollo
make help          # Ver todos los comandos disponibles
make build         # Construir la aplicación
make clean         # Limpiar cache y archivos de build
make fmt           # Formatear código
make vet           # Análisis estático del código

# Testing
make test          # Ejecutar todas las pruebas
make test-race     # Pruebas con detección de race conditions
make test-short    # Pruebas rápidas

# Calidad de código
make check         # Ejecutar todas las verificaciones de calidad
```

## 📈 Rendimiento

- **Framework Fiber**: Basado en Fasthttp, uno de los más rápidos para Go
- **Conexiones de base de datos**: Pool de conexiones optimizado con GORM
- **Middlewares eficientes**: Procesamiento mínimo de overhead
- **Paginación**: Implementada para consultas de gran volumen

## 🤝 Contribución

1. Fork el proyecto
2. Crea tu rama de feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 👨‍💻 Desarrollador

Desarrollado con ❤️ por [Kahyberth](https://github.com/Kahyberth)

---

## 📞 Soporte

Si tienes preguntas o necesitas ayuda, por favor:
- Abre un [issue](https://github.com/Kahyberth/Altheia-Backend/issues)
- Contacta al desarrollador

---

*Altheia Backend - Transformando la gestión de historiales médicos digitales* 🏥✨

## Nuevo Endpoint: Historiales Clínicos de Clínica con Paginación

### GET `/medical-history/clinic/{clinicId}`

Este endpoint permite obtener todos los historiales clínicos de los pacientes asociados a una clínica específica, con paginación incluida.

#### Parámetros de URL:
- `clinicId` (string, requerido): ID único de la clínica

#### Parámetros de Query:
- `page` (int, opcional): Número de página (por defecto: 1)
- `size` (int, opcional): Tamaño de página (por defecto: 10, máximo: 100)

#### Ejemplo de Request:
```bash
GET /medical-history/clinic/clinic-123?page=1&size=10
```

#### Ejemplo de Response:
```json
{
  "success": true,
  "data": [
    {
      "patient": {
        "id": "patient-123",
        "name": "Juan Pérez",
        "age": 35,
        "gender": "Male",
        "dob": "1988-05-15",
        "mrn": "MRN-23456",
        "avatar": "/placeholder.svg?height=128&width=128&text=J"
      },
      "medicalRecords": [
        {
          "id": "REC-HISTORY-hist-123",
          "patientId": "patient-123",
          "type": "diagnoses",
          "title": "Historia Clínica - 2024-01-15",
          "date": "2024-01-15",
          "provider": "Sistema",
          "status": "active",
          "content": {
            "consult_reason": "Dolor de cabeza recurrente",
            "description": "Dolor de cabeza recurrente",
            "personal_info": "Paciente masculino de 35 años",
            "family_info": "Historia familiar de migraña",
            "allergies": "Penicilina",
            "observations": "Paciente presenta episodios frecuentes",
            "notes": "Historia clínica creada. Motivo: Dolor de cabeza recurrente. Observaciones: Paciente presenta episodios frecuentes"
          },
          "documents": []
        }
      ],
      "lastUpdate": "2024-01-15",
      "recordCount": 3
    }
  ],
  "pagination": {
    "currentPage": 1,
    "pageSize": 10,
    "totalPages": 1,
    "totalRecords": 5,
    "hasNext": false,
    "hasPrevious": false
  },
  "summary": {
    "totalPatients": 5,
    "totalMedicalRecords": 15,
    "recentActivity": "2024-01-15",
    "mostActivePatient": "Juan Pérez",
    "lastUpdated": "2024-01-15 14:30:25"
  }
}
```

#### Códigos de Estado:
- `200`: Éxito - Retorna los historiales clínicos paginados
- `400`: Error de validación - clinicId requerido o parámetros de paginación inválidos
- `500`: Error interno del servidor

#### Características:
- **Paginación**: Soporte completo para paginación con información detallada
- **Información del Paciente**: Datos básicos del paciente incluyendo edad calculada
- **Historiales Médicos**: Registros médicos completos con consultas y prescripciones
- **Resumen de Clínica**: Estadísticas generales de la clínica
- **Validación**: Validación de parámetros con valores por defecto seguros
- **Rendimiento**: Consultas optimizadas con preload selectivo

#### Notas Técnicas:
- Los parámetros de paginación tienen valores por defecto seguros
- La edad se calcula automáticamente a partir de la fecha de nacimiento
- Los MRN se generan automáticamente a partir del ID del paciente
- Los avatars se generan con la inicial del nombre del paciente
- El endpoint maneja correctamente pacientes sin historiales médicos

# Altheia Backend - API Documentation

## Descripción del Proyecto
Sistema backend para gestión de historiales clínicos con soporte completo para documentos e imágenes.

## 🏥 Sistema de Historiales Clínicos

### Arquitectura del Sistema
El sistema está diseñado con un enfoque escalable:
- **Una historia clínica base** por paciente (datos permanentes)
- **Múltiples consultas médicas** asociadas a cada historia
- **Documentos e imágenes** asociados tanto a historias como a consultas específicas

---

## 📄 **SISTEMA DE DOCUMENTOS Y ARCHIVOS**

### **Funcionalidades Implementadas:**

✅ **Subir documentos al crear historia clínica**  
✅ **Subir documentos al crear consulta médica**  
✅ **Agregar documentos a historia clínica existente**  
✅ **Agregar documentos a consulta existente**  
✅ **Consultar documentos por historia clínica**  
✅ **Consultar documentos por consulta específica**  
✅ **Soporte para Base64 y URLs**  
✅ **Validación de tipos de archivo**  
✅ **Metadata completa de documentos**

### **Tipos de archivos soportados:**
- **Imágenes**: JPG, JPEG, PNG, GIF
- **Documentos**: PDF, DOC, DOCX, TXT
- **Otros**: Cualquier tipo con MIME type personalizado

---

## 🔗 **ENDPOINTS DE DOCUMENTOS**

### 1. **Agregar documentos a historia clínica existente**
```http
POST /medical-history/documents/add
```

**Request Body:**
```json
{
  "medical_history_id": "hist-123",
  "uploaded_by": "physician-456",
  "documents": [
    {
      "name": "radiografia_torax.jpg",
      "type": "jpg",
      "description": "Radiografía de tórax del paciente",
      "base64_data": "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEAYABgAAD..."
    },
    {
      "name": "analisis_sangre.pdf", 
      "type": "pdf",
      "description": "Resultados de análisis de sangre",
      "url": "https://lab.ejemplo.com/resultados/123.pdf"
    }
  ]
}
```

### 2. **Agregar documentos a consulta específica**
```http
POST /medical-history/consultation/documents/add
```

**Request Body:**
```json
{
  "consultation_id": "cons-789",
  "uploaded_by": "physician-456", 
  "documents": [
    {
      "name": "prescripcion_medica.pdf",
      "type": "pdf", 
      "description": "Prescripción médica detallada",
      "base64_data": "data:application/pdf;base64,JVBERi0xLjQK..."
    }
  ]
}
```

### 3. **Consultar documentos por historia clínica**
```http
GET /medical-history/documents/{medicalHistoryId}
```

**Response:**
```json
{
  "success": true,
  "documents": [
    {
      "id": "doc-001",
      "name": "1634567890_radiografia_torax.jpg",
      "original_name": "radiografia_torax.jpg", 
      "type": "jpg",
      "size": 245760,
      "mime_type": "image/jpeg",
      "url": "/api/documents/xyz123_radiografia_torax.jpg",
      "description": "Radiografía de tórax del paciente",
      "uploaded_by": "physician-456",
      "uploaded_at": "2023-05-18T14:30:45",
      "is_public": false
    }
  ],
  "count": 1
}
```

### 4. **Consultar documentos por consulta**
```http
GET /medical-history/consultation/documents/{consultationId}
```

---

## 📤 **MÉTODOS DE SUBIDA DE ARCHIVOS**

### **Método 1: Base64 (Recomendado para archivos pequeños)**
```json
{
  "name": "documento.pdf",
  "type": "pdf",
  "description": "Descripción del documento",
  "base64_data": "data:application/pdf;base64,JVBERi0xLjQK..."
}
```

### **Método 2: URL externa**
```json
{
  "name": "resultado_laboratorio.pdf",
  "type": "pdf", 
  "description": "Resultados de laboratorio",
  "url": "https://laboratorio.com/resultados/paciente123.pdf"
}
```

---

## 🔄 **FLUJOS DE TRABAJO COMPLETOS**

### **Escenario 1: Crear historia clínica con documentos**
```json
POST /medical-history/create
{
  "patient_id": "patient-123",
  "physician_id": "physician-456",
  "consult_reason": "Consulta de rutina",
  "personal_info": "Paciente de 35 años...",
  "observations": "Sin antecedentes relevantes",
  "documents": [
    {
      "name": "historia_previa.pdf",
      "type": "pdf",
      "description": "Historia clínica previa",
      "base64_data": "data:application/pdf;base64,..."
    }
  ]
}
```

### **Escenario 2: Crear consulta con documentos**
```json
POST /medical-history/consultation/create
{
  "patient_id": "patient-123",
  "physician_id": "physician-456",
  "symptoms": "Dolor de cabeza frecuente",
  "diagnosis": "Migraña tensional",
  "treatment": "Ibuprofeno 400mg cada 8 horas",
  "documents": [
    {
      "name": "resonancia_magnetica.jpg",
      "type": "jpg",
      "description": "Resonancia magnética cerebral",
      "base64_data": "data:image/jpeg;base64,..."
    }
  ]
}
```

### **Escenario 3: Agregar documentos posteriormente**
```json
POST /medical-history/documents/add
{
  "medical_history_id": "hist-123",
  "uploaded_by": "physician-456",
  "documents": [
    {
      "name": "nueva_radiografia.jpg",
      "type": "jpg", 
      "description": "Radiografía de control",
      "base64_data": "data:image/jpeg;base64,..."
    }
  ]
}
```

---

## 💾 **ALMACENAMIENTO Y SEGURIDAD**

### **Estructura de almacenamiento:**
- **Base de datos**: Metadata de documentos
- **Sistema de archivos**: Archivos físicos en `/uploads/medical-documents/`
- **URLs**: Generación automática de URLs de acceso

### **Seguridad:**
- ✅ Nombres de archivo sanitizados
- ✅ Validación de tipos MIME
- ✅ Control de acceso (is_public flag)
- ✅ IDs únicos para prevenir colisiones
- ✅ Asociación segura a historias/consultas

### **Metadata completa:**
```json
{
  "id": "doc-unique-id",
  "medical_history_id": "hist-123",
  "consultation_id": "cons-456", 
  "name": "safe_filename",
  "original_name": "nombre_original.pdf",
  "type": "pdf",
  "size": 1024000,
  "mime_type": "application/pdf",
  "file_path": "/uploads/medical-documents/...",
  "url": "/api/documents/...",
  "description": "Descripción del documento",
  "uploaded_by": "user-id",
  "uploaded_at": "2023-05-18T14:30:45",
  "is_public": false
}
```

---

## 🎯 **CASOS DE USO PRÁCTICOS**

### **Para médicos:**
1. Subir prescripciones médicas en PDF
2. Adjuntar imágenes de radiografías
3. Documentar resultados de laboratorio
4. Agregar notas de seguimiento

### **Para pacientes:**
1. Subir documentos de historia médica previa
2. Adjuntar imágenes de síntomas
3. Compartir resultados de exámenes externos

### **Para laboratorios:**
1. Enviar resultados vía URL
2. Integración automática de reportes
3. Actualización de estados de análisis

---

## ⚡ **RENDIMIENTO Y LIMITACIONES**

### **Recomendaciones:**
- **Archivos pequeños (< 5MB)**: Usar Base64
- **Archivos grandes (> 5MB)**: Usar URLs externas
- **Máximo por request**: 10 documentos
- **Tipos recomendados**: PDF, JPG, PNG

### **Validaciones implementadas:**
- ✅ Verificación de existencia de historia/consulta
- ✅ Validación de tipos de archivo
- ✅ Sanitización de nombres
- ✅ Control de tamaños
- ✅ Transacciones atómicas

---

## 🔧 **CONFIGURACIÓN TÉCNICA**

### **Base de datos:**
Nueva tabla `medical_documents` con todas las relaciones necesarias.

### **Migraciones:**
Automáticas al ejecutar el servidor - tabla creada automáticamente.

### **Estructura de carpetas:**
```
/uploads/
  /medical-documents/
    /nanoid_timestamp_filename.ext
```

Este sistema proporciona una solución completa para gestión de documentos médicos con máxima flexibilidad y seguridad. 🏥📄✨