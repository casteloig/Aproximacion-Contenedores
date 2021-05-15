# Tutorial ontenedor en Go

!!! warning
    Es muy recomendable utilizar una MV para realizar los siguientes pasos ya que pueden provocar problemas en el sistema si no se realizan correctamente.

Nuestro objetivo es que, al acabar este tutorial, tengamos un programa que sea capaz de crear un proceso y aislarlo con _namespaces_ y _cgroups_. De hecho, intentaremos que la interacción con éste sea muy parecida a la que tendríamos cuando ejecutamos un contenedor Docker:

```console
    # Nosotros vamos a ejecutarlo así:
root@bar:~$ go     run contenedor.go run         <command> <args>
    # Una ejecución en Docker sería algo de este estilo:
root@bar:~$ docker                   run <image> <command> <args>
```

Cabe destacar que es necesario que **todos los ficheros estén en una carpeta cuyo grupo y usuario pertenezcan a root**, así como realizar todos los comandos con privilegios de root.

## Paso 1: Creación del código base
El primer paso consiste en **escribir las primeras dos funciones**: `main`, que simplemente comprobará que se ha ejecutado el comando correcto en terminal y `run`, que imprimirá en pantalla los datos del proceso y ejecutará otro nuevo.

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

<details>
<summary>Explicación</summary>

La función `run` simplemente imprime por pantalla información útil sobre el proceso que estamos ejecutando y que, más adelante, creará el contenedor. De momento, lo único que estamos haciendo es indicarle que queremos ejecutar un comando con la función `Command` del paquete [`exec`](https://golang.org/pkg/os/exec/) indicándole los argumentos. Este comando devuelve una estructura del tipo `Cmd` en la que tenemos que especificarle el _Stdin_ _Stdout_ y _Stderr_.

También podemos ejecutar otros comandos dentro del contenedor, como `go run contenedor.go run /bin/bash`, en cuyo caso se abrirá una nueva terminal.
</details>

</br>

## Paso 2: Aislando con Namespace UTS (Hostname)
Este _namespace_ permite cambiar tanto el _hostname_ como el _domain-name_ del contenedor sin que afecte a estos campos del host.

Para lograr este aislamiento debemos **añadir las siguientes líneas**, además de importar el paquete necesario [`syscall`](https://golang.org/pkg/syscall/):

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS,
    }
```

Con estos cambios lograremos que ,al iniciar el contenedor con `go run contenedor.go run /bin/bash`, **podamos cambiar los el _hostname_ dentro del contenedor** y que, al salir del contenedor (saliendo del bash con un `exit`) no haya cambiado en el host.

<details>
<summary>Código completo hasta ahora</summary>

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
</details>

## Paso 3: Aislando con Namespace USER (username)
Con la inclusión de este _namespace_ vamos a **separar las tablas de UID y GID entre el host y el contenedor**, de tal forma que dentro del contenedor no haya los mismos usuarios que fuera.

Para crear el nuevo _namespace_ simplemente es **necesario añadir una _flag_** más al código:

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
                Cloneflags: syscall.CLONE_NEWUTS |
                            syscall.CLONE_NEWUSER,
        }
```

El problema es que ahora mismo, cuando ejecutamos el contenedor, nos informa que el proceso que lo llama tiene UID 0 pero si comprobamos el usuario que se nos asignó en la nueva tabla de UID nos asigna un usuario "aleatorio":

```console
root@bar:~$ go run contenedor.go run /bin/bash
Corriendo '[/bin/bash]' con User ID 0 en PID 2792
root@bar:~$ id
uid=65534(nobody) gid=65534(nogroup) groups=65534(nogroup)
```

Así que le vamos a indicar que **mapee el usuario fuera del contenedor (UID 0) con el que queramos dentro**:

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWUTS |
                syscall.CLONE_NEWUSER,
    UidMappings: []syscall.SysProcIDMap{
        {
            ContainerID: 0,           // UID dentro del container
            HostID: os.Getuid(),      // UID en el host
            Size: 1,                  // Quiero mapear solo unuser
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

<details>
<summary>Código completo hasta ahora</summary>

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
</details>

</br>

## Paso 4: Aislando con Namespace NS (Mount)
Este fue el primer _Namespace_ que se incluyó en el kernel de Linux y uno de los más sencillos: simplemente aisla los puntos de montaje. De tal forma que podemos **esconder los montajes del host en el contenedor y viceversa**.

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

<details>
<summary>Código completo hasta ahora</summary>

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
</details>

</br>

## Paso 5: Aislando con Namespace PID
El PID _namespace_ permite separar los árboles de procesos, de tal forma que dentro del **contenedor no se pueden ver los procesos del host**.

Para añadir este _namespace_ simplemente incluimos la _flag_ apropiada:

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWUTS  |
                syscall.CLONE_NEWUSER |
                syscall.CLONE_NEWNS   |
                syscall.CLONE_NEWPID,
}
```

Sin embargo, cuando ejecutamos un `ps a` seguimos pudiendo ver los mismos procesos de antes.

<details>
<summary>Explicación: Mount Namespace no aísla los procesos</summary>

Es importante saber que **_/proc_** es un pseudo-filesystem montado por el sistema operativo por defecto donde se muestra la información sobre los procesos. Cuando hacemos un **_ps a_**, lo que está pasando realmente es que esta instrucción consulta los datos del directorio anteriormente nombrado.
</details>

La solución es asignar un nuevo **_/proc_** en la raíz del contenedor. Para ello necesitamos un nuevo _root filesystem_ como Alpine (que continene únicamente los archivos necesarios para que funcione un contenedor).

### Añadiendo un Filesystem para el contenedor
Para realizar este paso necesitamos descargar el [mini-root](https://alpinelinux.org/downloads/) de Alpine. Lo descomprimimos y lo llamamos, por ejemplo, `alpinefs` y le cambiamos el usuario con `chown root alpinefs/`.

