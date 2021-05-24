# Tutorial ontenedor en Go

!!! warning "CUIDADO"
    Es muy recomendable utilizar una MV para realizar los siguientes pasos ya que pueden provocar problemas en el sistema si no se realizan correctamente.

Nuestro objetivo es que, al acabar este tutorial, tengamos un programa que sea capaz de crear un proceso y aislarlo con _namespaces_ y _cgroups_. De hecho, intentaremos que la interacción con éste sea muy parecida a la que tendríamos cuando ejecutamos un contenedor Docker:

```console
    # Nosotros vamos a ejecutarlo así:
root@bar:~$ go     run contenedor.go run         <command> <args>
    # Una ejecución en Docker sería algo de este estilo:
root@bar:~$ docker                   run <image> <command> <args>
```
!!! info ""
    Cabe destacar que es necesario que **todos los ficheros estén en una carpeta cuyo grupo y usuario pertenezca a root**, así como realizar todos los comandos con privilegios de root.

## Paso 1: Creación del código base
El primer paso consiste en **escribir las primeras dos funciones**: `main`, que simplemente comprobará que se ha ejecutado el comando correcto en terminal y `run`, que imprimirá en pantalla los datos del proceso y ejecutará otro nuevo que le indiquemos en parámetro.

```go
package main

import (
    "fmt"
    "os"
    "os/exec"
)

func main() {
    switch os.Args[1] {
        case "run":
            run()
        default:
            panic("¿Argumento Invalido?")
    }
}

func run() {
    fmt.Printf("Corriendo '%v' con User ID %d en PID %d \n", os.Args[2:], os.Getuid(), os.Getpid())

    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Run()
}
```

Ahora puedes probar a ejecutar la instrucción `go run contenedor.go run ps a` y comprobar que se existen dos procesos: el **contenedor.go** y el **ps a** que le hemos indicado que ejecute dentro del "pre-contenedor".

??? info "Explicación"

    La función `run` simplemente imprime por pantalla información útil sobre el proceso que estamos ejecutando y que, más adelante, creará el contenedor. De momento, lo único que estamos haciendo es indicarle que queremos ejecutar un comando con la función `Command` del paquete [`exec`](https://golang.org/pkg/os/exec/) indicándole los argumentos. Este comando devuelve una estructura del tipo `Cmd` en la que tenemos que especificarle el _Stdin_ _Stdout_ y _Stderr_.

    También podemos ejecutar otros comandos dentro del contenedor, como `go run contenedor.go run /bin/bash`, en cuyo caso se abrirá una nueva terminal.


## Paso 2: Aislando con Namespace UTS (Hostname)
Este _namespace_ permite cambiar tanto el _hostname_ como el _domain-name_ del contenedor sin que afecte a estos campos del host.

Para lograr este aislamiento debemos **añadir las siguientes líneas**, además de importar el paquete necesario [`syscall`](https://golang.org/pkg/syscall/):

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS,
    }
```

Con estos cambios lograremos que, al iniciar el contenedor con `go run contenedor.go run /bin/bash`, **podamos cambiar el _hostname_ dentro del contenedor** y que, al salir del mismo (saliendo del bash con un `exit`) no haya cambiado en el host.

??? abstract "Código completo hasta ahora"

    ```go
    package main

    import (
        "fmt"
        "os"
       "os/exec"
       "syscall"
    )

    func main() {
        switch os.Args[1] {
            case "run":
               run()
           default:
                panic("¿Argumento Invalido?")
        }
    }

    func run() {
        fmt.Printf("Corriendo '%v' con User ID %d en PID %d \n", os.Args[2:], os.Getuid(), os.Getpid())

        cmd := exec.Command(os.Args[2], os.Args[3:]...)
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        cmd.SysProcAttr = &syscall.SysProcAttr{
            Cloneflags: syscall.CLONE_NEWUTS,
        }

        cmd.Run()
    }
    ```

## Paso 3: Aislando con Namespace USER (username)
Con la inclusión de este _namespace_ vamos a **separar las tablas de UID y GID entre el host y el contenedor**, de tal forma que dentro del contenedor no haya los mismos usuarios que fuera.

Para crear el nuevo _namespace_ simplemente es **necesario añadir una _flag_** más al código:

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
                Cloneflags: syscall.CLONE_NEWUTS |
                            syscall.CLONE_NEWUSER,
        }
```

El problema es que ahora mismo, cuando ejecutamos el contenedor, nos informa que el proceso que lo llama tiene UID 0 pero si comprobamos el usuario que se nos asignó en la nueva tabla de UID (o sea, en la tabla del contenedor) es un usuario "aleatorio":

```console
root@bar:~$ go run contenedor.go run /bin/bash
Corriendo '[/bin/bash]' con User ID 0 en PID 2792
root@bar:~$ id
uid=65534(nobody) gid=65534(nogroup) groups=65534(nogroup)
```

Así que le vamos a indicar que **mapee el usuario de fuera del contenedor (UID 0) con el que queramos dentro**:

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWUTS |
                syscall.CLONE_NEWUSER,
    UidMappings: []syscall.SysProcIDMap{
        {
            ContainerID: 0,           // UID dentro del contenedor
            HostID: os.Getuid(),      // UID en el host
            Size: 1,                  // Quiero mapear solo un usuario
        },
    },
    GidMappings: []syscall.SysProcIDMap{
        {
            ContainerID: 0,
            HostID: os.Getgid(),
            Size: 1,
        },
    },
}
```

??? abstract "Código completo hasta ahora"

    ```go
    package main

    import (
        "fmt"
        "os"
        "os/exec"
        "syscall"
    )

    func main() {
        switch os.Args[1] {
            case "run":
                run()
            default:
                panic("¿Argumento Invalido?")
        }
    }

    func run() {
        fmt.Printf("Corriendo '%v' con User ID %d en PID %d \n", os.Args[2:], os.Getuid(), os.Getpid())

        cmd := exec.Command(os.Args[2], os.Args[3:]...)
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        cmd.SysProcAttr = &syscall.SysProcAttr{
            Cloneflags: syscall.CLONE_NEWUTS  |
                        syscall.CLONE_NEWUSER,
            UidMappings: []syscall.SysProcIDMap{
                {
                    ContainerID: 0,             // UID dentro del container
                    HostID: os.Getuid(),        // UID en el host
                    Size: 1,                    // Quiero mapear solo unuser
                },
            },
            GidMappings: []syscall.SysProcIDMap{
                {
                    ContainerID: 0,
                    HostID: os.Getpid(),
                    Size: 1,
                },
            },
        }
        cmd.Run()
    }
    ```


## Paso 4: Aislando con Namespace NS (Mount)
Este fue el primer _Namespace_ que se incluyó en el kernel de Linux y uno de los más sencillos: simplemente aisla los puntos de montaje. De esta forma podemos **esconder los _mounts_ entre el host y el contenedor y viceversa**.

Para ver los puntos de montaje usados en cada una de las máquinas con el comando `mount`.

Para añadir esta característica debemos incluir la _flag_ apropiada:

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWUTS  |
                syscall.CLONE_NEWNS   |
                syscall.CLONE_NEWUSER,

                {...}
}
```

??? abstract "Código completo hasta ahora"

    ```go
    package main

    import (
        "fmt"
        "os"
        "os/exec"
        "syscall"
    )

    func main() {
        switch os.Args[1] {
            case "run":
                run()
            default:
                panic("¿Argumento Invalido?")
        }
    }

    func run() {
        fmt.Printf("Corriendo '%v' con User ID %d en PID %d \n", os.Args[2:], os.Getuid(), os.Getpid())

        cmd := exec.Command(os.Args[2], os.Args[3:]...)
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        cmd.SysProcAttr = &syscall.SysProcAttr{
            Cloneflags: syscall.CLONE_NEWUTS |
                        syscall.CLONE_
                        syscall.CLONE_NEWUSER,
            UidMappings: []syscall.SysProcIDMap{
                {
                    ContainerID: 0,             // UID dentro del container
                    HostID: os.Getuid(),        // UID en el host
                    Size: 1,                    // Quiero mapear solo unuser
                },
            },
            GidMappings: []syscall.SysProcIDMap{
                {
                    ContainerID: 0,
                    HostID: os.Getpid(),
                    Size: 1,
                },
            },
        }
        cmd.Run()
    }
    ```


## Paso 5: Aislando con Namespace PID
El PID _namespace_ permite separar los árboles de procesos, de tal forma que dentro del **contenedor no se pueden ver los procesos del host**.

Para añadir este _namespace_ simplemente incluimos la _flag_ apropiada:

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWUTS  |
                syscall.CLONE_NEWUSER |
                syscall.CLONE_NEWNS   |
                syscall.CLONE_NEWPID,

                {...}
}
```

Sin embargo, cuando ejecutamos un `ps a` seguimos pudiendo ver los mismos procesos de antes.

??? info "Explicación: Mount Namespace no aísla los procesos"

    Es importante saber que **_/proc_** es un pseudo-filesystem montado por el sistema operativo por defecto donde se muestra la información sobre los procesos. Cuando hacemos un **_ps a_**, lo que está pasando realmente es que esta instrucción está consultando los datos del directorio anteriormente nombrado.


La solución es asignar un nuevo **_/proc_** en la raíz del contenedor. Para ello necesitamos un nuevo _root filesystem_ como Alpine (que continene únicamente los archivos necesarios para que funcione un contenedor).

## Paso 6: Añadiendo un Filesystem para el contenedor
Para realizar este paso necesitamos descargar el [mini-root](https://alpinelinux.org/downloads/) de Alpine. Lo descomprimimos y lo llamamos, por ejemplo, `alpinefs` y le cambiamos el usuario con `chown root alpinefs/`.

### Montamos nuestro propio /proc
Necesitamos un nuevo directorio **_proc_** para que el comando `ps a` pueda acceder a él para acceder a la información de los procesos del contnedor.

??? info "Explicación: directorio /proc"

    Otra cosa que se podría intuir es que es necesario añadir el _Namespace NS (de Mount)_ para aislar ambos directorios. Pero no, este último comentario es falso pese a que existan muchas referencias en la red a que es completamente necesario: cuando un proceso como `ps` quiere comprobar los procesos activos en `/proc` lo que hace es ir directamente a ese archivo. Nuestro proceso, tanto con el _Namespace NS_ como sin él, va a seguir mirando los procesos en la carpeta `/proc`, es decir, la que está justo debajo del directorio raíz y no en la del nuevo _root filesystem_ de alpine. Así que podríamos montar nuestro nuevo `proc/` sin el _Namespace NS_.


La solución de que se muestren únicamente los procesos activos de nuestro contenedor se divide en dos pasos, pero antes, debemos cambiar un poco la forma en la que habíamos planteado el programa en un principio.

Ahora, en vez de ejecutar desde la función `run` la instrucción indicada en los parámetros, vamos a duplicar el proceso actual llamando a `/proc/self/exe` para que en esta segunda ejecución se cambie el flujo del programa y no pase por la función `run`, sino por la función `child`.

```go
cmd := exec.Command ("/proc/self/exe", append([]string {"child"}, os.Args[2:]...)...)
```

De esta forma, habría otra función dentro del programa que se ejecutaría la segunda vez, donde implementamos la solución a nuestro último problema:

1. Hacer la nueva raíz de nuestro contenedor la raíz del _filesystem_ que acabamos de descargar (`alpinefs/`) para que al acceder a `/proc` esté accediendo al nuevo y no al del Host. Esto se puede hacer tanto con la llamada al sistema `chroot` o `pivot_root`. La segunda opción es más segura, aunque más complicada. Por lo tanto, para evitar aumentar demasiado la complejidad se utilizará el primer método (anexando el segundo al final del tutorial).


2. Montar el _filesystem_ `proc` para que el sistema pueda utilizarlo para almacenar información sobre los procesos.

```go
func child() {
        fmt.Printf("Running '%v' as user %d in PID %d \n", os.Args[2:], os.Getuid(), os.Getpid())

        must(syscall.Chroot("alpinefs/"))
        must(os.Chdir("/"))

        must(syscall.Mount("proc", "proc", "proc", 0, ""))

        cmd := exec.Command(os.Args[2], os.Args[3:]...)
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        defer func() {
            must(syscall.Unmount("proc", 0))
        }()

        must(cmd.Run())
}
```

??? info "Explicación: ¿por qué crear una nueva función?"

    Ahora no sólo vamos a añadir _namespaces_ y ejecutar una instrucción sino que vamos a realizar otras acciones. Si cogemos el flujo de la función `run` y realizamos las nuevas acciones después de `cmd.Run()` no se estarían completando hasta que acabara esta última orden. A su vez, si introducimos las acciones antes de `cmd.Run()` no se habrían creado aún los_namespaces: es justo mientras transcurre en `cmd.Run()` cuando queremos modificar el contenedor.

    Por eso una opción es obligar al proceso a llamarse a una copia de sí mismo y cambiar el flujo del programa a la nueva función `child`.

    Cabe destacar que el _filesystem_ propuesto de Alpine no cuenta con Bash, así que tendríamos que mandar ejecutar `/bin/sh`

    Por otro lado, es recomendable que a partir de ahora empezemos a manejar los errores que nos puedan aparecer:

    ```go
    func must(err error) {
        if err != nil {
            panic(err)
        }
    }
    ```



??? abstract "Código completo hasta ahora"

    ```go
    package main

    import (
            "fmt"
            "os"
            "os/exec"
            "syscall"
    )


    func main() {
        switch os.Args[1] {
            case "run":
                run()
            case "child":
                child()
            default:
                panic("¿Argumento Invalido?")
        }
    }

    func run() {
            fmt.Printf("Corriendo '%v' con User ID %d en PID %d \n", os.Args[2:], os.Getuid(), os.Getpid())

            cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
            cmd.Stdin = os.Stdin
            cmd.Stdout= os.Stdout
            cmd.Stderr = os.Stderr

            cmd.SysProcAttr = &syscall.SysProcAttr{
                    Cloneflags: syscall.CLONE_NEWUTS  |
                                syscall.CLONE_NEWUSER |
                                syscall.CLONE_NEWNS   |
                                syscall.CLONE_NEWPID,
                    UidMappings: []syscall.SysProcIDMap{
                            {
                                    ContainerID: 0,
                                    HostID: os.Getuid(),
                                    Size: 1,
                            },
                    },
                    GidMappings: []syscall.SysProcIDMap{
                            {
                                    ContainerID: 0,
                                    HostID: os.Getgid(),
                                    Size: 1,
                            },
                    },
            }

            must(cmd.Run())
    }

    func child() {
            fmt.Printf("Corriendo '%v' con User ID %d en PID %d \n", os.Args[2:], os.Getuid(), os.Getpid())

            must(syscall.Sethostname([]byte("container")))

            must(syscall.Chroot("alpinefs/"))
            must(os.Chdir("/"))
            must(syscall.Mount("proc", "proc", "proc", 0, ""))

            cmd := exec.Command(os.Args[2], os.Args[3:]...)
            cmd.Stdin = os.Stdin
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr

            defer func() {
                    must(syscall.Unmount("proc", 0))
            }()

            must(cmd.Run())
    }


    func must(err error) {
            if err != nil {
                    panic(err)
            }
    }
    ```


??? info "¿QUÉ HEMOS CONSEGUIDO HASTA AHORA?"

    En estos momentos hemos conseguido introducir unos cuantos _namespaces_, al menos los más significativos para realizar en este tutorial.

    El <b>hostname namespace</b> se puede comprobar de esta forma:

    ```console
        # Fuera del contenedor
    root@bar:~$ hostname
    host
    root@bar:~$ go run contenedor.go run /bin/sh
    Corriendo '[/bin/sh]' con User ID 0 en PID 72724

        # Dentro del contenedor
    root@bar:~$ sethostname contenedor
    root@bar:~$ exit

        # Fuera del contenedor
    root@bar:~$ hostname
    host
    ```

    El <b>user namespace</b> lo hemos conseguido introducir añadiendo los mapeos de usuario a root dentro del contenedor. Lo podemos comprobar de esta forma:

    ```console
        # Fuera del contenedor
    root@bar:~$ go run contenedor.go run /bin/sh
    Corriendo '[/bin/sh]' con User ID 0 en PID 72724

        # Dentro del contenedor
    root@bar:~$ id
    uid=0(root) gid=0(root) groups=0(root)
    ```

    El <b>mount namespace</b> se puede comprobar de una forma muy sencilla:

    ```console
        # Fuera del contenedor
    root@bar:~$ mount
    ###### Aparecen muchos puntos de montaje usados por el host
    root@bar:~$ go run contenedor.go run /bin/sh
    Corriendo '[/bin/sh]' con User ID 0 en PID 72724

        # Dentro del contenedor
    root@bar:~$ mount
    proc on /proc type proc (rw,relatime)
    ```

    El <b>pid namespace</b> lo podemos comprobar realizando las siguientes instrucciones:

    ```console
        # Fuera del contenedor
    root@bar:~$ go run contenedor.go run /bin/sh
    Corriendo '[/bin/sh]' con User ID 0 en PID 72724

        # Dentro del contenedor
    root@bar:~$ ps a
    PID     USER    TIME   COMMAND
        1   root     0:00  /proc/self/exe child /bin/sh
        5   root     0:00  /bin/sh
    11   root     0:00  ps a
    ```


## Añadiendo Cgroups (memoria y PID)

En este ejemplo añadiremos un límite al número máximo de procesos en el cgroup (y, por lo tanto, en el contenedor) permitidos. Para ello necesitamos crear un nuevo directorio en `/sys/fs/cgroup/pids/`. Al crear el directorio automáticamente el sistema añade los archivos necesarios para mostrar los datos del nuevo _Cgroup_ y para modificar los límites que se le quieran añadir.

En nuestro caso el grupo se llamará `demo`. Para modificar el número máximo de procesos que se permite en el contenedor sólo es necesario modificar el archivo donde se indica el número (pondremos como máximo 12 procesos) y otro donde se introduce al proceso del contenedor en el grupo de control.

Además añadiremos un número máximo de bytes de memoria que se le asignan al contenedor, aunque esto es más complicado de comprobar que funciona correctamente, pero los pasos son los mismos que en el anterior caso.

```go
func cg()
    cgroups := "/sys/fs/cgroup"
    // Creando cgroup para PIDs
    pids := filepath.Join(cgroups, "pids/demo")
    if _, err := os.Stat(pids); os.IsNotExist(err) {
        must(os.Mkdir(pids, 0755))
    }
    // Creando cgroup para PIDs
    memory := filepath.Join(cgroups, "memory/demo")
    if _, err := os.Stat(memory); os.IsnotExist(err) {
        must(os.Mkdir(memory, 0755))
    }

    //Establecemos limite y metemos al proceso dentro del grupo de procesos
    must(ioutil.WriteFile(filepath.Join(pids, "pids.max"), []byte("10"), 0700))
    must(ioutil.WriteFile(filepath.Join(pids, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))

    must(ioutil.WriteFile(filepath.Join(memory, "memory.limit_in_bytes"), []byte("2M"), 0700))
    must(ioutil.WriteFile(filepath.Join(memory, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}
```

Sólo hace falta llamar a esta función desde el principio de la función `child`.


## Anexo: mejora con pivot_root

Cuando se introdujo en el tutorial la llamada al sistema `chroot` se mencionó la posibilidad de utilizar otra más segura: `pivot_root`.

Aunque antiguamente, en los primeros contenedores, se utilizaba la primera opción, se ha llegado a la conclusión de que tiene varios fallos de seguridad que permiten "salir o escapar" del contenedor. `pivot_root` aprovecha el _mount namespace_ ya que permite hacer _unmount_ del antiguo root y no lo hace accesible en el _namespace_ del contenedor.

Si usamos únicamente `chroot` podemos acceder al Host con el siguiente comando: `chroot /proc/1/root`.

Lo que hay que saber para poder usar `pivot_root` es que necesita dos argumentos, el primero es la dirección del nuevo directorio raíz (no viene en la documentación pero debe estar montado sobre sí mismo con la opción `bind`) y el segundo es la dirección donde se va a situar el antiguo directorio raíz.

El código completo del tutorial quedaría así:

??? abstract "CÓDIGO COMPLETO"

    ```go
    package main

    import (
        "fmt"
        "os"
        "os/exec"
        "syscall"
        "io/ioutil"
        "strconv"
        "path/filepath"
    )


    func main() {
        switch os.Args[1] {
        case "run":
            run()
        case "child":
            child()
        default:
            panic("¿Argumento Invalido?")
        }
    }


    func run() {
        fmt.Printf("Corriendo '%v' con User ID %d en PID %d \n", os.Args[2:], os.Getuid(), os.Getpid())

        cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
        cmd.Stdin = os.Stdin
        cmd.Stdout= os.Stdout
        cmd.Stderr = os.Stderr

        cmd.SysProcAttr = &syscall.SysProcAttr{
            Cloneflags: syscall.CLONE_NEWUTS  |
                    syscall.CLONE_NEWUSER |
                    syscall.CLONE_NEWNS   |
                    syscall.CLONE_NEWPID,
            UidMappings: []syscall.SysProcIDMap{
                {
                    ContainerID: 0,
                    HostID: os.Getuid(),
                    Size: 1,
                },
            },
            GidMappings: []syscall.SysProcIDMap{
                {
                    ContainerID: 0,
                    HostID: os.Getgid(),
                    Size: 1,
                },
            },
        }

        must(cmd.Run())
    }



    func child() {
        fmt.Printf("Corriendo '%v' con User ID %d en PID %d \n", os.Args[2:], os.Getuid(), os.Getpid())

        cg()

        must(syscall.Sethostname([]byte("container")))

        pivot()
        must(syscall.Mount("proc", "proc", "proc", 0, ""))

        cmd := exec.Command(os.Args[2], os.Args[3:]...)
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        must(syscall.Unmount(".old_root", syscall.MNT_DETACH))
        must(os.Remove(".old_root"))

        defer func() {
            must(syscall.Unmount("proc", 0))
        }()


        must(cmd.Run())
    }


    func cg() {
        cgroups := "/sys/fs/cgroup"

        pids := filepath.Join(cgroups, "pids/demo")
        if _, err := os.Stat(pids); os.IsNotExist(err) {
            must(os.Mkdir(pids, 0755))
        }

        memory := filepath.Join(cgroups, "memory/demo")
        if _, err := os.Stat(memory); os.IsNotExist(err) {
            must(os.Mkdir(memory, 0755))
        }
        
        must(ioutil.WriteFile(filepath.Join(pids, "pids.max"), []byte("22"), 0700))
        must(ioutil.WriteFile(filepath.Join(pids, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))

        must(ioutil.WriteFile(filepath.Join(memory, "memory.limit_in_bytes"), []byte("2M"), 0700))
        must(ioutil.WriteFile(filepath.Join(memory, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
    }

    func pivot() {
        must(syscall.Mount("alpinefs", "alpinefs", "", syscall.MS_BIND|syscall.MS_REC, ""))
        if _, err := os.Stat("alpinefs/.old_root"); os.IsNotExist(err) {
            must(os.Mkdir("alpinefs/.old_root", 0700))
        }
        must(syscall.PivotRoot("alpinefs", "alpinefs/.old_root"))
        must(os.Chdir("/"))
    }

    func must(err error) {
        if err != nil {
            panic(err)
        }
    }
    ```
