from typing import List, Dict, Any
from qdrant_client import QdrantClient
from qdrant_client.http import models
from elasticsearch import AsyncElasticsearch
from loguru import logger
from config.settings import settings


class IndexerService:
    def __init__(self):
        self.qdrant = None
        self.es = None

    async def initialize(self):
        self.qdrant = QdrantClient(
            host=settings.qdrant_host,
            port=settings.qdrant_port,
        )
        self.es = AsyncElasticsearch(settings.elasticsearch_url)
        await self._ensure_collections()
        logger.info("Indexer service initialized")

    async def _ensure_collections(self):
        collections = self.qdrant.get_collections().collections
        exists = any(c.name == settings.qdrant_collection for c in collections)
        if not exists:
            self.qdrant.create_collection(
                collection_name=settings.qdrant_collection,
                vectors_config=models.VectorParams(
                    size=1536,
                    distance=models.Distance.COSINE,
                ),
            )
            logger.info(f"Created Qdrant collection: {settings.qdrant_collection}")

        index_exists = await self.es.indices.exists(index=settings.elasticsearch_index)
        if not index_exists:
            await self.es.indices.create(
                index=settings.elasticsearch_index,
                body={
                    "settings": {
                        "number_of_shards": 2,
                        "number_of_replicas": 1,
                        "analysis": {
                            "analyzer": {
                                "default": {
                                    "type": "standard",
                                }
                            }
                        },
                    },
                    "mappings": {
                        "properties": {
                            "document_id": {"type": "keyword"},
                            "user_id": {"type": "keyword"},
                            "content": {"type": "text"},
                            "chunk_index": {"type": "integer"},
                            "created_at": {"type": "date"},
                        }
                    },
                },
            )
            logger.info(f"Created ES index: {settings.elasticsearch_index}")

    async def index_document(self, doc_id: str, chunks: List[Dict],
                             metadata: Dict[str, Any]):
        points = []
        for chunk in chunks:
            point_id = f"{doc_id}_{chunk['chunk_index']}"
            points.append(models.PointStruct(
                id=hash(point_id),
                vector=chunk["embedding"],
                payload={
                    "document_id": doc_id,
                    "chunk_index": chunk["chunk_index"],
                    "content": chunk["content"],
                    "token_count": chunk["token_count"],
                    **metadata,
                },
            ))

            await self.es.index(
                index=settings.elasticsearch_index,
                id=point_id,
                body={
                    "document_id": doc_id,
                    "user_id": metadata.get("user_id", ""),
                    "content": chunk["content"],
                    "chunk_index": chunk["chunk_index"],
                },
            )

        self.qdrant.upsert(
            collection_name=settings.qdrant_collection,
            points=points,
        )
        logger.info(f"Indexed {len(chunks)} chunks for document {doc_id}")

    async def search_vector(self, query_vector: List[float], top_k: int = 10,
                            score_threshold: float = 0.7) -> List[Dict]:
        results = self.qdrant.search(
            collection_name=settings.qdrant_collection,
            query_vector=query_vector,
            limit=top_k,
            score_threshold=score_threshold,
        )
        return [
            {
                "document_id": r.payload.get("document_id"),
                "content": r.payload.get("content"),
                "score": r.score,
                "chunk_index": r.payload.get("chunk_index"),
            }
            for r in results
        ]

    async def search_fulltext(self, query: str, page: int = 1,
                               page_size: int = 20) -> Dict:
        response = await self.es.search(
            index=settings.elasticsearch_index,
            body={
                "query": {
                    "multi_match": {
                        "query": query,
                        "fields": ["content^2", "document_id"],
                    }
                },
                "from": (page - 1) * page_size,
                "size": page_size,
                "highlight": {
                    "fields": {"content": {}},
                },
            },
        )

        hits = response["hits"]["hits"]
        return {
            "results": [
                {
                    "document_id": h["_source"]["document_id"],
                    "score": h["_score"],
                    "content": h["_source"]["content"],
                    "highlights": h.get("highlight", {}).get("content", []),
                }
                for h in hits
            ],
            "total": response["hits"]["total"]["value"],
            "took_ms": response["took"],
        }
