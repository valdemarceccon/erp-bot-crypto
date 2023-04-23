#!/bin/sh

USERNAME=dev
DATABASE=dev
PASSWORD=dev_123456
ADM_USERNAME=valdemar.ceccon

docker-compose exec db psql -U "$USERNAME" -d "$DATABASE" -c
