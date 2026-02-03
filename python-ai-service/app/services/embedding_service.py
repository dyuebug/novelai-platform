"""
向量嵌入服务
"""
from typing import Optional

import structlog

from app.config import settings
from app.providers.openai_provider import OpenAIProvider

logger = structlog.get_logger()


class EmbeddingService:
    """向量嵌入服务"""

    def __init__(self):
        self.openai_provider = OpenAIProvider()

    async def generate_embedding(
        self,
        text: str,
        model: Optional[str] = None,
    ) -> list[float]:
        """生成单个文本的向量嵌入"""
        model = model or settings.default_embedding_model
        embeddings = await self.generate_embeddings([text], model)
        return embeddings[0]

    async def generate_embeddings(
        self,
        texts: list[str],
        model: Optional[str] = None,
    ) -> list[list[float]]:
        """批量生成向量嵌入"""
        model = model or settings.default_embedding_model

        logger.info(
            "Generating embeddings",
            text_count=len(texts),
            model=model,
        )

        # 使用 OpenAI 生成嵌入
        embeddings = await self.openai_provider.generate_embedding(texts, model)

        return embeddings

    def get_dimension(self, model: Optional[str] = None) -> int:
        """获取模型的向量维度"""
        model = model or settings.default_embedding_model

        dimensions = {
            "text-embedding-ada-002": 1536,
            "text-embedding-3-small": 1536,
            "text-embedding-3-large": 3072,
        }

        return dimensions.get(model, 1536)
