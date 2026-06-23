from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import List, Optional, Dict

from app.models.workflow import WorkflowStatus
from config.settings import settings

router = APIRouter()


class WorkflowStep(BaseModel):
    type: str
    params: Dict = {}


class CreateWorkflowRequest(BaseModel):
    name: str
    steps: List[WorkflowStep]
    user_id: Optional[str] = None
    document_id: Optional[str] = None


class WorkflowResponse(BaseModel):
    id: str
    name: str
    status: str
    progress: float
    current_step: int
    total_steps: int
    created_at: str


@router.post("")
async def create_workflow(request: CreateWorkflowRequest):
    from app.main import app
    engine = app.state.workflow_engine

    workflow = await engine.create_workflow(
        name=request.name,
        steps=[s.dict() for s in request.steps],
        user_id=request.user_id,
        document_id=request.document_id,
    )

    return WorkflowResponse(
        id=workflow.id,
        name=workflow.name,
        status=workflow.status.value,
        progress=workflow.progress,
        current_step=workflow.current_step,
        total_steps=len(workflow.steps),
        created_at=workflow.created_at.isoformat(),
    )


@router.post("/{workflow_id}/execute")
async def execute_workflow(workflow_id: str):
    from app.main import app
    engine = app.state.workflow_engine

    try:
        await engine.execute_workflow(workflow_id)
        return {"message": "Workflow execution started", "workflow_id": workflow_id}
    except ValueError as e:
        raise HTTPException(status_code=404, detail=str(e))


@router.get("/{workflow_id}")
async def get_workflow(workflow_id: str):
    from app.main import app
    db = app.state.db

    async with db.get_session() as session:
        from sqlalchemy import select
        from app.models.workflow import Workflow
        result = await session.execute(
            select(Workflow).where(Workflow.id == workflow_id)
        )
        workflow = result.scalar_one_or_none()
        if not workflow:
            raise HTTPException(status_code=404, detail="Workflow not found")

        return WorkflowResponse(
            id=workflow.id,
            name=workflow.name,
            status=workflow.status.value,
            progress=workflow.progress,
            current_step=workflow.current_step,
            total_steps=len(workflow.steps),
            created_at=workflow.created_at.isoformat(),
        )


@router.get("")
async def list_workflows(page: int = 1, page_size: int = 20):
    from app.main import app
    db = app.state.db

    async with db.get_session() as session:
        from sqlalchemy import select, func
        from app.models.workflow import Workflow

        count_result = await session.execute(func.count(Workflow.id))
        total = count_result.scalar()

        result = await session.execute(
            select(Workflow)
            .order_by(Workflow.created_at.desc())
            .offset((page - 1) * page_size)
            .limit(page_size)
        )
        workflows = result.scalars().all()

        return {
            "workflows": [
                WorkflowResponse(
                    id=w.id, name=w.name, status=w.status.value,
                    progress=w.progress, current_step=w.current_step,
                    total_steps=len(w.steps),
                    created_at=w.created_at.isoformat(),
                )
                for w in workflows
            ],
            "total": total,
            "page": page,
            "page_size": page_size,
        }
