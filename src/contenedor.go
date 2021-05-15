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

    fmt.Printf(strconv.Itoa(os.Getpid()))
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