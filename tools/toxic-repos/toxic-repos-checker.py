import json
import os
import requests
import sys


def get_json(url, filename):
    r = requests.get(url)
    if r.status_code != 200:
        sys.stderr.write(f'An error occured while fetching the data from toxic-repos\n'
                         f'Response status code: {r.status_code}\nError: {r.text}\n')
        sys.exit(1)
    return r.json()


def find_item(data, name):
    for item in data:
        if name.lower() in item['name'].lower() or name.lower() in \
        item['commit_link'].lower():
            return item
    return {}


def normalize_item(item, name):
    inaccessible_category = 'INACCESSIBLE'
    warning_category = 'WARNING'
    dangerous_category = 'DANGEROUS'
    safe_category = 'SAFE'

    problem_categories = {
        'ip_block': inaccessible_category,
        'political_slogans': warning_category,
        'hostile_actions': dangerous_category,
        'malware': dangerous_category,
        'ddos': dangerous_category,
        'broken_assembly': dangerous_category,
        'none': safe_category
    }

    risks = {
        safe_category: 'Low',
        inaccessible_category: 'Medium',
        warning_category: 'Medium',
        dangerous_category: 'High',
    }
    
    scores = {
        safe_category: 10,
        inaccessible_category: 5,
        warning_category: 5,
        dangerous_category: 2,
    }

    problem_category = problem_categories[item.get('problem_type', 'none')]

    normalized_item = {
        'name': name,
        'description': item.get('description', ''),
        'problem_type': item.get('problem_type', ''),
        'problem_category': problem_categories[item.get('problem_type',
                                                        'none')],
        'risk': risks[problem_category],
        'score': scores[problem_category]
    }

    return normalized_item


if __name__ == '__main__':
    url = 'https://raw.githubusercontent.com/toxic-repos/toxic-repos/main/data/json/toxic-repos.json'
    filename = 'toxic-repos.json'
    name = sys.argv[1]

    data = get_json(url, filename)
    item = find_item(data, name)
    result_item = normalize_item(item, name)

    report = {
        "checkName": "toxic-repos.check",
        "score": result_item['score'],
        "risk": result_item['risk'],
        "description": result_item['problem_type']
    }

    json.dump(report, sys.stdout)
