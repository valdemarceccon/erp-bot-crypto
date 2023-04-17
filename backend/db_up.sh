#!/bin/sh
export PATH="$(pwd)/venv/bin:$PATH"
docker-compose exec backend alembic upgrade head