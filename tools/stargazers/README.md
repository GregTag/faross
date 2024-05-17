## Popularity

Питоновская обёртка принимает два обязательных аргумента
```
python3 script.py OWNER REPO
```
Скрипт делает апи-запрос к гитхабу, получает `stargazers_count`, переводит это число по логарифмической шкале в `score` от 0 до 10, и выводит json в согласованном формате. Если произошла ошибка (например, репозиторий не был найден), `score` выставляется равным `0` и в `stderr` логируется ошибка.

Пример:
```
> python3 script.py microsoft OSSGadget 100
> {'Popularity': {'score': 4, 'risk': 'low', 'desc': 'Scores popularity based on number of stargazers of repo on GitHub'}}
```