import asyncio
from typing import Dict, Any, Callable, List
from datetime import datetime, timezone
from loguru import logger

from app.core.kafka import KafkaService
from app.core.db import DatabaseService
from app.models.workflow import Workflow, WorkflowStatus
from config.settings import settings


class WorkflowEngine:
    def __init__(self, kafka: KafkaService, db: DatabaseService):
        self.kafka = kafka
        self.db = db
        self._handlers: Dict[str, Callable] = {}
        self._running = False

    async def start(self):
        self._running = True
        self._register_default_handlers()
        logger.info("Workflow engine started")

    async def stop(self):
        self._running = False
        logger.info("Workflow engine stopped")

    def _register_default_handlers(self):
        self.register_handler("document_upload", self._handle_document_upload)
        self.register_handler("document_process", self._handle_document_process)
        self.register_handler("notify_user", self._handle_notify_user)
        self.register_handler("index_document", self._handle_index_document)
        self.register_handler("generate_report", self._handle_generate_report)

    def register_handler(self, step_type: str, handler: Callable):
        self._handlers[step_type] = handler

    async def create_workflow(self, name: str, steps: List[Dict],
                               user_id: str = None,
                               document_id: str = None) -> Workflow:
        async with self.db.get_session() as session:
            workflow = Workflow(
                name=name,
                user_id=user_id,
                document_id=document_id,
                steps=steps,
                context={"user_id": user_id, "document_id": document_id},
            )
            session.add(workflow)
            await session.commit()
            await session.refresh(workflow)
            logger.info(f"Workflow created: {workflow.id} - {name}")
            return workflow

    async def execute_workflow(self, workflow_id: str):
        async with self.db.get_session() as session:
            workflow = await session.get(Workflow, workflow_id)
            if not workflow:
                raise ValueError(f"Workflow {workflow_id} not found")

            workflow.status = WorkflowStatus.RUNNING
            await session.commit()

        try:
            async with self.db.get_session() as session:
                workflow = await session.get(Workflow, workflow_id)
                for i in range(workflow.current_step, len(workflow.steps)):
                    if not self._running:
                        break

                    step = workflow.steps[i]
                    step_type = step.get("type")
                    handler = self._handlers.get(step_type)

                    if not handler:
                        raise ValueError(f"No handler for step type: {step_type}")

                    logger.info(f"Executing step {i}: {step_type}")
                    result = await handler(workflow.context, step.get("params", {}))

                    workflow.current_step = i + 1
                    workflow.progress = ((i + 1) / len(workflow.steps)) * 100
                    workflow.context.update(result or {})
                    await session.commit()

                workflow.status = WorkflowStatus.COMPLETED
                workflow.completed_at = datetime.now(timezone.utc)
                workflow.progress = 100.0
                await session.commit()
                logger.info(f"Workflow {workflow_id} completed")

        except Exception as e:
            logger.error(f"Workflow {workflow_id} failed: {e}")
            async with self.db.get_session() as session:
                workflow = await session.get(Workflow, workflow_id)
                workflow.status = WorkflowStatus.FAILED
                workflow.error = str(e)
                await session.commit()
            raise

    async def _handle_document_upload(self, context: Dict, params: Dict) -> Dict:
        logger.info(f"Handling document upload: {params}")
        return {"upload_status": "completed"}

    async def _handle_document_process(self, context: Dict, params: Dict) -> Dict:
        doc_id = context.get("document_id")
        pipeline = params.get("pipeline", "default")
        logger.info(f"Processing document {doc_id} with pipeline {pipeline}")

        await self.kafka.publish(
            "document-processing",
            doc_id,
            {
                "document_id": doc_id,
                "pipeline_type": pipeline,
                "user_id": context.get("user_id"),
            },
        )
        return {"processing_status": "queued", "pipeline": pipeline}

    async def _handle_index_document(self, context: Dict, params: Dict) -> Dict:
        logger.info(f"Indexing document: {context.get('document_id')}")
        return {"index_status": "completed"}

    async def _handle_notify_user(self, context: Dict, params: Dict) -> Dict:
        user_id = context.get("user_id")
        message = params.get("message", "Workflow completed")
        logger.info(f"Notifying user {user_id}: {message}")

        await self.kafka.publish(
            "notifications",
            user_id,
            {
                "type": "workflow_notification",
                "user_id": user_id,
                "title": params.get("title", "Workflow Update"),
                "body": message,
                "channel": "in_app",
            },
        )
        return {"notification_sent": True}

    async def _handle_generate_report(self, context: Dict, params: Dict) -> Dict:
        logger.info(f"Generating report: {params}")
        return {"report_status": "generated", "report_url": "reports/sample.pdf"}
