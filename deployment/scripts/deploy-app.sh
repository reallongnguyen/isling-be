source .env
sudo docker run -v ./database/postgres/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "$PG_URL?sslmode=disable" up
sudo docker compose build app
sudo docker compose -f ./deployment/docker-compose.yml -d app
