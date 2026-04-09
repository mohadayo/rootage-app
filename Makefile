.PHONY: db-up db-down migrate seed run-backend run-frontend setup

include .env
export

db-up:
	docker compose up -d

db-down:
	docker compose down

migrate:
	@for f in $$(ls backend/migrations/*.sql | sort); do \
		echo "Applying $$f..."; \
		psql "$(DB_URL)" -f "$$f"; \
	done

seed:
	psql "$(DB_URL)" -f backend/seed/seed.sql

run-backend:
	cd backend && go run ./cmd/server

run-frontend:
	cd frontend && npm run dev

setup: db-up
	@echo "Waiting for PostgreSQL..."
	@sleep 3
	$(MAKE) migrate
	$(MAKE) seed
	@echo "Setup complete!"
