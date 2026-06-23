from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager
from loguru import logger

from app.api.routes import router
from app.core.kafka import KafkaService
from app.core.db import DatabaseService
from app.workers.workflow_engine import WorkflowEngine
from config.settings import settings


@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info(f"Starting {settings.service_name}")
    db = DatabaseService()
    await db.initialize()
    kafka = KafkaService()
    await kafka.start()
    workflow_engine = WorkflowEngine(kafka, db)
    await workflow_engine.start()

    app.state.db = db
    app.state.kafka = kafka
    app.state.workflow_engine = workflow_engine

    yield

    await workflow_engine.stop()
    await kafka.stop()
    logger.info("Orchestration service stopped")


app = FastAPI(
    title="Orchestration Service",
    version="1.0.0",
    lifespan=lifespan,
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(router, prefix="/api/v1/workflows")


@app.get("/health")
async def health():
    return {"status": "UP", "service": settings.service_name}
