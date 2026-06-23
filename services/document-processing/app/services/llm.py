from typing import Optional, List
from openai import AsyncOpenAI
from loguru import logger
from tenacity import retry, stop_after_attempt, wait_exponential
from config.settings import settings


class LLMService:
    def __init__(self):
        self.client = None
        if settings.openai_api_key:
            self.client = AsyncOpenAI(api_key=settings.openai_api_key)

    @retry(
        stop=stop_after_attempt(3),
        wait=wait_exponential(multiplier=1, min=2, max=10),
    )
    async def summarize(self, text: str, max_length: int = 500) -> Optional[str]:
        if not self.client:
            return "LLM not configured"
        try:
            response = await self.client.chat.completions.create(
                model=settings.openai_model,
                messages=[
                    {
                        "role": "system",
                        "content": "Summarize the following document concisely. "
                                   "Focus on key points, entities, and main topics.",
                    },
                    {"role": "user", "content": text[:8000]},
                ],
                max_tokens=max_length,
                temperature=0.3,
            )
            return response.choices[0].message.content
        except Exception as e:
            logger.error(f"LLM summarization failed: {e}")
            raise

    @retry(stop=stop_after_attempt(3), wait=wait_exponential(multiplier=1, min=2, max=10))
    async def extract_entities(self, text: str) -> dict:
        if not self.client:
            return {"entities": []}
        try:
            response = await self.client.chat.completions.create(
                model=settings.openai_model,
                messages=[
                    {
                        "role": "system",
                        "content": "Extract named entities (people, organizations, "
                                   "locations, dates) from the text as JSON.",
                    },
                    {"role": "user", "content": text[:4000]},
                ],
                response_format={"type": "json_object"},
                temperature=0.1,
            )
            return response.choices[0].message.content
        except Exception as e:
            logger.error(f"Entity extraction failed: {e}")
            return {"entities": []}

    async def classify(self, text: str, categories: List[str]) -> dict:
        if not self.client:
            return {"category": "unknown", "confidence": 0.0}
        try:
            response = await self.client.chat.completions.create(
                model=settings.openai_model,
                messages=[
                    {
                        "role": "system",
                        "content": f"Classify the text into one of: {', '.join(categories)}. "
                                   "Respond with JSON: {\"category\": \"...\", \"confidence\": 0.0}",
                    },
                    {"role": "user", "content": text[:2000]},
                ],
                response_format={"type": "json_object"},
                temperature=0.1,
            )
            return response.choices[0].message.content
        except Exception as e:
            logger.error(f"Classification failed: {e}")
            return {"category": "unknown", "confidence": 0.0}
