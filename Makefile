migrate-new:
	migrate create -ext sql -dir db/migration -seq $(name)

migrate-up:
	migrate -path db/migration \
	-database "postgresql://root:pass@127.0.0.1:5432/users?sslmode=disable&application_name=user-service" \
	-verbose up

migrate-down:
	migrate -path db/migration \
	-database "postgresql://root:pass@127.0.0.1:5432/users?sslmode=disable&application_name=user-service" \
	-verbose down

up:
	docker compose up -d