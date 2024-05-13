## Пример запуска контейнера

```sh
docker build -t depsdev .
docker run -it --rm depsdev pkg:npm/%40colors/colors@1.5.0
```
