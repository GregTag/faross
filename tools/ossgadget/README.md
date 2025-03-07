#  OSSGadget
Выполяняет проверку **Backdoors**

Скрипт `run.sh` принимает на вход purl, запускает утилиту oss-detect-backdoor инструмента OSSGadget, записывает её вывод в файл `data.json`. На основании количества тэгов хотя бы с одной проверкой > low-confidence, вычисляет `score` и выводит json в установленном формате. Если произошла ошибка (например, директория не была найдена), скрипт завершает работу с `error_code=1`, ничего не выводит в поток вывода, логирует ошибку в `stderr`. 

Пример:
```
> ./run.sh pkg:npm/execa    
> {"Backdoors": {"score": 8, "risk": "medium", "desc": "Uses regular expressions to identify backdoors"}}
```

Для корректной работы, утилита `oss-detect-backdoor` должна быть доступна как команда терминала. У OSSGadget есть написанный [докер-файл](https://github.com/microsoft/OSSGadget/wiki/Docker-Image)