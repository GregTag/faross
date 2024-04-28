## Всплески issue

Скрипт принимает 2 обязательных и 1 оптиональный параметр.
```
./run.sh OWNER REPO [DAYS_AGO]
```
По умолчанию `DAYS_AGO=30`
На выход скрипт отдаёт единственное число - количество открытых issue, созданных за последние `DAYS_AGO`. Для своей работы скрипт создаёт файл `issues.txt`, куда записывает уже отфильтрованные (без pull-requests) issues.

Пример:
```
> ./run.sh microsoft OSSGadget 100
> 2
```