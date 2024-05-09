docker run -p 5432:5432 --name postgres -e POSTGRES_PASSWORD=admin -d postgres

docker run -p 54321:80 --name pgadmin -e PGADMIN_DEFAULT_EMAIL=admin@molta9a.com -e PGADMIN_DEFAULT_PASSWORD=admin -d dpage/pgadmin4
