import json
import requests
import sys
import dependencies


def get_json(url):
    r = requests.get(url)
    if r.status_code != 200:
        raise Exception(
            f"An error occured while fetching the data from toxic-repos\n"
            f"Response status code: {r.status_code}\nError: {r.text}\n"
        )
    return r.json()


def find_items(data, name):
    items = []
    for item in data:
        if (
            name.lower() in item["name"].lower()
            or name.lower() in item["commit_link"].lower()
        ):
            items.append(item)
    return items


def normalize_item(item):
    PROBLEM2RISK = {
        "ip_block": "Medium",
        "political_slogans": "Medium",
        "hostile_actions": "Critical",
        "malware": "Critical",
        "ddos": "High",
        "broken_assembly": "Medium",
    }

    problem = item["problem_type"]
    risk = PROBLEM2RISK[problem]
    description = item["description"]

    normalized_item = {
        "description": f"{problem}: {description}",
        "risk": risk,
    }

    return normalized_item


RISK2SCORE = {
    "Low": 10,
    "Medium": 8,
    "High": 5,
    "Critical": 0,
    "?": "?",
}


def is_higher(lhs, rhs):
    return RISK2SCORE[lhs] < RISK2SCORE[rhs]


if __name__ == "__main__":
    risk, description = "Low", "No problems found"
    try:
        url = "https://raw.githubusercontent.com/toxic-repos/toxic-repos/main/data/json/toxic-repos.json"
        name = sys.argv[1]

        data = get_json(url)
        dependency_names = dependencies.main()

        for name in dependency_names:
            items = find_items(data, name)
            for item in items:
                result_item = normalize_item(item)
                
                if is_higher(result_item["risk"], risk):
                    risk = result_item["risk"]
                    description = f"{result_item['description']}\n"

    except Exception as e:
        risk, description = "?", str(e)
        sys.exit(1)
    finally:
        report = {
            "checkName": "toxic-repos.check",
            "score": RISK2SCORE[risk],
            "risk": risk,
            "description": description,
        }

        json.dump(report, sys.stdout)
