## Примеры запуска контейнера
```sh
docker build -t scorecard .
docker run -it --rm scorecard pkg:pypi/django@1.11.1
```
Не работает, если исходный репозиторий не на гитхабе