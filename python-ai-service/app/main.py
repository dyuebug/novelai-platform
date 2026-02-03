"""
NovelAI Python AI Service
AI 服务层入口
"""
import asyncio
from contextlib import asynccontextmanager

import structlog
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.config import settings
from app.grpc_server import serve_grpc

logger = structlog.get_logger()


@asynccontextmanager
async def lifespan(app: FastAPI):
    """应用生命周期管理"""
    logger.info("Starting AI Service", port=settings.grpc_port)

    # 启动 gRPC 服务器
    grpc_task = asyncio.create_task(serve_grpc())

    yield

    # 关闭 gRPC 服务器
    grpc_task.cancel()
    try:
        await grpc_task
    except asyncio.CancelledError:
        pass

    logger.info("AI Service stopped")


app = FastAPI(
    title="NovelAI AI Service",
    description="AI 服务层：LLM 调用、RAG 检索、质量评估",
    version="1.0.0",
    lifespan=lifespan,
)

# CORS 中间件
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/health")
async def health_check():
    """健康检查"""
    return {"status": "healthy", "service": "ai-service"}


@app.get("/")
async def root():
    """根路径"""
    return {
        "service": "NovelAI AI Service",
        "version": "1.0.0",
        "grpc_port": settings.grpc_port,
    }


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=settings.http_port,
        reload=settings.debug,
    )
