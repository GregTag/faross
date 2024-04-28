## Issues

Питоновская обёртка принимает два обязательных аргумента
```
python3 script.py OWNER REPO
```
Скрипт делает несколько апи-запросов к гитхабу, вычисляет `score` от 0 до 10 (=всплеска issue нет), и выводит json в согласованном формате. Если произошла ошибка (например, репозиторий не был найден), `score` выставляется равным `0` и в `stderr` логируется ошибка.

Пример:
```
> python3.9 script.py kamranahmedse developer-roadmap
> {'Issues': {'score': 10, 'risk': 'medium', 'desc': 'Scores issue splash based on increase in issue number over the last month'}}
```

#### Скрипт `api.sh`
Скрипт принимает 2 обязательных и 1 оптиональный параметр.
```
./run.sh OWNER REPO [DAYS_AGO]
```
По умолчанию `DAYS_AGO=30`
На выход скрипт отдаёт единственное число - количество всех issue, созданных за последние `DAYS_AGO`. Для своей работы скрипт создаёт промежуточный файл `temp.txt`, который потом сам и удаляет.

Пример:
```
> ./run.sh microsoft OSSGadget 100
> 2
```