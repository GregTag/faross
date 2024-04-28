import sys
import requests
from urllib.parse import quote


def getResponse(url: str) -> None:
    response = requests.get(url)

    if response.status_code == 200:
        print(response.text)
    else:
        sys.stderr.write(f"An error occured while fetching data from api.deps.dev\n"
                         f"Response status code: {response.status_code}\nError:\n{response.text}\n")
        sys.exit(1)


if __name__ == "__main__":
    purl = sys.argv[1]
    percent_encoded_purl = quote(purl, safe='')

    # Только метод PurlLookup поддерживает запрос через Purl
    base_purl_lookup_url = "https://api.deps.dev/v3alpha/purl/"
    # TODO: JSON не соответсвует образцу
    # TODO: используется alpha версия API

    url = base_purl_lookup_url + percent_encoded_purl
    getResponse(url)
