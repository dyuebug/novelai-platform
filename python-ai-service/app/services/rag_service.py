"""
RAG 检索服务
"""
from dataclasses import dataclass
from typing import Optional

import structlog

from app.config import settings
from app.services.embedding_service import EmbeddingService
from app.repositories.embedding_repository import EmbeddingRepository, EmbeddingRecord

logger = structlog.get_logger()


@dataclass
class RetrievedItem:
    """检索项"""
    content: str
    score: float
    source_type: str
    source_id: str
    chapter_number: Optional[int] = None


@dataclass
class RetrieveResult:
    """检索结果"""
    query: str
    results: dict[str, list[RetrievedItem]]
    total_tokens: int


class RAGService:
    """RAG 检索服务"""

    def __init__(self):
        self.embedding_service = EmbeddingService()
        self.embedding_repo = EmbeddingRepository()

    async def retrieve_context(
        self,
        project_id: str,
        query: str,
        top_k: int = 10,
        layers: Optional[list[str]] = None,
        chapter_range: Optional[tuple[int, int]] = None,
        min_score: float = 0.7,
    ) -> RetrieveResult:
        """
        分层检索上下文

        Layers:
        - settings: 世界观设定
        - summaries: 章节摘要
        - chunks: 内容片段
        - characters: 角色信息
        """
        layers = layers or ["settings", "summaries", "chunks"]

        logger.info(
            "Retrieving context",
            project_id=project_id,
            query=query[:50],
            top_k=top_k,
            layers=layers,
        )

        # 生成查询向量
        query_embedding = await self.embedding_service.generate_embedding(query)

        results = {}
        total_tokens = 0

        for layer in layers:
            layer_results = await self._retrieve_layer(
                project_id=project_id,
                query_embedding=query_embedding,
                layer=layer,
                top_k=top_k,
                chapter_range=chapter_range,
                min_score=min_score,
            )
            results[layer] = layer_results
            total_tokens += sum(len(item.content.split()) for item in layer_results)

        return RetrieveResult(
            query=query,
            results=results,
            total_tokens=total_tokens,
        )

    async def _retrieve_layer(
        self,
        project_id: str,
        query_embedding: list[float],
        layer: str,
        top_k: int,
        chapter_range: Optional[tuple[int, int]] = None,
        min_score: float = 0.7,
    ) -> list[RetrievedItem]:
        """检索单个层级"""
        logger.info(
            "Retrieving layer",
            layer=layer,
            top_k=top_k,
        )

        # 执行向量搜索
        records = await self.embedding_repo.search(
            project_id=project_id,
            query_embedding=query_embedding,
            content_type=layer,
            top_k=top_k,
            chapter_range=chapter_range,
            min_score=min_score,
        )

        # 转换为 RetrievedItem
        items = [
            RetrievedItem(
                content=record.content,
                score=record.score,
                source_type=record.content_type,
                source_id=str(record.content_id),
                chapter_number=record.chapter_number,
            )
            for record in records
        ]

        logger.info(
            "Layer retrieval completed",
            layer=layer,
            results_count=len(items),
        )

        return items

    async def update_embeddings(
        self,
        project_id: str,
        content_type: str,
        content_id: str,
        content: str,
        chapter_number: Optional[int] = None,
        metadata: Optional[dict] = None,
    ) -> None:
        """更新向量嵌入"""
        logger.info(
            "Updating embeddings",
            project_id=project_id,
            content_type=content_type,
            content_id=content_id,
        )

        # 生成嵌入
        embedding = await self.embedding_service.generate_embedding(content)

        # 保存到数据库
        await self.embedding_repo.upsert(
            project_id=project_id,
            content_type=content_type,
            content_id=content_id,
            content=content,
            embedding=embedding,
            chapter_number=chapter_number,
            metadata=metadata,
        )

        logger.info(
            "Embeddings updated successfully",
            project_id=project_id,
            content_type=content_type,
            content_id=content_id,
        )

    async def delete_embeddings(
        self,
        project_id: str,
        content_type: str,
        content_id: str,
    ) -> None:
        """删除向量嵌入"""
        logger.info(
            "Deleting embeddings",
            project_id=project_id,
            content_type=content_type,
            content_id=content_id,
        )

        await self.embedding_repo.delete_by_content_id(
            project_id=project_id,
            content_type=content_type,
            content_id=content_id,
        )

    async def delete_project_embeddings(self, project_id: str) -> None:
        """删除项目的所有向量嵌入"""
        logger.info(
            "Deleting all project embeddings",
            project_id=project_id,
        )

        await self.embedding_repo.delete_by_project(project_id)

    async def batch_update_embeddings(
        self,
        project_id: str,
        items: list[dict],
    ) -> int:
        """
        批量更新向量嵌入

        items: [
            {
                "content_type": "chunks",
                "content_id": "uuid",
                "content": "text",
                "chapter_number": 1,  # optional
                "metadata": {}  # optional
            }
        ]
        """
        logger.info(
            "Batch updating embeddings",
            project_id=project_id,
            items_count=len(items),
        )

        # 批量生成嵌入
        contents = [item["content"] for item in items]
        embeddings = await self.embedding_service.generate_embeddings(contents)

        # 逐个保存
        for item, embedding in zip(items, embeddings):
            await self.embedding_repo.upsert(
                project_id=project_id,
                content_type=item["content_type"],
                content_id=item["content_id"],
                content=item["content"],
                embedding=embedding,
                chapter_number=item.get("chapter_number"),
                metadata=item.get("metadata"),
            )

        logger.info(
            "Batch update completed",
            project_id=project_id,
            updated_count=len(items),
        )

        return len(items)
