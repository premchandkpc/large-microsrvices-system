from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession, async_sessionmaker
from sqlalchemy.orm import DeclarativeBase
from loguru import logger
from config.settings import settings


class Base(DeclarativeBase):
    pass


class DatabaseService:
    def __init__(self):
        self.engine = None
        self.session_factory = None

    async def initialize(self):
        self.engine = create_async_engine(
            settings.postgres_url,
            pool_size=10,
            max_overflow=20,
            pool_pre_ping=True,
            echo=False,
        )
        self.session_factory = async_sessionmaker(
            self.engine, class_=AsyncSession, expire_on_commit=False
        )
        async with self.engine.begin() as conn:
            await conn.run_sync(Base.metadata.create_all)
        logger.info("Database initialized")

    async def get_session(self) -> AsyncSession:
        return self.session_factory()

    async def close(self):
        if self.engine:
            await self.engine.dispose()
