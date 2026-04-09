#!/usr/bin/env bash
set -e

# フロントエンドをビルド
cd frontend
npm install
npm run build
cd ..

# ビルド成果物をbackend/publicに配置
rm -rf backend/public
cp -r frontend/dist backend/public

# バックエンドをビルド
cd backend
go build -o server ./cmd/server/main.go
