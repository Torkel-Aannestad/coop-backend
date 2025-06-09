include .env

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	@cd sql/migrations && goose postgres ${DSN_DB_DEV_LOCAL_MACHINE} up && cd ../..

