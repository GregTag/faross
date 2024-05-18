# FAROSS

## Build Image
```bash
$ docker build -t faross .
```

## Run container
```bash
$ docker run -it -p 5002:5002 -v /tmp:/tmp -v /var/run/docker.sock:/var/run/docker.sock --name faross-main faross
```
