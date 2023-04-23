.PHONY: up down new auto

EXEC_BACK=docker-compose exec backend

up:
	$(EXEC_BACK) alembic upgrade head
down:
	$(EXEC_BACK) alembic downgrade base
new:
	$(EXEC_BACK) alembic revision -m "$(MESSAGE)"
auto:
	$(EXEC_BACK) alembic revision --autogenerate -m "$(MESSAGE)"
