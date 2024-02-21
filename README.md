# Instructions

## Setep - 1
create the .env file 
```bash
touch .env
echo "PASSWORD=password\nDATABASE_NAME=db_name\nPORT=:3030" >> .env
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
