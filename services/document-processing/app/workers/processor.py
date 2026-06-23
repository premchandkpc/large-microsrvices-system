import asyncio
import json
import time
from typing import Optional
from aiokafka import AIOKafkaConsumer
from loguru import logger

from app.core.kafka import KafkaService
from app.core.embeddings import EmbeddingService
from app.core.storage import StorageService
from app.services.llm import LLMService
from app.services.indexer import IndexerService
from config.settings import settings


class DocumentProcessor:
    def __init__(self, kafka: KafkaService, embeddings: EmbeddingService,
                 storage: StorageService):
        self.kafka = kafka
        self.embeddings = embeddings
        self.storage = storage
        self.llm_service = LLMService()
        self.indexer = IndexerService()

    async def process(self, message):
        start_time = time.time()
        try:
            data = message.value if hasattr(message, 'value') else message
            doc_id = data.get("document_id", "")
            pipeline = data.get("pipeline_type", "default")
            user_id = data.get("user_id", "")
            options = data.get("options", {})

            logger.info(f"Processing document {doc_id} with pipeline {pipeline}")

            s3_key = f"{user_id}/documents/{doc_id}"
            content = await self.storage.get_text_content(s3_key)
            if not content:
                raise ValueError(f"Empty content for {s3_key}")

            chunks = self.embeddings.chunk_text(content)
            logger.info(f"Document {doc_id}: {len(chunks)} chunks created")

            all_embeddings = []
            for i, chunk in enumerate(chunks):
                embedding = await self.embeddings.embed_text(chunk)
                all_embeddings.append({
                    "chunk_index": i,
                    "content": chunk,
                    "embedding": embedding,
                    "token_count": self.embeddings.count_tokens(chunk),
                })

            await self.indexer.index_document(
                doc_id=doc_id,
                chunks=all_embeddings,
                metadata={
                    "user_id": user_id,
                    "pipeline": pipeline,
                    "num_chunks": len(chunks),
                    "content_type": data.get("content_type", "text/plain"),
                },
            )

            summary = None
            if pipeline in ("summarize", "full"):
                summary = await self.llm_service.summarize(
                    content[:settings.max_chunk_size * 3]
                )

            processing_time = (time.time() - start_time) * 1000

            result = {
                "document_id": doc_id,
                "status": "completed",
                "num_chunks": len(chunks),
                "summary": summary,
                "processing_time_ms": processing_time,
                "user_id": user_id,
            }
            await self.kafka.publish(
                settings.kafka_output_topic, doc_id, result
            )

            analytics = {
                "documentId": doc_id,
                "userId": user_id,
                "action": "process",
                "fileType": data.get("content_type", "text/plain"),
                "processingTimeMs": processing_time,
                "success": True,
            }
            await self.kafka.publish(
                settings.kafka_analytics_topic, doc_id, analytics
            )

            logger.info(f"Document {doc_id} processed in {processing_time:.0f}ms")

        except Exception as e:
            processing_time = (time.time() - start_time) * 1000
            logger.error(f"Processing failed: {e}")

            error_result = {
                "document_id": doc_id if 'doc_id' in dir() else "unknown",
                "status": "failed",
                "error": str(e),
                "processing_time_ms": processing_time,
            }
            await self.kafka.publish(
                settings.kafka_output_topic, error_result.get("document_id"), error_result
            )


async def start_consumer(kafka: KafkaService, embeddings: EmbeddingService,
                         storage: StorageService):
    processor = DocumentProcessor(kafka, embeddings, storage)
    task = asyncio.create_task(
        kafka.consume(settings.kafka_input_topic, processor.process)
    )
    logger.info("Document processing consumer started")
    return task
