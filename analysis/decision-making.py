import json
import os

def parse_json_files(directory):
    weighted_scores = []
    total_weights = []
    has_critical = False

    risk_weight = {'low': 1, 'medium': 2, 'high': 3, 'critical': 4}

    for filename in os.listdir(directory):
        if filename.endswith('.json'):
            with open(os.path.join(directory, filename), 'r') as file:
                data = json.load(file)
                score = data.get('score')
                risk = data.get('risk')

                if risk == 'critical' and score > 0:
                    has_critical = True
                    break

                weight = risk_weight.get(risk, 1)
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
    directory = "input"
    output_path = "output/output.json"
    aggregate_score = parse_json_files(directory)
    if aggregate_score is not None:
        print(f"Aggregate score: {aggregate_score}")
        save_results(aggregate_score, output_path)
    else:
        print("No valid scores found to aggregate.")
