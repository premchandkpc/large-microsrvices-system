import io
from typing import Optional
import boto3
from botocore.config import Config
from loguru import logger
from config.settings import settings


class StorageService:
    def __init__(self):
        self.client = None
        self.bucket = settings.s3_bucket

    async def initialize(self):
        self.client = boto3.client(
            "s3",
            endpoint_url=settings.s3_endpoint,
            aws_access_key_id=settings.s3_access_key,
            aws_secret_access_key=settings.s3_secret_key,
            region_name=settings.s3_region,
            config=Config(
                connect_timeout=5,
                read_timeout=30,
                retries={"max_attempts": 3},
            ),
        )
        logger.info(f"Storage service initialized (endpoint={settings.s3_endpoint})")

    async def download(self, key: str) -> Optional[bytes]:
        try:
            response = self.client.get_object(Bucket=self.bucket, Key=key)
            return response["Body"].read()
        except Exception as e:
            logger.error(f"S3 download failed for {key}: {e}")
            return None

    async def get_text_content(self, key: str) -> str:
        data = await self.download(key)
        if data is None:
            return ""
        try:
            return data.decode("utf-8")
        except UnicodeDecodeError:
            return data.decode("latin-1")
