# フロントエンドビルド
FROM node:20-alpine AS frontend
ARG VITE_ENABLE_RAG=false
ENV VITE_ENABLE_RAG=$VITE_ENABLE_RAG
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# バックエンドビルド
FROM golang:1.25-alpine AS backend
WORKDIR /app
COPY backend/ ./
RUN go build -o server ./cmd/server/main.go

# 本番イメージ
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=backend /app/server .
COPY --from=frontend /app/frontend/dist ./public
COPY backend/migrations ./migrations
COPY backend/seed ./seed
EXPOSE 8080
CMD ["./server"]
