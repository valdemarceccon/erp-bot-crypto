.PHONY: up down new auto

EXEC_BACK=docker-compose exec backend
EXEC_DB=docker-compose exec db

up:
	$(EXEC_BACK) alembic upgrade head
down:
	$(EXEC_BACK) alembic downgrade base
new:
	$(EXEC_BACK) alembic revision -m "$(MESSAGE)"
auto:
	$(EXEC_BACK) alembic revision --autogenerate -m "$(MESSAGE)"
resetdb:
	$(EXEC_DB) psql -U dev -d postgres -c "drop database dev;"
	$(EXEC_DB) psql -U dev -d postgres -c "create database dev;"
