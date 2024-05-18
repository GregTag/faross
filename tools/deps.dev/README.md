## Пример запуска контейнера

```sh
docker build -t depsdev .
docker run -it --rm depsdev pkg:golang/github.com/gin-gonic/gin@v1.10.0
```
