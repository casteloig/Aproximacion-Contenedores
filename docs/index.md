# Creación de un contenedor en Go
Nuestro objetivo es que, al acabar este tutorial, tengamos un programa que sea capaz de crear un proceso y aislarlo con _namespaces_ y _cgroups_. De hecho, intentaremos que la interacción con éste sea muy parecida a la que tendríamos cuando ejecutamos un contenedor Docker:

```console
    # Nosotros vamos a ejecutarlo así:
root@bar:~$ go     run contenedor.go run         <command> <args>
    # Una ejecución en Docker sería algo de este estilo:
root@bar:~$ docker                   run <image> <command> <args>
```
> Cabe destacar que es necesario que **todos los ficheros estén en una carpeta cuyo grupo y usuario pertenezcan a root**, así como realizar todos los comandos con privilegios de root.

## Paso 1: Creación del código base
El primer paso consiste en escribir las primeras dos funciones: `main`, que simplemente comprobará que se ha ejecutado el comando correcto en terminal y `run`, que imprimirá en pantalla los datos del proceso y ejecutará otro nuevo.

<details>
<summary>Código en Go necesario</summary>

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
</details>

Ahora puedes probar a ejecutar la instrucción `go run contenedor.go run ps a` y comprobar que se existen dos procesos: el **contenedor.go** y el **ps a** que le hemos indicado que ejecute dentro del "pre-contenedor".

<details>
<summary>Explicación</summary>

>La función `run` simplemente imprime por pantalla información útil sobre el proceso que estamos ejecutando y que, más adelante, creará el contenedor. De momento, lo único que estamos haciendo es indicarle que queremos ejecutar un comando con la función `Command` del paquete [`exec`](https://golang.org/pkg/os/exec/) indicándole los argumentos. Este comando devuelve una estructura del tipo `Cmd` en la que tenemos que especificarle el _Stdin_ _Stdout_ y _Stderr_.


>También podemos ejecutar otros comandos dentro del contenedor, como `go run contenedor.go run /bin/bash`, en cuyo caso se abrirá una nueva terminal.
</details>








For full documentation visit [mkdocs.org](https://www.mkdocs.org).

## Commands

* `mkdocs new [dir-name]` - Create a new project.
* `mkdocs serve` - Start the live-reloading docs server.
* `mkdocs build` - Build the documentation site.
* `mkdocs -h` - Print help message and exit.

## Project layout

    mkdocs.yml    # The configuration file.
    docs/
        index.md  # The documentation homepage.
        ...       # Other markdown pages, images and other files.
