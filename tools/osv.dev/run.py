import sys
import requests
import json
import cvss


def getResponse(url: str, data: str) -> dict[str]:
    response = requests.post(url, data=data)

    if response.status_code == 200:
        return response.json()
    else:
        error_message = response.json()["message"]
        raise Exception(
            f"""An error occured while fetching data from {url} 
            Response status code: {response.status_code}
            Error: {error_message}"""
        )


def getDanderousVulnsCount(vulns: list) -> tuple[int, int]:
    CVSS_CALCULATOR = {
        "CVSS_V2": cvss.CVSS2,
        "CVSS_V3": cvss.CVSS3,
        "CVSS_V4": cvss.CVSS4,
    }

    high_cnt, critical_cnt = 0, 0
    for v in vulns:
        if v["id"].startswith("GHSA") and "severity" in v:
            for sev_object in v["severity"]:
                try:
                    calc = CVSS_CALCULATOR[sev_object["type"]]
                    cvss_vector = calc((sev_object["score"]))
                    base_severity = cvss_vector.severities()[0]
                    if base_severity == "High":
                        high_cnt += 1
                    if base_severity == "Critical":
                        critical_cnt += 1
                except KeyError:
                    continue
                except Exception:
                    raise Exception(
                        f"Unknown error in calculation of severity for {v['id']}"
                    )
                else:
                    break

    return high_cnt, critical_cnt


def getScore(count: int) -> int:
    return {
        count >= 100: 0,
        75 <= count < 100: 1,
        55 <= count < 75: 2,
        35 <= count < 55: 3,
        20 <= count < 35: 4,
        10 <= count < 20: 5,
        7 <= count < 10: 6,
        5 <= count < 7: 7,
        3 <= count < 5: 8,
        2 <= count < 3: 9,
        count < 2: 10,
    }[True]


if __name__ == "__main__":
    score = 10
    description = "No vulnerabilities have been detected or such PURL doesn't exist."

    try:
        purl = sys.argv[1]
        base_url = "https://api.osv.dev/v1/query"

        request_body = f'{{"package": {{"purl": "{purl}"}}}}'
        response = getResponse(url=base_url, data=request_body)

        if response:
            high_cnt, critical_cnt = getDanderousVulnsCount(response["vulns"])
            score = getScore(high_cnt + critical_cnt * 2)
            description = (
                f"Detected {high_cnt} vulnerability (-ies) with HIGH severity"
                f" and {critical_cnt} vulnerability (-ies) with CRITICAL severity."
            )
    except Exception as e:
        score = "?"
        description = str(e)
        sys.exit(1)
    finally:
        report = [
            {
                "checkName": "CVE-check",
                "score": score,
<<<<<<< HEAD
                "risk": "medium",
=======
                "risk": "Medium",
>>>>>>> 9f7f0e8 ([DAS-136] modified deps.dev script for toxic-repos)
                "description": description,
            }
        ]

        json.dump(report, sys.stdout)
