import sys
import requests
from urllib.parse import quote


def getResponse(url: str) -> None:
    response = requests.get(url)

    if response.status_code == 200:
        sys.stderr.write(f"Запрос к {url} успешно выполнен!")
        print(response.json())
    else:
        sys.stderr.write(f"Произошла ошибка, {response.status_code}\n")
        sys.stderr.write(response.text)
        sys.exit(1)


purl = sys.argv[1]
percent_encoded_purl = quote(purl, safe='')

# Только метод PurlLookup поддерживает запрос через Purl
base_purl_lookup_url = "https://api.deps.dev/v3alpha/purl/"
# TODO: JSON не соответсвует образцу
# TODO: используется alpha версия API

url = base_purl_lookup_url + percent_encoded_purl
getResponse(url)