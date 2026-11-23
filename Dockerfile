# ==========================================
# Stage 1: Build Frontend (Svelte SPA)
# ==========================================
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# ==========================================
# Stage 2: Build Backend (Golang Binary)
# ==========================================
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

# ==========================================
# Stage 3: Minimal Production Runtime
# ==========================================
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Salin binary backend
COPY --from=backend-builder /app/backend/server /app/server

# Salin aset uploads default (termasuk diagram SVG)
COPY backend/uploads /app/uploads

# Salin hasil build static frontend Svelte ke /app/static
COPY --from=frontend-builder /app/frontend/dist /app/static

EXPOSE 8080

CMD ["./server"]
