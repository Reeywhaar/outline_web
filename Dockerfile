### build backend
FROM golang:1.22-bookworm AS go-builder
WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
COPY ./backend ./backend

# Download dependencies and build Go binary
RUN go build -o bin/main backend/main.go

### build frontend
FROM node:22-bookworm AS frontend-builder
WORKDIR /app

COPY package*.json ./
COPY ./front ./front
COPY ./esbuild.js ./esbuild.js
COPY ./tsconfig.json ./tsconfig.json

RUN npm install && npm run build:front

### Final image
FROM debian:bookworm-slim

WORKDIR /app

COPY ./templates ./templates
COPY --from=go-builder /app/bin /app/bin
COPY --from=frontend-builder /app/static /app/static

CMD ["./bin/main"]