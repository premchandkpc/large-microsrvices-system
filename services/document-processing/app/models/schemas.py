from pydantic import BaseModel
from typing import Optional, List, Dict
from datetime import datetime


class HealthResponse(BaseModel):
    status: str
    service: str
    version: str


class ProcessingRequest(BaseModel):
    document_id: str
    pipeline_type: str = "default"
    options: Dict[str, str] = {}


class ProcessingResponse(BaseModel):
    job_id: str
    status: str
    document_id: str


class ChunkResult(BaseModel):
    chunk_index: int
    content: str
    embedding: List[float]
    token_count: int


class ProcessingResult(BaseModel):
    document_id: str
    status: str
    chunks: List[ChunkResult] = []
    summary: Optional[str] = None
    extracted_metadata: Dict[str, str] = {}
    processing_time_ms: float
    error: Optional[str] = None
