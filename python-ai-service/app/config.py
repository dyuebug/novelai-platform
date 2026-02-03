"""
配置管理
"""
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    """应用配置"""

    # 服务配置
    debug: bool = False
    http_port: int = 8001
    grpc_port: int = 50051

    # 数据库配置
    database_url: str = "postgresql+asyncpg://postgres:postgres@localhost:5432/novelai"

    # AI 提供商配置
    openai_api_key: str = ""
    openai_base_url: str = "https://api.openai.com/v1"
    anthropic_api_key: str = ""
    gemini_api_key: str = ""

    # 默认模型配置
    default_provider: str = "openai"
    default_model: str = "gpt-4o"
    default_temperature: float = 0.7
    default_max_tokens: int = 4096

    # 向量配置
    default_embedding_model: str = "text-embedding-3-small"
    default_embedding_dim: int = 1536

    # gRPC TLS 配置
    grpc_tls_enabled: bool = False
    grpc_cert_file: str = ""  # 服务端证书
    grpc_key_file: str = ""   # 服务端私钥
    grpc_ca_file: str = ""    # CA 证书 (用于验证客户端)

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"


settings = Settings()
