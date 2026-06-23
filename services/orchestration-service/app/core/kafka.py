import json
from aiokafka import AIOKafkaProducer, AIOKafkaConsumer
from loguru import logger
from config.settings import settings


class KafkaService:
    def __init__(self):
        self.producer = None

    async def start(self):
        self.producer = AIOKafkaProducer(
            bootstrap_servers=settings.kafka_brokers,
            value_serializer=lambda v: json.dumps(v, default=str).encode(),
            acks="all",
        )
        await self.producer.start()
        logger.info("Orchestration Kafka producer started")

    async def stop(self):
        if self.producer:
            await self.producer.stop()

    async def publish(self, topic: str, key: str, value: dict):
        await self.producer.send(
            topic=topic, key=key.encode(), value=value
        )
        logger.debug(f"Published to {topic}: {key}")
