"""
Anthropic Provider
"""
from typing import AsyncGenerator, Optional

import structlog
from anthropic import AsyncAnthropic

from app.config import settings

logger = structlog.get_logger()


class AnthropicProvider:
    """Anthropic 提供商"""

    def __init__(self):
        self.client = AsyncAnthropic(api_key=settings.anthropic_api_key)

    async def generate_stream(
        self,
        prompt: str,
        model: str = "claude-sonnet-4-20250514",
        temperature: float = 0.7,
        max_tokens: int = 4096,
        system_prompt: Optional[str] = None,
    ) -> AsyncGenerator[str, None]:
        """流式生成文本"""
        try:
            async with self.client.messages.stream(
                model=model,
                max_tokens=max_tokens,
                messages=[{"role": "user", "content": prompt}],
                system=system_prompt or "",
                temperature=temperature,
            ) as stream:
                async for text in stream.text_stream:
                    yield text

        except Exception as e:
            logger.error("Anthropic generation error", error=str(e))
            raise

    async def generate(
        self,
        prompt: str,
        model: str = "claude-sonnet-4-20250514",
        temperature: float = 0.7,
        max_tokens: int = 4096,
        system_prompt: Optional[str] = None,
    ) -> str:
        """非流式生成文本"""
        response = await self.client.messages.create(
            model=model,
            max_tokens=max_tokens,
            messages=[{"role": "user", "content": prompt}],
            system=system_prompt or "",
            temperature=temperature,
        )

        return response.content[0].text if response.content else ""
