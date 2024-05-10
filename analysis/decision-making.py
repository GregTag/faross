from typing import List, Any, Tuple

import json
import sys
import os

STATIC_TOOLS = ['osv.dev', 'toxic-repos', 'packj-static']
DYNAMIC_TOOLS = ['packj-trace']


def process_part(data: dict, tag: str) -> Tuple[List[Any], List[Any]]:
    results = data.get(tag, None)
    if not results:
        return [], []
    scores = []
    risks = []
    tools = DYNAMIC_TOOLS if tag == 'dynamic_analysis' else STATIC_TOOLS
    for tool in tools:
        res = results.get(tool).get('result')
        if not res:
            continue
        scores.append(res.get('score'))
        risks.append(res.get('risk').lower())
    return scores, risks


def parse_json_file(input_file: str):
    weighted_scores = []
    total_weights = []
    has_critical = False

    risk_weight = {'no risk': 0, 'low': 1, 'medium': 2, 'high': 3, 'critical': 4}
    
    with open(input_file, 'r') as file:
        data = json.load(file)
        scores_static, risks_static = process_part(data, 'static_analysis')
        scores_dynamic, risks_dynamic = process_part(data, 'dynamic_analysis')

    # TODO: replace with proper function
    scores_static.extend(scores_dynamic)
    risks_static.extend(risks_dynamic)

    for score, risk in zip(scores_static, risks_static):
        if risk.lower() == 'critical' and score > 0:
            has_critical = True
            break
        weight = risk_weight.get(risk.lower(), 1)
        weighted_score = score * weight
        weighted_scores.append(weighted_score)
        total_weights.append(weight)

    if has_critical:
        return 0

    if weighted_scores:
        average_score = sum(weighted_scores) / sum(total_weights)
        normalized_score = (average_score / max(total_weights)) * 10
        return min(normalized_score, 10)

    return None


def save_results(aggregate_score, output_path):
    decision = "yes" if aggregate_score > 6 else "no"
    results = {
        "aggregate-score": aggregate_score,
        "decision": decision
    }
    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    with open(output_path, 'w') as file:
        json.dump(results, file, indent=4)


if __name__ == "__main__":
    aggregate_score = parse_json_file("/usr/src/app/input/input.json")
    if aggregate_score is not None:
        json.dump({
            "score": aggregate_score
        }, sys.stdout)
        # save_results(aggregate_score, output_path)
    else:
        json.dump({
            "score": 10
        }, sys.stdout)
