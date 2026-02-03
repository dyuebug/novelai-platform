"""
宪法约束服务 - 内容约束检查与豁免管理
"""
import json
import re
from dataclasses import dataclass, field
from typing import Optional
from enum import Enum

import structlog

from app.services.llm_service import LLMService

logger = structlog.get_logger()


class ConstraintType(str, Enum):
    """约束类型"""
    HARD = "hard"  # 硬约束 - 必须遵守
    SOFT = "soft"  # 软约束 - 可以豁免


class ConstraintCategory(str, Enum):
    """约束分类"""
    CHARACTER = "character"  # 角色约束
    PLOT = "plot"  # 情节约束
    WORLDVIEW = "worldview"  # 世界观约束
    STYLE = "style"  # 风格约束
    CONTENT = "content"  # 内容约束


@dataclass
class Constraint:
    """约束定义"""
    id: str
    name: str
    description: str
    type: ConstraintType
    category: ConstraintCategory
    rule: str  # 约束规则描述
    examples: list[str] = field(default_factory=list)  # 违规示例
    is_active: bool = True


@dataclass
class ConstraintViolation:
    """约束违规"""
    constraint_id: str
    constraint_name: str
    constraint_type: ConstraintType
    severity: str  # critical, major, minor
    description: str
    location: str  # 违规位置描述
    suggestion: str  # 修改建议
    can_exempt: bool = False  # 是否可以豁免


@dataclass
class ConstraintCheckResult:
    """约束检查结果"""
    passed: bool
    violations: list[ConstraintViolation] = field(default_factory=list)
    hard_violations: int = 0
    soft_violations: int = 0
    exempted_count: int = 0


@dataclass
class ExemptionRequest:
    """豁免请求"""
    constraint_id: str
    chapter_id: str
    reason: str
    approved: bool = False
    approved_by: Optional[str] = None


class ConstitutionService:
    """宪法约束服务"""

    def __init__(self):
        self.llm_service = LLMService()
        # 默认约束集
        self.default_constraints = self._init_default_constraints()
        # 豁免记录
        self.exemptions: dict[str, list[ExemptionRequest]] = {}

    def _init_default_constraints(self) -> list[Constraint]:
        """初始化默认约束"""
        return [
            # 硬约束 - 角色
            Constraint(
                id="char_death",
                name="主角死亡约束",
                description="主角不能在非结局章节死亡",
                type=ConstraintType.HARD,
                category=ConstraintCategory.CHARACTER,
                rule="主角在非结局章节中不能真正死亡，可以濒死但必须存活",
                examples=["主角被杀死", "主角永远闭上了眼睛"],
            ),
            Constraint(
                id="char_ooc",
                name="角色OOC约束",
                description="角色行为必须符合设定性格",
                type=ConstraintType.HARD,
                category=ConstraintCategory.CHARACTER,
                rule="角色的言行举止必须与其设定的性格特征一致",
                examples=["冷酷角色突然变得话痨", "胆小角色毫无理由地勇敢"],
            ),
            # 硬约束 - 世界观
            Constraint(
                id="world_power",
                name="力量体系约束",
                description="不能违反已设定的力量体系规则",
                type=ConstraintType.HARD,
                category=ConstraintCategory.WORLDVIEW,
                rule="角色的能力不能超出力量体系的设定范围",
                examples=["低级修士击败高级修士", "无魔法天赋者使用高级魔法"],
            ),
            Constraint(
                id="world_logic",
                name="世界逻辑约束",
                description="不能违反世界基本逻辑",
                type=ConstraintType.HARD,
                category=ConstraintCategory.WORLDVIEW,
                rule="情节发展必须符合世界观设定的基本逻辑",
                examples=["科技世界出现魔法", "古代世界出现现代科技"],
            ),
            # 软约束 - 情节
            Constraint(
                id="plot_pacing",
                name="节奏约束",
                description="情节节奏应保持合理",
                type=ConstraintType.SOFT,
                category=ConstraintCategory.PLOT,
                rule="避免情节过于拖沓或过于仓促",
                examples=["连续多章无情节推进", "重大事件一笔带过"],
            ),
            Constraint(
                id="plot_foreshadow",
                name="伏笔约束",
                description="重要伏笔应及时回收",
                type=ConstraintType.SOFT,
                category=ConstraintCategory.PLOT,
                rule="埋下的伏笔应在合理章节内回收",
                examples=["伏笔超过50章未回收", "伏笔被遗忘"],
            ),
            # 软约束 - 风格
            Constraint(
                id="style_tone",
                name="基调约束",
                description="保持作品整体基调一致",
                type=ConstraintType.SOFT,
                category=ConstraintCategory.STYLE,
                rule="章节风格应与作品整体基调保持一致",
                examples=["严肃作品突然搞笑", "轻松作品突然沉重"],
            ),
            # 软约束 - 内容
            Constraint(
                id="content_sensitive",
                name="敏感内容约束",
                description="避免过度敏感内容",
                type=ConstraintType.SOFT,
                category=ConstraintCategory.CONTENT,
                rule="避免过度血腥、暴力或不当内容",
                examples=["过度血腥描写", "不当内容"],
            ),
        ]

    def _parse_json_response(self, text: str) -> dict:
        """解析 LLM 返回的 JSON"""
        json_match = re.search(r'```json\s*([\s\S]*?)\s*```', text)
        if json_match:
            text = json_match.group(1)
        else:
            json_match = re.search(r'\{[\s\S]*\}', text)
            if json_match:
                text = json_match.group(0)

        try:
            return json.loads(text)
        except json.JSONDecodeError:
            logger.warning("Failed to parse JSON response", text=text[:200])
            return {}

    async def check_constraints(
        self,
        chapter_id: str,
        content: str,
        project_constraints: Optional[list[Constraint]] = None,
        context: Optional[dict] = None,
        provider: str = "openai",
        model: str = "gpt-4o-mini",
    ) -> ConstraintCheckResult:
        """检查内容是否违反约束"""
        logger.info(
            "Checking constraints",
            chapter_id=chapter_id,
            content_length=len(content),
        )

        # 合并默认约束和项目约束
        constraints = self.default_constraints.copy()
        if project_constraints:
            constraints.extend(project_constraints)

        # 过滤激活的约束
        active_constraints = [c for c in constraints if c.is_active]

        # 获取该章节的豁免列表
        chapter_exemptions = self.exemptions.get(chapter_id, [])
        exempted_ids = {e.constraint_id for e in chapter_exemptions if e.approved}

        # 构建约束描述
        constraints_desc = "\n".join([
            f"- [{c.type.value.upper()}] {c.name}: {c.rule}"
            for c in active_constraints
        ])

        # 构建上下文信息
        context_info = ""
        if context:
            if context.get("characters"):
                chars = context["characters"][:5]
                context_info += "\n角色设定:\n" + "\n".join([
                    f"- {c.get('name')}: {c.get('personality', '无')}"
                    for c in chars
                ])
            if context.get("world_settings"):
                settings = context["world_settings"][:5]
                context_info += "\n世界设定:\n" + "\n".join([
                    f"- {s.get('title')}: {s.get('content', '')[:50]}"
                    for s in settings
                ])

        prompt = f"""你是一位小说内容审核专家。请检查以下章节内容是否违反约束规则。

约束规则:
{constraints_desc}
{context_info}

章节内容:
---
{content[:3000]}
---

请检查内容是否违反上述约束，以 JSON 格式返回:
```json
{{
    "violations": [
        {{
            "constraint_name": "约束名称",
            "constraint_type": "hard/soft",
            "severity": "critical/major/minor",
            "description": "违规描述",
            "location": "违规位置(引用原文)",
            "suggestion": "修改建议"
        }}
    ]
}}
```

如果没有违规，返回空数组。"""

        try:
            response = await self.llm_service.generate(
                prompt=prompt,
                provider=provider,
                model=model,
                temperature=0.3,
                max_tokens=1000,
            )

            data = self._parse_json_response(response)
            violations = []
            hard_count = 0
            soft_count = 0
            exempted_count = 0

            for v in data.get("violations", []):
                # 查找对应的约束
                constraint = next(
                    (c for c in active_constraints if c.name == v.get("constraint_name")),
                    None
                )

                constraint_id = constraint.id if constraint else "unknown"
                constraint_type = ConstraintType(v.get("constraint_type", "soft"))

                # 检查是否已豁免
                if constraint_id in exempted_ids:
                    exempted_count += 1
                    continue

                violation = ConstraintViolation(
                    constraint_id=constraint_id,
                    constraint_name=v.get("constraint_name", "未知约束"),
                    constraint_type=constraint_type,
                    severity=v.get("severity", "minor"),
                    description=v.get("description", ""),
                    location=v.get("location", ""),
                    suggestion=v.get("suggestion", ""),
                    can_exempt=constraint_type == ConstraintType.SOFT,
                )
                violations.append(violation)

                if constraint_type == ConstraintType.HARD:
                    hard_count += 1
                else:
                    soft_count += 1

            return ConstraintCheckResult(
                passed=hard_count == 0,
                violations=violations,
                hard_violations=hard_count,
                soft_violations=soft_count,
                exempted_count=exempted_count,
            )

        except Exception as e:
            logger.error("Constraint check failed", error=str(e))
            return ConstraintCheckResult(passed=True, violations=[])

    async def request_exemption(
        self,
        chapter_id: str,
        constraint_id: str,
        reason: str,
        user_id: str,
    ) -> ExemptionRequest:
        """请求约束豁免"""
        logger.info(
            "Requesting exemption",
            chapter_id=chapter_id,
            constraint_id=constraint_id,
        )

        # 查找约束
        constraint = next(
            (c for c in self.default_constraints if c.id == constraint_id),
            None
        )

        if not constraint:
            raise ValueError(f"Constraint {constraint_id} not found")

        if constraint.type == ConstraintType.HARD:
            raise ValueError("Cannot exempt hard constraints")

        # 创建豁免请求
        exemption = ExemptionRequest(
            constraint_id=constraint_id,
            chapter_id=chapter_id,
            reason=reason,
            approved=True,  # 软约束自动批准
            approved_by=user_id,
        )

        # 存储豁免
        if chapter_id not in self.exemptions:
            self.exemptions[chapter_id] = []
        self.exemptions[chapter_id].append(exemption)

        return exemption

    def get_exemptions(self, chapter_id: str) -> list[ExemptionRequest]:
        """获取章节的豁免列表"""
        return self.exemptions.get(chapter_id, [])

    def revoke_exemption(self, chapter_id: str, constraint_id: str) -> bool:
        """撤销豁免"""
        if chapter_id not in self.exemptions:
            return False

        self.exemptions[chapter_id] = [
            e for e in self.exemptions[chapter_id]
            if e.constraint_id != constraint_id
        ]
        return True

    def get_all_constraints(self) -> list[Constraint]:
        """获取所有约束"""
        return self.default_constraints

    def add_project_constraint(
        self,
        project_id: str,
        name: str,
        description: str,
        rule: str,
        constraint_type: ConstraintType = ConstraintType.SOFT,
        category: ConstraintCategory = ConstraintCategory.PLOT,
    ) -> Constraint:
        """添加项目自定义约束"""
        constraint = Constraint(
            id=f"custom_{project_id}_{len(self.default_constraints)}",
            name=name,
            description=description,
            type=constraint_type,
            category=category,
            rule=rule,
        )
        # 注意：实际应用中应该存储到数据库
        return constraint
