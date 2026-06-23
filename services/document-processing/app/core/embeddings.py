import tiktoken
from typing import List, Optional
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_openai import OpenAIEmbeddings
from loguru import logger
from config.settings import settings


class EmbeddingService:
    def __init__(self):
        self.encoder = None
        self.embeddings = None
        self.text_splitter = None

    async def initialize(self):
        self.encoder = tiktoken.get_encoding("cl100k_base")
        self.text_splitter = RecursiveCharacterTextSplitter(
            chunk_size=settings.max_chunk_size,
            chunk_overlap=settings.chunk_overlap,
            length_function=len,
            separators=["\n\n", "\n", ".", " ", ""],
        )
        if settings.openai_api_key:
            self.embeddings = OpenAIEmbeddings(
                model=settings.embedding_model,
                openai_api_key=settings.openai_api_key,
            )
        logger.info(f"Embedding service initialized (ollama={settings.use_ollama})")

    def chunk_text(self, text: str) -> List[str]:
        return self.text_splitter.split_text(text)

    async def embed_text(self, text: str) -> List[float]:
        if not self.embeddings:
            return [0.0] * 1536
        try:
            embedding = await self.embeddings.aembed_query(text)
            return embedding
        except Exception as e:
            logger.error(f"Embedding failed: {e}")
            return [0.0] * 1536

    async def embed_batch(self, texts: List[str]) -> List[List[float]]:
        results = []
        for text in texts:
            embedding = await self.embed_text(text)
            results.append(embedding)
        return results

    def count_tokens(self, text: str) -> int:
        return len(self.encoder.encode(text))
