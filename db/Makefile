up:
	migrate -database ${POSTGRESQL_URL} -path /opt/db/migrations up
down:
	migrate -database ${POSTGRESQL_URL} -path /opt/db/migrations down
create:
	migrate create -ext sql -dir /opt/db/migrations -seq $(ARG)
