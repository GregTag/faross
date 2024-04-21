import sys
import requests
import json


def getResponse(url: str, data: str) -> dict[str]:
    response = requests.post(url, data=data)

    if response.status_code == 200:
        sys.stderr.write(f"Запрос к {url} успешно выполнен!\n")
        return response.json()
    else:
        sys.stderr.write(f"Произошла ошибка {response.status_code}.\n")
        sys.stderr.write(response.json()["message"])
        sys.exit(1)


def calculateVulnsCount(vulns: list) -> int:
    count = 0
    seen = set()
    for v in vulns:
        if v["id"] not in seen:
            seen.add(v["id"])
            count += 1
        seen.update(v["aliases"])

    return count


purl = sys.argv[1]
base_url = "https://api.osv.dev/v1/query"

request_body = f'{{"package": {{"purl": "{purl}"}}}}'
response = getResponse(url=base_url, data=request_body)

# TODO: Есть небольшой нюанс с тем, что если по данному purl такого пакета не существует,
# TODO: то API выдает такой же ответ, что и при отсутствии уязвимостей -- пустой json ({})

vulns_count = calculateVulnsCount(response["vulns"]) if response else 0
score = 10 - vulns_count if vulns_count <= 10 else 0

report = [
    {
        "checkName": "CVE-check",
        "score": score,
        "risk": "High",
        "description": f"Detected {vulns_count} vulnerabilities.",
    }
]

json.dump(report, sys.stdout)