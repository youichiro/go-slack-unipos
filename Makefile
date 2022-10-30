psql:
	psql -h localhost -p 5432 -U postgres go_slack_unipos_development
sqlboiler:
	sqlboiler psql --config ./configs/sqlboiler.toml
migrate-up:
	docker compose run --rm migrate make up
migrate-down:
	docker compose run --rm migrate make down
migrate-force:
	docker compose run --rm migrate make force
