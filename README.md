# dock try to re-implement docker

dock try to implement most feature in docker

## Todo

- Daemon mode
- Container image
- Container volume
- Container network
- Container management
- Container logs
- ...

## Usage

### Run a container

```sh
[root@localhost dock]# ./dock run -ti /bin/sh
sh-4.2# ps
PID     TTY     TIME CMD
1       pts/0   00:00:00 sh
5       pts/0   00:00:00 ps
```

### Limit memory usage and cpushare

```sh
dock run -ti -m 100m - cpushare 512 sh
```

### License

dock is release under MIT LICENSE