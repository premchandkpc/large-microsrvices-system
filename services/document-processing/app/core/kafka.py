import asyncio
import json
from typing import Callable, Optional
from aiokafka import AIOKafkaProducer, AIOKafkaConsumer
from loguru import logger
from config.settings import settings


class KafkaService:
    def __init__(self):
        self.producer: Optional[AIOKafkaProducer] = None
        self.consumer: Optional[AIOKafkaConsumer] = None

    async def start(self):
        self.producer = AIOKafkaProducer(
            bootstrap_servers=settings.kafka_brokers,
            value_serializer=lambda v: json.dumps(v, default=str).encode(),
            acks="all",
            compression_type="snappy",
        )
        await self.producer.start()
        logger.info("Kafka producer started")

    async def stop(self):
        if self.producer:
            await self.producer.stop()
        if self.consumer:
            await self.consumer.stop()
        logger.info("Kafka stopped")

    async def publish(self, topic: str, key: str, value: dict):
        await self.producer.send(
            topic=topic,
            key=key.encode(),
            value=value,
        )
        logger.debug(f"Published to {topic}: {key}")

    async def consume(self, topic: str, handler: Callable, group_id: Optional[str] = None):
        self.consumer = AIOKafkaConsumer(
            topic,
            bootstrap_servers=settings.kafka_brokers,
            group_id=group_id or settings.kafka_group_id,
            value_deserializer=lambda v: json.loads(v.decode()),
            auto_offset_reset="earliest",
            enable_auto_commit=True,
        )
        await self.consumer.start()
        logger.info(f"Consumer started for topic: {topic}")

        try:
            async for msg in self.consumer:
                await handler(msg)
        except asyncio.CancelledError:
            logger.info("Consumer cancelled")
        finally:
            await self.consumer.stop()
