import sys
import requests
import json
from urllib.parse import quote


def getPathParameters(url: str) -> tuple[str, str, str]:
    response = requests.get(url)

    if response.status_code == 200:
        version_key = response.json()["version"]["versionKey"]
        return version_key["system"], version_key["name"], version_key["version"]
    else:
        error_message = response.text
        raise Exception(
            f"""An error occured while fetching data from {url} 
            Response status code: {response.status_code}
            Error: {error_message}"""
        )


def checkErrors(error_message: str) -> None:
    if error_message:
        raise Exception(
            f"""An internal server error associated with the dependency graph
            An error here may imply the graph as a whole is incorrect 
            Error: {error_message}"""
        )


def getDependencies(url: str) -> list[str]:
    response = requests.get(url)

    if response.status_code == 200:
        dependency_names = []
        checkErrors(response.json()["error"])
        nodes = response.json()["nodes"]
        for v in nodes:
            try:
                checkErrors(v["errors"])
            except Exception:
                continue
            version_key = v["versionKey"]
            dependency_names.append(version_key["name"])
        return dependency_names
    else:
        error_message = response.text
        raise Exception(
            f"""An error occured while fetching data from {url} 
            Response status code: {response.status_code}
            Error: {error_message}"""
        )


if __name__ == "__main__":
    try:
        purl = sys.argv[1]
        base_url = "https://api.deps.dev"

        purl_lookup_url = base_url + "/v3alpha/purl/" + quote(purl, safe="")
        system, name, version = getPathParameters(purl_lookup_url)

        encoded_name = quote(name, safe="")
        dependencies_url = (
            base_url
            + f"/v3/systems/{system}/packages/{encoded_name}/versions/{version}:dependencies"
        )

        dependency_names = getDependencies(dependencies_url)
        json.dump(dependency_names, sys.stdout)

    except Exception as e:
        sys.stderr.write(str(e))
        sys.exit(1)
