from typing import List, Any, Tuple

import json
import sys
import os

STATIC_TOOLS = ["osv.dev", "toxic-repos", "packj-static", "ossgadget", "application-inspector-operations", 
"application-inspector-filetypes", "scorecard"]
DYNAMIC_TOOLS = ["packj-trace"]


def process_part(data: dict, tag: str) -> Tuple[List[Any], List[Any], List[Any]]:
    results = data.get(tag, {})
    scores = []
    risks = []
    scores_info = []
    tools = DYNAMIC_TOOLS if tag == "dynamic_analysis" else STATIC_TOOLS
    for tool in tools:
        exit_code = int(results.get(tool, {}).get("exit_code", 1))
        if exit_code != 0:
            continue
        res = results.get(tool, {}).get("result", {})
        if not isinstance(res, dict):
            continue
        score = res.get("score")
        if score is not None and (isinstance(score, int) or isinstance(score, float)):
            scores.append(score)
            risks.append(res.get("risk", "").lower())
            scores_info.append(res)
    return scores, risks, scores_info


def parse_json_file(input_file: str) -> Tuple[float, List[Any]]:
    weighted_scores: List[Tuple[int, Any]] = []  # по каждой взвешенной оценке хранит мета-информацию
    total_weights = []
    has_critical = False

    risk_weight = {"no risk": 0, "low": 1, "medium": 2, "high": 3, "critical": 4}

    with open(input_file, "r") as file:
        data = json.load(file)
        scores_static, risks_static, scores_info_static = process_part(data, "static_analysis")
        scores_dynamic, risks_dynamic, scores_info_dynamic = process_part(data, "dynamic_analysis")

    scores_static.extend(scores_dynamic)
    risks_static.extend(risks_dynamic)
    scores_info_static.extend(scores_info_dynamic)

    for score, risk, info in zip(scores_static, risks_static, scores_info_static):
        if score == 0:
            has_critical = True
        weight = risk_weight.get(risk.lower(), 1)
        weighted_score = score * weight
        weighted_scores.append((weighted_score, info))
        total_weights.append(weight)

    if weighted_scores:
        if has_critical:
            average_score = 0
        else:
            average_score = sum(weighted_score for weighted_score, _ in weighted_scores) / sum(total_weights)
        impactful_scores = [score_info for _, score_info in sorted(weighted_scores, key=lambda x: -x[0])][:5]
        return min(average_score, 10), impactful_scores

    return 0, []  # If all scores failed to parse, we do not approve


def calculate_decision(score, impactful_scores):
    decision: bool = score < 6
    return {"score": score, "is_quarantined": decision, "impactful_scores": impactful_scores}


if __name__ == "__main__":
    aggregate_score, impact_scores = parse_json_file("/usr/src/app/input/input.json")
    json.dump(calculate_decision(aggregate_score, impact_scores), sys.stdout)
