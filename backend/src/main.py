from fastapi import FastAPI, Depends, HTTPException
from pydantic import BaseModel
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from typing import List
from datetime import datetime
import bcrypt

app = FastAPI()

class User(BaseModel):
    user_id: str
    oauth2_hash: str
    jwt_token: str
    created_at: datetime

class BotStatus(BaseModel):
    user_id: str
    bot_status: str

class PaymentProof(BaseModel):
    payment_id: str
    user_id: str
    transaction_hash: str
    created_at: datetime

users = []
bot_statuses = []
payment_proofs = []

def get_current_user(token: str):
    for user in users:
        if user.jwt_token == token:
            return user
    raise HTTPException(status_code=404, detail="User not found")

app = FastAPI()

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

users = []
bot_statuses = []
payment_proofs = []

app = FastAPI()

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")

users = []
bot_statuses = []
payment_proofs = []

def hash_password(password: str):
    return bcrypt.hashpw(password.encode(), bcrypt.gensalt()).decode()

def verify_password(password: str, hashed_password: str):
    return bcrypt.checkpw(password.encode(), hashed_password.encode())

def authenticate_user(email: str, password: str):
    for user in users:
        if user.email == email and verify_password(password, user.hashed_password):
            return user
    return None

def get_current_user(token: str = Depends(oauth2_scheme)):
    for user in users:
        if user.jwt_token == token:
            return user
    raise HTTPException(status_code=404, detail="User not found")

def authenticate_user(email: str, password: str):
    for user in users:
        if user.email == email and user.oauth2_hash == password:
            return user
    return None

def get_current_user(token: str = Depends(oauth2_scheme)):
    for user in users:
        if user.jwt_token == token:
            return user
    raise HTTPException(status_code=404, detail="User not found")

def hash_password(password: str):
    return bcrypt.hashpw(password.encode(), bcrypt.gensalt()).decode()

def verify_password(password: str, hashed_password: str):
    return bcrypt.checkpw(password.encode(), hashed_password.encode())

def authenticate_user(email: str, password: str):
    for user in users:
        if user.email == email and verify_password(password, user.hashed_password):
            return user
    return None

def get_current_user(token: str = Depends(oauth2_scheme)):
    for user in users:
        if user.jwt_token == token:
            return user
    raise HTTPException(status_code=404, detail="User not found")

@app.post("/v1/users")
def create_user(email: str, password: str, jwt_token: str):
    hashed_password = hash_password(password)
    user = User(email=email, hashed_password=hashed_password, jwt_token=jwt_token)
    users.append(user)
    return {"email": user.email, "message": "User created successfully"}

@app.post("/token")
async def login(form_data: OAuth2PasswordRequestForm = Depends()):
    user = authenticate_user(form_data.username, form_data.password)
    if not user:
        raise HTTPException(status_code=400, detail="Invalid email or password")
    return {"access_token": user.jwt_token, "token_type": "bearer"}

@app.get("/v1/users/me", response_model=User)
def get_self_info(current_user: User = Depends(get_current_user)):
    return current_user

@app.get("/v1/admin/users", response_model=List[User])
def get_user_list():
    return users

@app.get("/v1/admin/users/{user_id}", response_model=User)
def get_user_details(user_id: str):
    for user in users:
        if user.user_id == user_id:
            return user
    raise HTTPException(status_code=404, detail="User not found")

@app.get("/v1/users/me/bot/status", response_model=BotStatus)
def get_bot_status(current_user: User = Depends(get_current_user)):
    for status in bot_statuses:
        if status.user_id == current_user.user_id:
            return status
    raise HTTPException(status_code=404, detail="Bot status not found")

@app.put("/v1/users/me/bot/status")
def set_bot_status(new_status: str, current_user: User = Depends(get_current_user)):
    for status in bot_statuses:
        if status.user_id == current_user.user_id:
            status.bot_status = new_status
            return {"user_id": current_user.user_id, "bot_status": new_status, "message": "Bot status updated successfully"}

    new_status = BotStatus(user_id=current_user.user_id, bot_status=new_status)
    bot_statuses.append(new_status)
    return {"user_id": current_user.user_id, "bot_status": new_status.bot_status, "message": "Bot status updated successfully"}

@app.post("/v1/users/me/payments")
def add_payment_proof(transaction_hash: str, current_user: User = Depends(get_current_user)):
    payment_proof = PaymentProof(payment_id=str(len(payment_proofs) + 1), user_id=current_user.user_id, transaction_hash=transaction_hash, created_at=datetime.utcnow())
    payment_proofs.append(payment_proof)
    return {"payment_id": payment_proof.payment_id, "user_id": current_user.user_id, "transaction_hash": transaction_hash, "created_at": payment_proof.created_at, "message": "Payment proof submitted successfully"}

@app.get("/")
def get_user_details():
    return {"message": "hello"}
