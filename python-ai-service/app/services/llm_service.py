"""
LLM 服务 - 统一 AI 调用接口
"""
from dataclasses import dataclass
from typing import AsyncGenerator, Optional

import structlog

from app.config import settings
from app.providers.openai_provider import OpenAIProvider
from app.providers.anthropic_provider import AnthropicProvider
from app.providers.gemini_provider import GeminiProvider

logger = structlog.get_logger()


@dataclass
class GenerateChunk:
    """生成块"""
    content: str
    done: bool = False
    error: Optional[str] = None
    usage: Optional[dict] = None


class LLMService:
    """统一 LLM 调用服务"""

    def __init__(self):
        self.providers = {
            "openai": OpenAIProvider(),
            "anthropic": AnthropicProvider(),
            "gemini": GeminiProvider(),
        }

    async def generate_stream(
        self,
        prompt: str,
        provider: str = "openai",
        model: Optional[str] = None,
        temperature: float = 0.7,
        max_tokens: int = 4096,
        system_prompt: Optional[str] = None,
    ) -> AsyncGenerator[GenerateChunk, None]:
        """流式生成文本"""
        provider_instance = self.providers.get(provider)
        if not provider_instance:
            yield GenerateChunk(
                content="",
                done=True,
                error=f"Unknown provider: {provider}",
            )
            return

        model = model or self._get_default_model(provider)

        logger.info(
            "Starting generation",
            provider=provider,
            model=model,
            temperature=temperature,
            max_tokens=max_tokens,
        )

        total_content = ""
        input_tokens = 0
        output_tokens = 0

        try:
            async for chunk in provider_instance.generate_stream(
                prompt=prompt,
                model=model,
                temperature=temperature,
                max_tokens=max_tokens,
                system_prompt=system_prompt,
            ):
                total_content += chunk
                output_tokens += 1  # 简化计算

                yield GenerateChunk(content=chunk, done=False)

            # 最终响应
            yield GenerateChunk(
                content="",
                done=True,
                usage={
                    "input_tokens": input_tokens,
                    "output_tokens": output_tokens,
                    "total_tokens": input_tokens + output_tokens,
                    "model": model,
                    "provider": provider,
                },
            )

        except Exception as e:
            logger.error("Generation error", error=str(e))
            yield GenerateChunk(
                content="",
                done=True,
                error=str(e),
            )

    def _get_default_model(self, provider: str) -> str:
        """获取提供商默认模型"""
        defaults = {
            "openai": "gpt-4o",
            "anthropic": "claude-sonnet-4-20250514",
            "gemini": "gemini-2.5-pro",
        }
        return defaults.get(provider, "gpt-4o")

    async def generate(
        self,
        prompt: str,
        provider: str = "openai",
        model: Optional[str] = None,
        temperature: float = 0.7,
        max_tokens: int = 4096,
        system_prompt: Optional[str] = None,
    ) -> str:
        """非流式生成文本"""
        content = ""
        async for chunk in self.generate_stream(
            prompt=prompt,
            provider=provider,
            model=model,
            temperature=temperature,
            max_tokens=max_tokens,
            system_prompt=system_prompt,
        ):
            if not chunk.done:
                content += chunk.content
            elif chunk.error:
                raise Exception(chunk.error)
        return content
