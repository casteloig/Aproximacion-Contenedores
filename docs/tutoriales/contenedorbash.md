# Contenedor en Bash

!!! warning
    Es muy recomendable utilizar una MV para realizar los siguientes pasos ya que pueden provocar problemas en el sistema si no se realizan correctamente.

## Namespaces
Para crear un contenedor en terminal necesitamos muchas menos instrucciones que para crearlo en Go: contamos con comandos que facilitan enormemente la tarea.

En el siguiente ejemplo se creará lo mismo que en el tutorial de Go pero en tan sólo unas pocas líneas.

Para el caso de los namespaces simplemente necesitamos utilizar la orden [`unshare`](https://man7.org/linux/man-pages/man1/unshare.1.html): en este caso crearemos los correspondientes a `--mount`, `--pid`, `--uts` y `--user`.

Además, vamos a especificar con `--fork` que la terminal que va a abrir el programa sea un hijo del proceso `unshare` (es útil precisamente porque hemos creado un nuevo _Namespace PID_). Con `--map-root-user` ahorramos todas las líneas que escribíamos en Go para hacer que el usuario del nuevo _Namespace USER_ sea root. Por último, podemos ahorrar realizar el `pivot_root` del que ya hemos hablado anteriormente con la opción `root=[dirección]`.

Por otro lado, es importante acordarse de montar el nuevo _pseudo-filesystem_ de `/proc`.


```console
root@bar~$ unshare --mount --pid --uts --user --fork --map-root-user --root=alpinefs /bin/sh
root@bar~$ mount -t proc proc proc
root@bar~$ hostname demo
```

## Cgroups
En el caso de los grupos de control, tiene la misma dificultad que en el ejemplo en Go: simplemente hay que crear los directorios correspondientes para cada tipo de _Cgroup_, eso sí, desde una terminal con root en el Host, ya que tenemos que introducir el PID en los correspondientes `cgroup.procs` del _Namespace USER_ superior (es decir, el que vemos desde el Host). 

```go
root@bar~$ ps aux | grep /sh
root    2552    pts/0   00:13   /bin/sh
root@bar~$ mkdir /sys/fs/cgroup/pids/demo
root@bar~$ mkdir /sys/fs/cgroup/memory/demo
root@bar~$ echo 20 /sys/fs/cgroup/pids/demo/pids.max
root@bar~$ echo 2552 /sys/fs/cgroup/pids/demo/cgroup.procs
root@bar~$ echo "2M" /sys/fs/cgroup/memory/demo/memory.limit_in_bytes
root@bar~$ echo 2552 /sys/fs/cgroup/memory/demo/cgroup.procs
```
