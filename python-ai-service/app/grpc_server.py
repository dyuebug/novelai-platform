"""
gRPC 服务器
"""
import asyncio
from concurrent import futures
from dataclasses import asdict

import grpc
import structlog

from app.config import settings
from app.services.llm_service import LLMService
from app.services.rag_service import RAGService
from app.services.quality_service import QualityService
from app.services.embedding_service import EmbeddingService
from app.services.constitution_service import ConstitutionService

# 导入生成的 protobuf 代码
from app.proto import ai_service_pb2, ai_service_pb2_grpc

logger = structlog.get_logger()


class AIServiceServicer:
    """AI 服务 gRPC 实现"""

    def __init__(self):
        self.llm_service = LLMService()
        self.rag_service = RAGService()
        self.quality_service = QualityService()
        self.embedding_service = EmbeddingService()
        self.constitution_service = ConstitutionService()

    async def GenerateStream(self, request, context):
        """流式文本生成"""
        logger.info(
            "GenerateStream called",
            provider=request.provider or settings.default_provider,
            model=request.model or settings.default_model,
            user_id=request.user_id,
        )

        try:
            async for chunk in self.llm_service.generate_stream(
                prompt=request.prompt,
                provider=request.provider or settings.default_provider,
                model=request.model or settings.default_model,
                temperature=request.temperature or settings.default_temperature,
                max_tokens=request.max_tokens or settings.default_max_tokens,
                system_prompt=request.system_prompt,
            ):
                # 返回流式响应
                yield {
                    "content": chunk.content,
                    "done": chunk.done,
                    "error": chunk.error or "",
                    "usage": chunk.usage,
                }
        except Exception as e:
            logger.error("GenerateStream error", error=str(e))
            yield {
                "content": "",
                "done": True,
                "error": str(e),
            }

    async def RetrieveContext(self, request, context):
        """RAG 检索"""
        logger.info(
            "RetrieveContext called",
            project_id=request.project_id,
            query=request.query[:50],
            top_k=request.top_k,
        )

        try:
            result = await self.rag_service.retrieve_context(
                project_id=request.project_id,
                query=request.query,
                top_k=request.top_k or 10,
                layers=list(request.layers) or ["settings", "summaries", "chunks"],
                chapter_range=(request.chapter_range_start, request.chapter_range_end),
            )
            return result
        except Exception as e:
            logger.error("RetrieveContext error", error=str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return {}

    async def AnalyzeReadingPower(self, request, context):
        """追读力分析 (T7.1)"""
        logger.info(
            "AnalyzeReadingPower called",
            chapter_id=request.chapter_id,
            content_length=len(request.content),
        )

        try:
            result = await self.quality_service.analyze_reading_power(
                chapter_id=request.chapter_id,
                content=request.content,
                provider=request.provider or settings.default_provider,
                model=request.model or settings.default_model,
            )
            return asdict(result)
        except Exception as e:
            logger.error("AnalyzeReadingPower error", error=str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return {}

    async def CheckConsistency(self, request, context):
        """一致性检查 (T7.2)"""
        logger.info(
            "CheckConsistency called",
            chapter_id=request.chapter_id,
            content_length=len(request.content),
        )

        try:
            # 构建上下文
            ctx = {
                "characters": list(request.characters) if request.characters else [],
                "world_settings": list(request.world_settings) if request.world_settings else [],
                "previous_summary": request.previous_summary or "",
            }

            result = await self.quality_service.check_consistency(
                chapter_id=request.chapter_id,
                content=request.content,
                context=ctx,
                provider=request.provider or settings.default_provider,
                model=request.model or settings.default_model,
            )
            return asdict(result)
        except Exception as e:
            logger.error("CheckConsistency error", error=str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return {}

    async def MultiAgentReview(self, request, context):
        """多Agent审查 (T7.3)"""
        logger.info(
            "MultiAgentReview called",
            chapter_id=request.chapter_id,
            content_length=len(request.content),
        )

        try:
            result = await self.quality_service.multi_agent_review(
                chapter_id=request.chapter_id,
                content=request.content,
                provider=request.provider or settings.default_provider,
                model=request.model or settings.default_model,
            )
            return asdict(result)
        except Exception as e:
            logger.error("MultiAgentReview error", error=str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return {}

    async def EvaluateQuality(self, request, context):
        """综合质量评估"""
        logger.info(
            "EvaluateQuality called",
            chapter_id=request.chapter_id,
        )

        try:
            # 构建上下文
            ctx = {
                "characters": list(request.characters) if request.characters else [],
                "world_settings": list(request.world_settings) if request.world_settings else [],
                "previous_summary": request.previous_summary or "",
            }

            result = await self.quality_service.evaluate_quality(
                chapter_id=request.chapter_id,
                content=request.content,
                context=ctx,
                provider=request.provider or settings.default_provider,
                model=request.model or settings.default_model,
            )
            return result
        except Exception as e:
            logger.error("EvaluateQuality error", error=str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return {}

    async def GenerateEmbedding(self, request, context):
        """生成向量嵌入"""
        logger.info(
            "GenerateEmbedding called",
            text_count=len(request.texts),
            model=request.model,
        )

        try:
            result = await self.embedding_service.generate_embeddings(
                texts=list(request.texts),
                model=request.model or settings.default_embedding_model,
            )
            return result
        except Exception as e:
            logger.error("GenerateEmbedding error", error=str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return {}

    async def AnalyzeContent(self, request, context):
        """内容分析 (兼容旧接口)"""
        logger.info(
            "AnalyzeContent called",
            chapter_id=request.chapter_id,
            content_length=len(request.content),
        )

        try:
            result = await self.quality_service.analyze_reading_power(
                chapter_id=request.chapter_id,
                content=request.content,
            )
            return asdict(result)
        except Exception as e:
            logger.error("AnalyzeContent error", error=str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return {}

    async def CheckConstraints(self, request, context):
        """约束检查 (T8.1)"""
        logger.info(
            "CheckConstraints called",
            project_id=request.project_id,
            chapter_id=request.chapter_id,
            content_length=len(request.content),
        )

        try:
            # 转换约束格式
            constraints = []
            for c in request.constraints:
                constraints.append({
                    "id": c.id,
                    "type": c.type,
                    "description": c.description,
                    "rule": c.rule,
                })

            result = await self.constitution_service.check_constraints(
                project_id=request.project_id,
                chapter_id=request.chapter_id,
                content=request.content,
                constraints=constraints,
            )

            # 转换响应格式
            violations = []
            for v in result.get("violations", []):
                violations.append({
                    "constraint_id": v.get("constraint_id", ""),
                    "constraint_type": v.get("constraint_type", ""),
                    "description": v.get("description", ""),
                    "severity": v.get("severity", ""),
                    "suggestion": v.get("suggestion", ""),
                    "can_exempt": v.get("can_exempt", False),
                })

            return {
                "violations": violations,
                "passed": result.get("passed", True),
            }
        except Exception as e:
            logger.error("CheckConstraints error", error=str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return {"violations": [], "passed": False}

    async def RequestExemption(self, request, context):
        """请求豁免 (T8.1)"""
        logger.info(
            "RequestExemption called",
            project_id=request.project_id,
            chapter_id=request.chapter_id,
            constraint_id=request.constraint_id,
        )

        try:
            result = await self.constitution_service.request_exemption(
                project_id=request.project_id,
                chapter_id=request.chapter_id,
                constraint_id=request.constraint_id,
                reason=request.reason,
                requested_by=request.requested_by,
            )

            return {
                "id": result.get("id", ""),
                "constraint_id": result.get("constraint_id", ""),
                "chapter_id": result.get("chapter_id", ""),
                "reason": result.get("reason", ""),
                "status": result.get("status", ""),
                "created_at": result.get("created_at", ""),
            }
        except Exception as e:
            logger.error("RequestExemption error", error=str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return {}

    async def RevokeExemption(self, request, context):
        """撤销豁免 (T8.1)"""
        logger.info(
            "RevokeExemption called",
            exemption_id=request.exemption_id,
        )

        try:
            await self.constitution_service.revoke_exemption(
                exemption_id=request.exemption_id,
            )
            return {
                "success": True,
                "message": "豁免已撤销",
            }
        except Exception as e:
            logger.error("RevokeExemption error", error=str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return {
                "success": False,
                "message": str(e),
            }

    async def UpdateEmbedding(self, request, context):
        """更新向量索引 (T2.3)"""
        logger.info(
            "UpdateEmbedding called",
            project_id=request.project_id,
            content_type=request.content_type,
            content_id=request.content_id,
        )

        try:
            # 转换 metadata
            metadata = dict(request.metadata) if request.metadata else None

            await self.rag_service.update_embeddings(
                project_id=request.project_id,
                content_type=request.content_type,
                content_id=request.content_id,
                content=request.content,
                chapter_number=request.chapter_number if request.chapter_number else None,
                metadata=metadata,
            )
            return {
                "success": True,
                "message": "向量索引更新成功",
            }
        except Exception as e:
            logger.error("UpdateEmbedding error", error=str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return {
                "success": False,
                "message": str(e),
            }

    async def BatchUpdateEmbeddings(self, request, context):
        """批量更新向量索引 (T2.3)"""
        logger.info(
            "BatchUpdateEmbeddings called",
            project_id=request.project_id,
            items_count=len(request.items),
        )

        try:
            # 转换 items 格式
            items = []
            for item in request.items:
                items.append({
                    "content_type": item.content_type,
                    "content_id": item.content_id,
                    "content": item.content,
                    "chapter_number": item.chapter_number if item.chapter_number else None,
                    "metadata": dict(item.metadata) if item.metadata else None,
                })

            updated_count = await self.rag_service.batch_update_embeddings(
                project_id=request.project_id,
                items=items,
            )
            return {
                "success": True,
                "updated_count": updated_count,
                "message": f"成功更新 {updated_count} 条向量索引",
            }
        except Exception as e:
            logger.error("BatchUpdateEmbeddings error", error=str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return {
                "success": False,
                "updated_count": 0,
                "message": str(e),
            }

    async def DeleteEmbedding(self, request, context):
        """删除向量索引 (T2.3)"""
        logger.info(
            "DeleteEmbedding called",
            project_id=request.project_id,
            content_type=request.content_type,
            content_id=request.content_id,
        )

        try:
            if not request.content_type and not request.content_id:
                # 删除项目所有嵌入
                await self.rag_service.delete_project_embeddings(
                    project_id=request.project_id,
                )
                return {
                    "success": True,
                    "message": "项目所有向量索引已删除",
                }
            else:
                # 删除指定内容的嵌入
                await self.rag_service.delete_embeddings(
                    project_id=request.project_id,
                    content_type=request.content_type,
                    content_id=request.content_id,
                )
                return {
                    "success": True,
                    "message": "向量索引已删除",
                }
        except Exception as e:
            logger.error("DeleteEmbedding error", error=str(e))
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(str(e))
            return {
                "success": False,
                "message": str(e),
            }


async def serve_grpc():
    """启动 gRPC 服务器"""
    server = grpc.aio.server(
        futures.ThreadPoolExecutor(max_workers=10),
        options=[
            ("grpc.max_send_message_length", 50 * 1024 * 1024),
            ("grpc.max_receive_message_length", 50 * 1024 * 1024),
        ],
    )

    # 注册服务
    ai_service_pb2_grpc.add_AIServiceServicer_to_server(
        AIServiceServicer(), server
    )

    listen_addr = f"[::]:{settings.grpc_port}"

    # 配置 TLS
    if settings.grpc_tls_enabled:
        server_credentials = load_server_credentials()
        if server_credentials:
            server.add_secure_port(listen_addr, server_credentials)
            logger.info("Starting gRPC server with mTLS", address=listen_addr)
        else:
            logger.warning("TLS enabled but credentials not loaded, falling back to insecure")
            server.add_insecure_port(listen_addr)
    else:
        server.add_insecure_port(listen_addr)
        logger.info("Starting gRPC server (insecure)", address=listen_addr)

    await server.start()

    try:
        await server.wait_for_termination()
    except asyncio.CancelledError:
        logger.info("Stopping gRPC server")
        await server.stop(grace=5)


def load_server_credentials():
    """加载 mTLS 服务端凭证"""
    try:
        # 读取证书文件
        with open(settings.grpc_cert_file, "rb") as f:
            server_cert = f.read()
        with open(settings.grpc_key_file, "rb") as f:
            server_key = f.read()

        # 读取 CA 证书 (用于验证客户端)
        root_certs = None
        if settings.grpc_ca_file:
            with open(settings.grpc_ca_file, "rb") as f:
                root_certs = f.read()

        # 创建服务端凭证
        return grpc.ssl_server_credentials(
            [(server_key, server_cert)],
            root_certificates=root_certs,
            require_client_auth=root_certs is not None,  # 如果有 CA 证书则要求客户端认证
        )
    except FileNotFoundError as e:
        logger.error("TLS certificate file not found", error=str(e))
        return None
    except Exception as e:
        logger.error("Failed to load TLS credentials", error=str(e))
        return None
