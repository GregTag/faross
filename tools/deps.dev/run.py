import sys
import requests
from urllib.parse import quote


purl = sys.argv[1]
percent_encoded_purl = quote(purl, safe='')

# Только метод PurlLookup поддерживает запрос через Purl
base_purl_lookup_url = "https://api.deps.dev/v3alpha/purl/"
purl_lookup_response = requests.get(base_purl_lookup_url + percent_encoded_purl)

if purl_lookup_response.status_code == 200:
    print("Запрос PurlLookup успешно выполнен!")
    print("Ответ сервера:")
    print(purl_lookup_response.json())
else:
    print("Произошла ошибка при выполнении запроса:")
    print("Код ошибки:", purl_lookup_response.status_code)
    print("Текст ошибки:", purl_lookup_response.text)