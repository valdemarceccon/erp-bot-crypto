#!/bin/sh

docker-compose exec backend alembic revision --autogenerate -m "$*"