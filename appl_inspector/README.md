# Application Inspector
Выполянет две проверки: **Unsafe operations** и **File types**

Скрипт `run.sh` принимает на вход путь до директории, которую необходимо проверить, запускает application inspector, записывает его вывод в файл `data.json`. Запускает по очереди все проверки, выводит json в согласованном формате. Если произошла ошибка (например, директория не была найдена), скрипт завершает работу с `error_code=1`, ничего не выводит в поток вывода, логирует ошибку в `stderr`.

## Unsafe operations 
Из всей информации в data.json берёт только теги, ищет среди них опасные, присваивает пакету штраф (risk\_low -> 1 штраф, risk\_high -> 3 единицы штрафа), вычисляет `score` как (1 - штраф / макс_штраф) * 10. 

Пример:
```
> ./run.sh ./pisc    
> {"Unsafe operations": {"score": 6.8, "risk": "medium", "desc": "Uses regular expressions to identify potentially dangerous constructs in the code, e.g. network connection or dynamic execution"}}
```