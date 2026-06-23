from typing import Optional
from pydantic import BaseSettings


class KafkaConfig(BaseSettings):
    brokers: str = "localhost:9092"
    input_topic: str = "document-processing"
    output_topic: str = "document-processed"
    analytics_topic: str = "document-analytics"
    group_id: str = "document-processing"

    class Config:
        env_prefix = "KAFKA_"
