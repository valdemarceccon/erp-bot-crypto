# backend/Dockerfile
FROM python:3.10

WORKDIR /app

COPY . .

RUN python setup.py install

CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8000"]

