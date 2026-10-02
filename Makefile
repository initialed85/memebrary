.PHONY: dev-backend dev-frontend test build

dev-backend:
	cd backend && DATA_DIR=./data AI_BASE_URL=$${AI_BASE_URL:-http://192.168.137.111:8088/v1} go run ./cmd/memebrary

dev-frontend:
	cd frontend && npm run dev

test:
	cd backend && go test ./...
	cd frontend && npm run check

build:
	cd backend && CGO_ENABLED=0 go build -o ./memebrary ./cmd/memebrary
	cd frontend && npm run build
