"""
向量嵌入仓储
"""
from dataclasses import dataclass
from typing import Optional
from uuid import UUID

import structlog
from pgvector.sqlalchemy import Vector
from sqlalchemy import Column, String, Text, Integer, DateTime, text, select, delete
from sqlalchemy.dialects.postgresql import UUID as PG_UUID, JSONB
from sqlalchemy.sql import func

from app.config import settings
from app.db import Base, get_session

logger = structlog.get_logger()


class EmbeddingModel(Base):
    """向量嵌入 ORM 模型"""
    __tablename__ = "embeddings"

    id = Column(PG_UUID(as_uuid=True), primary_key=True, server_default=text("uuid_generate_v4()"))
    project_id = Column(PG_UUID(as_uuid=True), nullable=False, index=True)
    content_type = Column(String(50), nullable=False, index=True)
    content_id = Column(PG_UUID(as_uuid=True), nullable=False, index=True)
    content = Column(Text, nullable=False)
    chapter_number = Column(Integer, nullable=True)
    embedding = Column(Vector(settings.default_embedding_dim), nullable=False)
    metadata = Column(JSONB, default={})
    created_at = Column(DateTime, server_default=func.now())
    updated_at = Column(DateTime, server_default=func.now(), onupdate=func.now())


@dataclass
class EmbeddingRecord:
    """嵌入记录"""
    id: UUID
    project_id: UUID
    content_type: str
    content_id: UUID
    content: str
    chapter_number: Optional[int]
    score: float  # 相似度分数


class EmbeddingRepository:
    """向量嵌入仓储"""

    async def upsert(
        self,
        project_id: str,
        content_type: str,
        content_id: str,
        content: str,
        embedding: list[float],
        chapter_number: Optional[int] = None,
        metadata: Optional[dict] = None,
    ) -> None:
        """插入或更新向量嵌入"""
        async with get_session() as session:
            # 检查是否存在
            stmt = select(EmbeddingModel).where(
                EmbeddingModel.project_id == UUID(project_id),
                EmbeddingModel.content_type == content_type,
                EmbeddingModel.content_id == UUID(content_id),
            )
            result = await session.execute(stmt)
            existing = result.scalar_one_or_none()

            if existing:
                # 更新
                existing.content = content
                existing.embedding = embedding
                existing.chapter_number = chapter_number
                if metadata:
                    existing.metadata = metadata
                logger.info(
                    "Updated embedding",
                    project_id=project_id,
                    content_type=content_type,
                    content_id=content_id,
                )
            else:
                # 插入
                new_embedding = EmbeddingModel(
                    project_id=UUID(project_id),
                    content_type=content_type,
                    content_id=UUID(content_id),
                    content=content,
                    embedding=embedding,
                    chapter_number=chapter_number,
                    metadata=metadata or {},
                )
                session.add(new_embedding)
                logger.info(
                    "Inserted embedding",
                    project_id=project_id,
                    content_type=content_type,
                    content_id=content_id,
                )

    async def search(
        self,
        project_id: str,
        query_embedding: list[float],
        content_type: str,
        top_k: int = 10,
        chapter_range: Optional[tuple[int, int]] = None,
        min_score: float = 0.7,
    ) -> list[EmbeddingRecord]:
        """向量相似度搜索"""
        async with get_session() as session:
            # 构建基础查询 - 使用余弦距离
            # pgvector 的 <=> 操作符返回余弦距离 (1 - 相似度)
            # 所以相似度 = 1 - 距离
            distance_expr = EmbeddingModel.embedding.cosine_distance(query_embedding)

            stmt = (
                select(
                    EmbeddingModel,
                    (1 - distance_expr).label("score"),
                )
                .where(
                    EmbeddingModel.project_id == UUID(project_id),
                    EmbeddingModel.content_type == content_type,
                )
                .order_by(distance_expr)
                .limit(top_k)
            )

            # 章节范围过滤
            if chapter_range and content_type in ("chunks", "summaries"):
                stmt = stmt.where(
                    EmbeddingModel.chapter_number >= chapter_range[0],
                    EmbeddingModel.chapter_number <= chapter_range[1],
                )

            result = await session.execute(stmt)
            rows = result.all()

            # 过滤低分结果并转换
            records = []
            for row in rows:
                embedding_model = row[0]
                score = float(row[1])

                if score >= min_score:
                    records.append(EmbeddingRecord(
                        id=embedding_model.id,
                        project_id=embedding_model.project_id,
                        content_type=embedding_model.content_type,
                        content_id=embedding_model.content_id,
                        content=embedding_model.content,
                        chapter_number=embedding_model.chapter_number,
                        score=score,
                    ))

            logger.info(
                "Vector search completed",
                project_id=project_id,
                content_type=content_type,
                results_count=len(records),
            )

            return records

    async def delete_by_content_id(
        self,
        project_id: str,
        content_type: str,
        content_id: str,
    ) -> None:
        """删除指定内容的嵌入"""
        async with get_session() as session:
            stmt = delete(EmbeddingModel).where(
                EmbeddingModel.project_id == UUID(project_id),
                EmbeddingModel.content_type == content_type,
                EmbeddingModel.content_id == UUID(content_id),
            )
            await session.execute(stmt)
            logger.info(
                "Deleted embedding",
                project_id=project_id,
                content_type=content_type,
                content_id=content_id,
            )

    async def delete_by_project(self, project_id: str) -> None:
        """删除项目的所有嵌入"""
        async with get_session() as session:
            stmt = delete(EmbeddingModel).where(
                EmbeddingModel.project_id == UUID(project_id),
            )
            result = await session.execute(stmt)
            logger.info(
                "Deleted project embeddings",
                project_id=project_id,
                deleted_count=result.rowcount,
            )
