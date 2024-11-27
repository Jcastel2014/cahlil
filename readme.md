sudo docker run --name postgresDatabase -e POSTGRES_USER=user -e POSTGRES_PASSWORD=user -e POSTGRES_DB=user -p 5432:5432 -d postgres
sudo docker start postgresDatabase