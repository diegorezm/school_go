# Instructions

## Setep - 1
create the .env file 
```bash
touch .env
echo "PASSWORD=password" >> .env
echo "DATABASE_NAME=db_name" >> .env
```

## Setep - 2
run docker compose
```bash
sudo docker-compose up -d
```

## Setep - 3
run the application
```bash
make run 
```
