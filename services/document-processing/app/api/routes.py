from fastapi import APIRouter, HTTPException
from app.models.schemas import ProcessingRequest, ProcessingResponse
from config.settings import settings

router = APIRouter()


@router.post("/process")
async def process_document(request: ProcessingRequest):
    from app.core.kafka import KafkaService

    kafka = KafkaService()
    await kafka.start()

    await kafka.publish(
        settings.kafka_input_topic,
        request.document_id,
        request.dict(),
    )
    await kafka.stop()

    return ProcessingResponse(
        job_id=request.options.get("job_id", request.document_id),
        status="queued",
        document_id=request.document_id,
    )


@router.get("/pipelines")
async def list_pipelines():
    return {
        "pipelines": [
            {"name": "default", "description": "Chunk + embed + index"},
            {"name": "summarize", "description": "Chunk + embed + index + summarize"},
            {"name": "full", "description": "Full processing with entity extraction"},
        ]
    }
