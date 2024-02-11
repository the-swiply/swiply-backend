from fastapi import FastAPI, APIRouter
from fastapi.middleware.cors import CORSMiddleware

from app.router import oracle
from app.snake import config

config.parse_yaml()

app = FastAPI(title=config.get("app").get("name"))

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

api_router = APIRouter()

app.include_router(oracle.router, prefix="/api")
