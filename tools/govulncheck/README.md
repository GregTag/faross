## Пример запуска контейнера

```sh
docker build -t govulncheck .
docker run -it --rm govulncheck | jq ".finding"
```