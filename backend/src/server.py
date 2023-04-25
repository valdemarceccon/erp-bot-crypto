import os

from fastapi import FastAPI
from fastapi import HTTPException
from fastapi.security import OAuth2PasswordBearer
from sqlalchemy import text
from src.dependencies.database import engine
from src.routers import auth
from src.routers import user

root_path = os.getenv("BASE_PATH", "")

app = FastAPI(root_path=root_path)
app.include_router(user.router)
app.include_router(auth.router)


@app.get("/")
def read_root():
    return {"Hello": "World"}


# EXTERNAL_API_URL = "https://example.com/api/health"


@app.get("/health")
async def health_check():
    health_status = {
        "database": "unknown",
        # "external_api": "unknown",
    }

    # Check the database connection
    try:
        with engine.connect() as connection:
            connection.execute(text("SELECT 1"))
        health_status["database"] = "OK"
    except Exception as e:
        health_status["database"] = "unavailable"

    # Check the external API
    # try:
    #     async with httpx.AsyncClient() as client:
    #         response = await client.get(EXTERNAL_API_URL)
    #         response.raise_for_status()
    #     health_status["external_api"] = "OK"
    # except Exception as e:
    #     health_status["external_api"] = "unavailable"

    if all(value == "OK" for value in health_status.values()):
        return health_status
    else:
        raise HTTPException(status_code=500, detail=health_status)


# OAuth2 setup
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
