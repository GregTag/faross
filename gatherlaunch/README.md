## Тестовый вариант gather-launch, не требующий обертки для инструментов

### Warning:

Так как это демонстрационный вариант, в данный момент эта реализация умеет работать __только__ с тестовым контейнером, который просто выводит параметры, данные ему на вход.
Не подлежит контейнеризации (или нужно установить docker внутрь docker-образа с приложением)

### Инструкция по запуску:

1. `docker pull imarenf/example-output:1.0`

2. `go run cmd/main.go {purl}`

### Пример вывода:

```
{
        "static analysis": { 
                "osv.dev": {
                        "result": {"checkName": "CVE-check", "score": 9, "risk": "High", "description": "Detected 1 vulnerabilities."},
                        "exit_code": 0
                },
                "toxic-repos": {
                        "result": {"checkName": "toxic-repos.check", "score": 10, "risk": "Low", "description": ""},
                        "exit_code": 0
                },
                "deps.dev": {
                        "result": "{\"version\":{\"versionKey\":{\"system\":\"PYPI\",\"name\":\"requests\",\"version\":\"2.28.0\"},\"purl\":\"pkg:pypi/requests@2.28.0\",\"publishedAt\":\"2022-06-09T14:44:34Z\",\"isDefault\":false,\"licenses\":[\"Apache-2.0\"],\"advisoryKeys\":[{\"id\":\"GHSA-j8r2-6x86-q33q\"},{\"id\":\"PYSEC-2023-74\"}],\"links\":[{\"label\":\"DOCUMENTATION\",\"url\":\"https://requests.readthedocs.io\"},{\"label\":\"SOURCE_REPO\",\"url\":\"https://github.com/psf/requests\"}],\"slsaProvenances\":[],\"registries\":[\"https://pypi.org/simple\"],\"relatedProjects\":[{\"projectKey\":{\"id\":\"github.com/psf/requests\"},\"relationProvenance\":\"UNVERIFIED_METADATA\",\"relationType\":\"SOURCE_REPO\"}]}}",
                        "exit_code": 0
                }
        },
        "dynamic analysis": {
                "packj-trace": {
                        "result": "",
                        "exit-code": 0
                }
        }
}
```
