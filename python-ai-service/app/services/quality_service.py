"""
质量评估服务 - 追读力分析、一致性检查、多Agent审查
"""
import asyncio
import json
import re
from dataclasses import dataclass, field, asdict
from typing import Optional, Any

import structlog

from app.services.llm_service import LLMService

logger = structlog.get_logger()


@dataclass
class HookAnalysis:
    """钩子分析"""
    hook_type: str  # crisis, mystery, emotion, choice, desire
    hook_strength: str  # strong, medium, weak
    hook_content: str
    score: float
    suggestions: list[str] = field(default_factory=list)


@dataclass
class CoolpointItem:
    """爽点项"""
    type: str  # achievement, revenge, romance, power_up, recognition
    content: str
    position: int  # 在文本中的位置百分比
    intensity: float  # 0-10


@dataclass
class MicropayoffItem:
    """微兑现项"""
    goal: str  # 设定的目标/期望
    payoff: str  # 兑现内容
    position: int  # 位置百分比
    satisfaction: float  # 满足度 0-10


@dataclass
class ConsistencyIssue:
    """一致性问题"""
    type: str  # character, timeline, worldview, logic
    severity: str  # critical, major, minor
    description: str
    chapter_refs: list[int] = field(default_factory=list)
    suggestion: str = ""


@dataclass
class AgentReview:
    """单个Agent审查结果"""
    agent_name: str
    role: str  # reader, editor, critic, writer, marketer, psychologist
    score: float
    strengths: list[str] = field(default_factory=list)
    weaknesses: list[str] = field(default_factory=list)
    suggestions: list[str] = field(default_factory=list)


@dataclass
class ReadingPowerResult:
    """追读力分析结果"""
    overall_score: float
    hook: Optional[HookAnalysis] = None
    coolpoints: list[CoolpointItem] = field(default_factory=list)
    micropayoffs: list[MicropayoffItem] = field(default_factory=list)
    dimension_scores: dict[str, float] = field(default_factory=dict)
    suggestions: list[str] = field(default_factory=list)


@dataclass
class ConsistencyResult:
    """一致性检查结果"""
    overall_score: float
    issues: list[ConsistencyIssue] = field(default_factory=list)
    character_score: float = 0.0
    timeline_score: float = 0.0
    worldview_score: float = 0.0


@dataclass
class MultiAgentResult:
    """多Agent审查结果"""
    overall_score: float
    reviews: list[AgentReview] = field(default_factory=list)
    consensus_strengths: list[str] = field(default_factory=list)
    consensus_weaknesses: list[str] = field(default_factory=list)
    final_suggestions: list[str] = field(default_factory=list)


class QualityService:
    """质量评估服务"""

    def __init__(self):
        self.llm_service = LLMService()

    def _parse_json_response(self, text: str) -> dict:
        """解析 LLM 返回的 JSON"""
        # 尝试提取 JSON 块
        json_match = re.search(r'```json\s*([\s\S]*?)\s*```', text)
        if json_match:
            text = json_match.group(1)
        else:
            # 尝试找到 { } 包围的内容
            json_match = re.search(r'\{[\s\S]*\}', text)
            if json_match:
                text = json_match.group(0)

        try:
            return json.loads(text)
        except json.JSONDecodeError:
            logger.warning("Failed to parse JSON response", text=text[:200])
            return {}

    # ==================== T7.1 追读力分析 ====================

    async def analyze_reading_power(
        self,
        chapter_id: str,
        content: str,
        provider: str = "openai",
        model: str = "gpt-4o-mini",
    ) -> ReadingPowerResult:
        """分析章节追读力"""
        logger.info(
            "Analyzing reading power",
            chapter_id=chapter_id,
            content_length=len(content),
        )

        # 并行执行三种分析
        hook_task = self._analyze_hook(content, provider, model)
        coolpoint_task = self._analyze_coolpoints(content, provider, model)
        micropayoff_task = self._analyze_micropayoffs(content, provider, model)

        hook, coolpoints, micropayoffs = await asyncio.gather(
            hook_task, coolpoint_task, micropayoff_task
        )

        # 计算各维度评分
        dimension_scores = {
            "hook": hook.score if hook else 5.0,
            "coolpoint": self._calc_coolpoint_score(coolpoints),
            "micropayoff": self._calc_micropayoff_score(micropayoffs),
        }

        # 计算综合评分 (钩子30% + 爽点40% + 微兑现30%)
        overall_score = (
            dimension_scores["hook"] * 0.3 +
            dimension_scores["coolpoint"] * 0.4 +
            dimension_scores["micropayoff"] * 0.3
        )

        # 生成建议
        suggestions = self._generate_reading_power_suggestions(
            hook, coolpoints, micropayoffs, dimension_scores
        )

        return ReadingPowerResult(
            overall_score=round(overall_score, 1),
            hook=hook,
            coolpoints=coolpoints,
            micropayoffs=micropayoffs,
            dimension_scores=dimension_scores,
            suggestions=suggestions,
        )

    async def _analyze_hook(
        self, content: str, provider: str, model: str
    ) -> Optional[HookAnalysis]:
        """分析章节钩子"""
        # 获取章节开头 (前 800 字)
        opening = content[:800] if len(content) > 800 else content

        prompt = f"""你是一位网文分析专家。请分析以下章节开头的"钩子"效果。

钩子类型说明:
- crisis: 危机钩子 - 主角面临紧迫危险或困境
- mystery: 悬念钩子 - 设置谜团或未解之谜
- emotion: 情感钩子 - 强烈的情感冲突或共鸣
- choice: 抉择钩子 - 主角面临重大选择
- desire: 欲望钩子 - 激发读者对某事物的渴望

章节开头:
---
{opening}
---

请以 JSON 格式返回分析结果:
```json
{{
    "hook_type": "类型(crisis/mystery/emotion/choice/desire)",
    "hook_strength": "强度(strong/medium/weak)",
    "hook_content": "钩子内容摘录(50字以内)",
    "score": 评分(0-10),
    "suggestions": ["改进建议1", "改进建议2"]
}}
```"""

        try:
            response = await self.llm_service.generate(
                prompt=prompt,
                provider=provider,
                model=model,
                temperature=0.3,
                max_tokens=500,
            )

            data = self._parse_json_response(response)
            if data:
                return HookAnalysis(
                    hook_type=data.get("hook_type", "mystery"),
                    hook_strength=data.get("hook_strength", "medium"),
                    hook_content=data.get("hook_content", opening[:50]),
                    score=float(data.get("score", 5.0)),
                    suggestions=data.get("suggestions", []),
                )
        except Exception as e:
            logger.error("Hook analysis failed", error=str(e))

        # 返回默认值
        return HookAnalysis(
            hook_type="mystery",
            hook_strength="medium",
            hook_content=opening[:50],
            score=5.0,
            suggestions=["建议加强开头的悬念设置"],
        )

    async def _analyze_coolpoints(
        self, content: str, provider: str, model: str
    ) -> list[CoolpointItem]:
        """分析爽点分布"""
        prompt = f"""你是一位网文分析专家。请分析以下章节中的"爽点"分布。

爽点类型说明:
- achievement: 成就爽点 - 主角达成目标、获得成功
- revenge: 复仇爽点 - 打脸、报复、以牙还牙
- romance: 情感爽点 - 感情进展、甜蜜互动
- power_up: 升级爽点 - 实力提升、获得宝物
- recognition: 认可爽点 - 被认可、被崇拜、扬名

章节内容:
---
{content[:3000]}
---

请以 JSON 格式返回分析结果(最多5个爽点):
```json
{{
    "coolpoints": [
        {{
            "type": "类型",
            "content": "爽点内容摘录(30字以内)",
            "position": 位置百分比(0-100),
            "intensity": 强度(0-10)
        }}
    ]
}}
```"""

        try:
            response = await self.llm_service.generate(
                prompt=prompt,
                provider=provider,
                model=model,
                temperature=0.3,
                max_tokens=800,
            )

            data = self._parse_json_response(response)
            if data and "coolpoints" in data:
                return [
                    CoolpointItem(
                        type=cp.get("type", "achievement"),
                        content=cp.get("content", ""),
                        position=int(cp.get("position", 50)),
                        intensity=float(cp.get("intensity", 5.0)),
                    )
                    for cp in data["coolpoints"][:5]
                ]
        except Exception as e:
            logger.error("Coolpoint analysis failed", error=str(e))

        return []

    async def _analyze_micropayoffs(
        self, content: str, provider: str, model: str
    ) -> list[MicropayoffItem]:
        """分析微兑现"""
        prompt = f"""你是一位网文分析专家。请分析以下章节中的"微兑现"。

微兑现是指: 在章节中设置的小目标、小期待得到满足的情节。
例如: 读者期待主角打败某个对手 → 主角成功击败对手

章节内容:
---
{content[:3000]}
---

请以 JSON 格式返回分析结果(最多5个微兑现):
```json
{{
    "micropayoffs": [
        {{
            "goal": "设定的目标/期望(20字以内)",
            "payoff": "兑现内容(30字以内)",
            "position": 位置百分比(0-100),
            "satisfaction": 满足度(0-10)
        }}
    ]
}}
```"""

        try:
            response = await self.llm_service.generate(
                prompt=prompt,
                provider=provider,
                model=model,
                temperature=0.3,
                max_tokens=800,
            )

            data = self._parse_json_response(response)
            if data and "micropayoffs" in data:
                return [
                    MicropayoffItem(
                        goal=mp.get("goal", ""),
                        payoff=mp.get("payoff", ""),
                        position=int(mp.get("position", 50)),
                        satisfaction=float(mp.get("satisfaction", 5.0)),
                    )
                    for mp in data["micropayoffs"][:5]
                ]
        except Exception as e:
            logger.error("Micropayoff analysis failed", error=str(e))

        return []

    def _calc_coolpoint_score(self, coolpoints: list[CoolpointItem]) -> float:
        """计算爽点评分"""
        if not coolpoints:
            return 4.0  # 无爽点给低分

        # 基于数量和强度计算
        avg_intensity = sum(cp.intensity for cp in coolpoints) / len(coolpoints)
        count_bonus = min(len(coolpoints) * 0.5, 2.0)  # 数量加成，最多+2

        return min(avg_intensity + count_bonus, 10.0)

    def _calc_micropayoff_score(self, micropayoffs: list[MicropayoffItem]) -> float:
        """计算微兑现评分"""
        if not micropayoffs:
            return 4.0

        avg_satisfaction = sum(mp.satisfaction for mp in micropayoffs) / len(micropayoffs)
        count_bonus = min(len(micropayoffs) * 0.5, 2.0)

        return min(avg_satisfaction + count_bonus, 10.0)

    def _generate_reading_power_suggestions(
        self,
        hook: Optional[HookAnalysis],
        coolpoints: list[CoolpointItem],
        micropayoffs: list[MicropayoffItem],
        scores: dict[str, float],
    ) -> list[str]:
        """生成追读力改进建议"""
        suggestions = []

        # 钩子建议
        if scores.get("hook", 0) < 6:
            suggestions.append("章节开头缺乏吸引力，建议增加悬念或冲突")
        if hook and hook.suggestions:
            suggestions.extend(hook.suggestions[:2])

        # 爽点建议
        if scores.get("coolpoint", 0) < 6:
            suggestions.append("爽点不足，建议增加主角的成就感或打脸情节")
        if len(coolpoints) < 2:
            suggestions.append("爽点数量偏少，建议每章至少设置2-3个爽点")

        # 微兑现建议
        if scores.get("micropayoff", 0) < 6:
            suggestions.append("微兑现不足，建议设置更多小目标并及时满足")
        if not micropayoffs:
            suggestions.append("缺少微兑现，读者可能感到阅读疲劳")

        return suggestions[:5]  # 最多返回5条建议

    # ==================== T7.2 一致性检查 ====================

    async def check_consistency(
        self,
        chapter_id: str,
        content: str,
        context: dict[str, Any],  # 包含角色、世界观等上下文
        provider: str = "openai",
        model: str = "gpt-4o-mini",
    ) -> ConsistencyResult:
        """检查章节一致性"""
        logger.info(
            "Checking consistency",
            chapter_id=chapter_id,
            content_length=len(content),
        )

        # 并行执行三种一致性检查
        char_task = self._check_character_consistency(content, context, provider, model)
        timeline_task = self._check_timeline_consistency(content, context, provider, model)
        worldview_task = self._check_worldview_consistency(content, context, provider, model)

        char_result, timeline_result, worldview_result = await asyncio.gather(
            char_task, timeline_task, worldview_task
        )

        # 合并问题
        all_issues = char_result[1] + timeline_result[1] + worldview_result[1]

        # 计算综合评分
        overall_score = (char_result[0] + timeline_result[0] + worldview_result[0]) / 3

        return ConsistencyResult(
            overall_score=round(overall_score, 1),
            issues=all_issues,
            character_score=char_result[0],
            timeline_score=timeline_result[0],
            worldview_score=worldview_result[0],
        )

    async def _check_character_consistency(
        self, content: str, context: dict, provider: str, model: str
    ) -> tuple[float, list[ConsistencyIssue]]:
        """检查角色一致性"""
        characters = context.get("characters", [])
        if not characters:
            return 10.0, []

        char_info = "\n".join([
            f"- {c.get('name')}: {c.get('personality', '无')} / {c.get('background', '无')}"
            for c in characters[:10]
        ])

        prompt = f"""你是一位小说编辑。请检查以下章节中角色表现是否与设定一致。

角色设定:
{char_info}

章节内容:
---
{content[:2500]}
---

请检查:
1. 角色性格是否与设定一致
2. 角色行为是否符合其背景
3. 角色对话风格是否统一

以 JSON 格式返回:
```json
{{
    "score": 一致性评分(0-10),
    "issues": [
        {{
            "type": "character",
            "severity": "critical/major/minor",
            "description": "问题描述",
            "suggestion": "修改建议"
        }}
    ]
}}
```"""

        try:
            response = await self.llm_service.generate(
                prompt=prompt, provider=provider, model=model,
                temperature=0.3, max_tokens=800,
            )
            data = self._parse_json_response(response)
            if data:
                issues = [
                    ConsistencyIssue(
                        type="character",
                        severity=i.get("severity", "minor"),
                        description=i.get("description", ""),
                        suggestion=i.get("suggestion", ""),
                    )
                    for i in data.get("issues", [])
                ]
                return float(data.get("score", 8.0)), issues
        except Exception as e:
            logger.error("Character consistency check failed", error=str(e))

        return 8.0, []

    async def _check_timeline_consistency(
        self, content: str, context: dict, provider: str, model: str
    ) -> tuple[float, list[ConsistencyIssue]]:
        """检查时间线一致性"""
        timeline = context.get("timeline", [])
        previous_summary = context.get("previous_summary", "")

        prompt = f"""你是一位小说编辑。请检查以下章节的时间线是否合理。

前情提要:
{previous_summary[:500] if previous_summary else "无"}

章节内容:
---
{content[:2500]}
---

请检查:
1. 时间顺序是否合理
2. 事件发生的时间间隔是否合理
3. 是否有时间矛盾

以 JSON 格式返回:
```json
{{
    "score": 一致性评分(0-10),
    "issues": [
        {{
            "type": "timeline",
            "severity": "critical/major/minor",
            "description": "问题描述",
            "suggestion": "修改建议"
        }}
    ]
}}
```"""

        try:
            response = await self.llm_service.generate(
                prompt=prompt, provider=provider, model=model,
                temperature=0.3, max_tokens=800,
            )
            data = self._parse_json_response(response)
            if data:
                issues = [
                    ConsistencyIssue(
                        type="timeline",
                        severity=i.get("severity", "minor"),
                        description=i.get("description", ""),
                        suggestion=i.get("suggestion", ""),
                    )
                    for i in data.get("issues", [])
                ]
                return float(data.get("score", 8.0)), issues
        except Exception as e:
            logger.error("Timeline consistency check failed", error=str(e))

        return 8.0, []

    async def _check_worldview_consistency(
        self, content: str, context: dict, provider: str, model: str
    ) -> tuple[float, list[ConsistencyIssue]]:
        """检查世界观一致性"""
        world_settings = context.get("world_settings", [])
        if not world_settings:
            return 10.0, []

        settings_info = "\n".join([
            f"- {s.get('title')}: {s.get('content', '')[:100]}"
            for s in world_settings[:10]
        ])

        prompt = f"""你是一位小说编辑。请检查以下章节是否与世界观设定一致。

世界观设定:
{settings_info}

章节内容:
---
{content[:2500]}
---

请检查:
1. 力量体系是否符合设定
2. 社会规则是否一致
3. 地理/环境描述是否准确

以 JSON 格式返回:
```json
{{
    "score": 一致性评分(0-10),
    "issues": [
        {{
            "type": "worldview",
            "severity": "critical/major/minor",
            "description": "问题描述",
            "suggestion": "修改建议"
        }}
    ]
}}
```"""

        try:
            response = await self.llm_service.generate(
                prompt=prompt, provider=provider, model=model,
                temperature=0.3, max_tokens=800,
            )
            data = self._parse_json_response(response)
            if data:
                issues = [
                    ConsistencyIssue(
                        type="worldview",
                        severity=i.get("severity", "minor"),
                        description=i.get("description", ""),
                        suggestion=i.get("suggestion", ""),
                    )
                    for i in data.get("issues", [])
                ]
                return float(data.get("score", 8.0)), issues
        except Exception as e:
            logger.error("Worldview consistency check failed", error=str(e))

        return 8.0, []

    # ==================== T7.3 多Agent审查 ====================

    async def multi_agent_review(
        self,
        chapter_id: str,
        content: str,
        provider: str = "openai",
        model: str = "gpt-4o-mini",
    ) -> MultiAgentResult:
        """多Agent并行审查"""
        logger.info(
            "Starting multi-agent review",
            chapter_id=chapter_id,
            content_length=len(content),
        )

        # 定义6个Agent角色
        agents = [
            ("reader", "普通读者", "从读者角度评价阅读体验、代入感、情感共鸣"),
            ("editor", "资深编辑", "从编辑角度评价文笔、节奏、结构"),
            ("critic", "文学评论家", "从文学角度评价主题深度、人物塑造、叙事技巧"),
            ("writer", "网文作者", "从同行角度评价商业性、追读力、爽点设置"),
            ("marketer", "市场分析师", "从市场角度评价受众定位、卖点、竞争力"),
            ("psychologist", "读者心理专家", "从心理角度评价情绪曲线、满足感、成瘾性"),
        ]

        # 并行执行6个Agent审查
        tasks = [
            self._single_agent_review(content, agent, provider, model)
            for agent in agents
        ]
        reviews = await asyncio.gather(*tasks)

        # 汇总结果
        overall_score = sum(r.score for r in reviews) / len(reviews)

        # 找出共识
        all_strengths = [s for r in reviews for s in r.strengths]
        all_weaknesses = [w for r in reviews for w in r.weaknesses]
        all_suggestions = [s for r in reviews for s in r.suggestions]

        # 统计出现频率高的项目
        consensus_strengths = self._find_consensus(all_strengths, threshold=2)
        consensus_weaknesses = self._find_consensus(all_weaknesses, threshold=2)
        final_suggestions = self._find_consensus(all_suggestions, threshold=2)

        return MultiAgentResult(
            overall_score=round(overall_score, 1),
            reviews=reviews,
            consensus_strengths=consensus_strengths[:5],
            consensus_weaknesses=consensus_weaknesses[:5],
            final_suggestions=final_suggestions[:5],
        )

    async def _single_agent_review(
        self,
        content: str,
        agent: tuple[str, str, str],
        provider: str,
        model: str,
    ) -> AgentReview:
        """单个Agent审查"""
        role_id, role_name, role_desc = agent

        prompt = f"""你是一位{role_name}。{role_desc}

请审查以下章节内容:
---
{content[:3000]}
---

以 JSON 格式返回你的审查意见:
```json
{{
    "score": 综合评分(0-10),
    "strengths": ["优点1", "优点2", "优点3"],
    "weaknesses": ["不足1", "不足2"],
    "suggestions": ["建议1", "建议2"]
}}
```"""

        try:
            response = await self.llm_service.generate(
                prompt=prompt, provider=provider, model=model,
                temperature=0.5, max_tokens=600,
            )
            data = self._parse_json_response(response)
            if data:
                return AgentReview(
                    agent_name=role_name,
                    role=role_id,
                    score=float(data.get("score", 7.0)),
                    strengths=data.get("strengths", [])[:3],
                    weaknesses=data.get("weaknesses", [])[:3],
                    suggestions=data.get("suggestions", [])[:3],
                )
        except Exception as e:
            logger.error(f"Agent {role_name} review failed", error=str(e))

        return AgentReview(
            agent_name=role_name,
            role=role_id,
            score=7.0,
            strengths=["内容完整"],
            weaknesses=[],
            suggestions=[],
        )

    def _find_consensus(self, items: list[str], threshold: int = 2) -> list[str]:
        """找出共识项目（出现次数>=threshold的项目）"""
        from collections import Counter

        # 简单的相似度匹配（基于关键词）
        normalized = []
        for item in items:
            # 提取关键词
            keywords = set(item.replace("，", " ").replace("、", " ").split())
            normalized.append((item, keywords))

        # 统计相似项
        result = []
        used = set()

        for i, (item1, kw1) in enumerate(normalized):
            if i in used:
                continue
            count = 1
            for j, (item2, kw2) in enumerate(normalized[i+1:], i+1):
                if j in used:
                    continue
                # 如果关键词重叠超过50%，认为是相似的
                overlap = len(kw1 & kw2) / max(len(kw1), len(kw2), 1)
                if overlap > 0.3:
                    count += 1
                    used.add(j)

            if count >= threshold:
                result.append(item1)
                used.add(i)

        return result

    # ==================== 综合评估接口 ====================

    async def evaluate_quality(
        self,
        chapter_id: str,
        content: str,
        context: dict[str, Any],
        provider: str = "openai",
        model: str = "gpt-4o-mini",
    ) -> dict:
        """综合质量评估"""
        logger.info(
            "Evaluating quality",
            chapter_id=chapter_id,
        )

        # 并行执行所有评估
        reading_power_task = self.analyze_reading_power(
            chapter_id, content, provider, model
        )
        consistency_task = self.check_consistency(
            chapter_id, content, context, provider, model
        )
        multi_agent_task = self.multi_agent_review(
            chapter_id, content, provider, model
        )

        reading_power, consistency, multi_agent = await asyncio.gather(
            reading_power_task, consistency_task, multi_agent_task
        )

        # 计算综合评分
        overall_score = (
            reading_power.overall_score * 0.4 +
            consistency.overall_score * 0.3 +
            multi_agent.overall_score * 0.3
        )

        return {
            "overall_score": round(overall_score, 1),
            "reading_power": asdict(reading_power),
            "consistency": asdict(consistency),
            "multi_agent": asdict(multi_agent),
            "summary": {
                "strengths": multi_agent.consensus_strengths,
                "weaknesses": multi_agent.consensus_weaknesses,
                "suggestions": (
                    reading_power.suggestions[:2] +
                    multi_agent.final_suggestions[:3]
                ),
            },
        }
