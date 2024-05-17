## Пример запуска контейнера

```
docker build -t toxic-repos .
docker run -it --rm toxic-repos pkg:npm/es5-ext@0.10.64
```

или

```
docker build -t toxic-repos .
docker run -it --rm toxic-repos es5-ext
```

