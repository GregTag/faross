## Пример запуска контейнера

```sh
docker build -t osv-dev .
docker run -it --rm osv-dev pkg:pypi/jinja2@2.4.1
```