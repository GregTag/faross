## Примеры запуска контейнера
```sh
docker build -t scorecard .
docker run -it --rm scorecard pkg:pypi/django@1.11.1
```
Не работает, если исходный репозиторий не на гитхабе

Пример:
```sh
$ sudo docker run -it --rm scorecard pkg:pypi/django@1.11.1
[{"name":"Maintained","score":10,"reason":"30 commit(s) and 0 issue activity found in the last 90 days -- score normalized to 10","risk":"High"},{"name":"Code-Review","score":8,"re
ason":"Found 24/27 approved changesets -- score normalized to 8","risk":"High"},{"name":"CII-Best-Practices","score":0,"reason":"no effort to earn an OpenSSF best practices badge d
etected","risk":"Low"},{"name":"License","score":10,"reason":"license file detected","risk":"Low"},{"name":"Signed-Releases","score":-1,"reason":"no releases found","risk":"High"},
{"name":"Branch-Protection","score":-1,"reason":"internal error: error during branchesHandler.setup: internal error: githubv4.Query: Resource not accessible by integration","risk":
"High"},{"name":"Packaging","score":-1,"reason":"packaging workflow not detected","risk":"Medium"},{"name":"Dangerous-Workflow","score":10,"reason":"no dangerous workflow patterns 
detected","risk":"Critical"},{"name":"Security-Policy","score":9,"reason":"security policy file detected","risk":"Medium"},{"name":"Token-Permissions","score":0,"reason":"detected 
GitHub workflow tokens with excessive permissions","risk":"High"},{"name":"SAST","score":0,"reason":"SAST tool is not run on all commits -- score normalized to 0","risk":"Medium"},
{"name":"Binary-Artifacts","score":10,"reason":"no binaries found in the repo","risk":"High"},{"name":"Fuzzing","score":10,"reason":"project is fuzzed","risk":"Medium"},{"name":"Vu
lnerabilities","score":10,"reason":"0 existing vulnerabilities detected","risk":"High"},{"name":"Pinned-Dependencies","score":0,"reason":"dependency not pinned by hash detected -- 
score normalized to 0","risk":"Medium"}]
```