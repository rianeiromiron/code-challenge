# ═══════════════════════════════════════════════════════════════════
# Etapa 1: Build (compilación)
# ═══════════════════════════════════════════════════════════════════
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Descargar dependencias primero (cache eficiente)
COPY go.mod ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilar el servidor estáticamente (sin dependencias de libc)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# ═══════════════════════════════════════════════════════════════════
# Etapa 2: Runtime (imagen final ligera)
# ═══════════════════════════════════════════════════════════════════
FROM alpine:latest

WORKDIR /app

# Copiar solo el binario compilado desde la etapa builder
COPY --from=builder /app/server .

# Puerto que expone la aplicación
EXPOSE 8080

# Comando por defecto al iniciar el contenedor
CMD ["./server"]