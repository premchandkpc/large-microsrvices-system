from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager
from loguru import logger

from app.api.routes import router
from app.core.kafka import KafkaService
from app.core.embeddings import EmbeddingService
from app.core.storage import StorageService
from app.models.schemas import HealthResponse
from config.settings import settings


@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info(f"Starting {settings.service_name} in {settings.environment} mode")
    kafka_service = KafkaService()
    embedding_service = EmbeddingService()
    storage_service = StorageService()
    app.state.kafka = kafka_service
    app.state.embeddings = embedding_service
    app.state.storage = storage_service
    await kafka_service.start()
    await embedding_service.initialize()
    from app.workers.processor import start_consumer
    app.state.consumer_task = await start_consumer(kafka_service, embedding_service, storage_service)
    yield
    await kafka_service.stop()
    if hasattr(app.state, 'consumer_task'):
        app.state.consumer_task.cancel()


app = FastAPI(
    title="Document Processing Service",
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

app.include_router(router, prefix="/api/v1")


@app.get("/health")
async def health():
    return HealthResponse(
        status="UP",
        service=settings.service_name,
        version="1.0.0",
    ).model_dump()
