Для запуска контейнера:

```
 docker build -t pa-image .
```

```bash
docker run --cgroupns=host --privileged --rm -ti \
 pa-image analyze \
 -dynamic-bucket file:///results/ -execution-log-bucket file:///results \
 -ecosystem pypi -package pillow -version 10.3.0 \
 && python3 console_log.py
```

ecosystems: pypi, npm