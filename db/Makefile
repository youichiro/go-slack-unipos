up:
	migrate -database ${POSTGRESQL_URL} -path /opt/db/migrations up
down:
	migrate -database ${POSTGRESQL_URL} -path /opt/db/migrations down
force:
	migrate -database ${POSTGRESQL_URL} -path /opt/db/migrations force 1
create:
	migrate create -ext sql -dir /opt/db/migrations -seq $(ARG)
