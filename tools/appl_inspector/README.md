# Application Inspector
Выполянет две проверки: **Unsafe operations** и **File types**

Скрипт `run.sh` принимает на вход путь до директории, которую необходимо проверить, запускает application inspector, записывает его вывод в файл `data.json`. Запускает по очереди все проверки, выводит результат в виде согласованных json-ов на отдельных строках. Если произошла ошибка (например, директория не была найдена), скрипт завершает работу с `error_code=1`, ничего не выводит в поток вывода, логирует ошибку в `stderr`.

Для корректной работы скрипта, команда `appinspector analyze` должна быть доступна как команда терминала. Application Inspector можно установить с помощью [dotnel tool](https://github.com/microsoft/ApplicationInspector/wiki/2.-NuGet-Support#installing-and-using-the-command-line-tool)

Пример:
```
> ./run.sh ./pisc         
> {"Unsafe operations": {"score": 6.8, "risk": "medium", "desc": "Uses regular expressions to identify potentially dangerous constructs in the code, e.g. network connection or dynamic execution"}}
{"File types": {"score": 10, "risk": "high", "desc": "Determines the presence of executable file extensions"}}
```

## Unsafe operations 
Из всей информации в data.json берёт только теги, ищет среди них опасные, присваивает пакету штраф (risk\_low -> 1 штраф, risk\_high -> 3 единицы штрафа), вычисляет `score` как (1 - штраф / макс_штраф) * 10. 

## Unsafe operations 
В `executable_extensions.txt` приведён список опасных (исполняемых) расширений файла, в `obfuscated_extensions.txt` - список обфусцированных расширений. Если в проекте встречается расширение из первой категории, то `score=3`, если из второй -  `score=1`. Если опасных раширений не встретилось, `score=10`. Количество разных опасных расширений не влияет на `score`, достаточно одного представителя от категории.
