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