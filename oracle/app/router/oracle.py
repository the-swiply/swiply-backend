from fastapi import APIRouter

router = APIRouter()


@router.get("/stub")
def stub():
    return {"status": "ok"}
