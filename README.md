# dock try to re-implement docker

## usage

### run a container

```sh
    dock run -ti /bin/sh

    [root@localhost dock]# ./dock run -ti /bin/sh
    sh-4.2# ps
      PID TTY          TIME CMD
        1 pts/0    00:00:00 sh
        5 pts/0    00:00:00 ps
```

### limit memory usage

```sh
    dock run -ti -m 100m /bin/sh
```

### license

dock is release under MIT LICENSE