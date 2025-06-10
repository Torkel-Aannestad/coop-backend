include .env

.PHONY: audit/api
audit/api:
	@echo 'Tidying and verifying module dependencies...'
	cd services/social-media-aggregator-api && go mod tidy
	cd services/social-media-aggregator-api && go mod verify
	@echo 'Formatting code...'
	cd services/social-media-aggregator-api && go fmt ./...
	@echo 'Vetting code...'
	cd services/social-media-aggregator-api && go vet ./...
	cd services/social-media-aggregator-api && staticcheck ./...
	@echo 'Running tests...'
	cd services/social-media-aggregator-api && go test -race -vet=off ./...



.PHONY: audit/mastodon
audit/mastodon:
	@echo 'Tidying and verifying module dependencies...'
	cd services/social-media-aggregator-mastodon && go mod tidy
	cd services/social-media-aggregator-mastodon && go mod verify
	@echo 'Formatting code...'
	cd services/social-media-aggregator-mastodon && go fmt ./...
	@echo 'Vetting code...'
	cd services/social-media-aggregator-mastodon && go vet ./...
	cd services/social-media-aggregator-mastodon && staticcheck ./...
	@echo 'Running tests...'
	cd services/social-media-aggregator-mastodon && go test -race -vet=off ./...

