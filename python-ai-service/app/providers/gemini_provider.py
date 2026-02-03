"""
Gemini Provider
"""
from typing import AsyncGenerator, Optional

import structlog
import google.generativeai as genai

from app.config import settings

logger = structlog.get_logger()


class GeminiProvider:
    """Gemini 提供商"""

    def __init__(self):
        genai.configure(api_key=settings.gemini_api_key)

    async def generate_stream(
        self,
        prompt: str,
        model: str = "gemini-2.5-pro",
        temperature: float = 0.7,
        max_tokens: int = 4096,
        system_prompt: Optional[str] = None,
    ) -> AsyncGenerator[str, None]:
        """流式生成文本"""
        try:
            model_instance = genai.GenerativeModel(
                model_name=model,
                system_instruction=system_prompt,
                generation_config=genai.GenerationConfig(
                    temperature=temperature,
                    max_output_tokens=max_tokens,
                ),
            )

            response = await model_instance.generate_content_async(
                prompt,
                stream=True,
            )

            async for chunk in response:
                if chunk.text:
                    yield chunk.text

        except Exception as e:
            logger.error("Gemini generation error", error=str(e))
            raise

    async def generate(
        self,
        prompt: str,
        model: str = "gemini-2.5-pro",
        temperature: float = 0.7,
        max_tokens: int = 4096,
        system_prompt: Optional[str] = None,
    ) -> str:
        """非流式生成文本"""
        model_instance = genai.GenerativeModel(
            model_name=model,
            system_instruction=system_prompt,
            generation_config=genai.GenerationConfig(
                temperature=temperature,
                max_output_tokens=max_tokens,
            ),
        )

        response = await model_instance.generate_content_async(prompt)
        return response.text or ""
