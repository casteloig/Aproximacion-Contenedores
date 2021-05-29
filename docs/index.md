# Una aproximación a los contenedores
## ¿Qué es un contenedor?
La virtualización es un proceso mediante el cual un software es usado para crear una abstracción sobre unos recursos, dando la sensación de que los elementos hardware se dividen en varias computadores virtuales.

Existen dos tipos técnicas principales de virtualización:

1. Máquinas virtuales.
2. Contenedores (virtualización ligera).

!!! note "Definición"
    Un **contenedor** es un límite lógico que se crea dentro de un sistema operativo proporcionado por el aislamiento de recursos hardware.

    La característica principal es que en esta técnica se utilizan herramientas que proporciona el Linux Kernel (como _cgroups_ y _namespaces_).

## Diferencia entre MV y contenedores

<figure>
<img src="./img/vm_vs_cont.png" width="650" />
<figcaption>Diferencia infraestructura: VM y contenedor</figcaption>
</figure>

* Un contenedor es una forma de **virtualización ligera**. 
* Normalmente envuelve a un **pequeño grupo de procesos**.
* Los contenedores **comparten el _kernel_** con el host.
* Dentro del contendor se encuentra **únicamente el código, librerías y ejecutables estrictamente necesarios**.


!!! info ""

    === "Usos principales"

        * **Microservicios**: los contenedores **son ligeros y envuelven servicios muy pequeños**, lo que los hace muy aptos para su uso en microservicios.
        * **DevOps + CI/CD**: facilita el ciclo _"build, test and deploy".
        * **Cloud**: los contenedores pueden **funcionar** de forma consistente en prácticamente **cualquier lugar**.

    === "Ventajas principales"

        * **Rapidez y ligereza**.
        * **Portabilidad e independencia de plataformas**.
        * **Escalabilidad**.

    === "Desventajas principales"

        * **Seguridad**: el aislamiento entre contenedores y el host es menor: puede provocar vulnerabilidades.
        * **Monitorización**: existen varias capas que monitorizar aunque sólo se tenga una aplicación en un contenedor.
        * **Complejidad a gran escala**.


## Historia de contenedores

### Chroot (1979)
`chroot` es una llamada al sistema que permite **cambiar el directorio raiz de un proceso y de sus hijos** a un nuevo lugar dentro del sistema de archivos. Esto hace que un grupo de procesos tengan una **visión limitada del almacenamiento** del sistema.

El principal **problema** de `chroot` es que los procesos con permisos de root se pueden **saltar el aislamiento** con facilidad.

### FreeBSD Jails (2000)
Fue el **primer "pseudo-contenedor"** y el impulsor de las tecnologías de los contenedores que existen hoy en día.

Su objetivo era **"confinar el root omnipotente"**, que da nombre al documento donde se presenta la herramienta.

Las Jails dan uso a `chroot` añadiendo nuevos mecanismos existentes en el sistema operativo. En cada _jail_ los procesos pueden **manipular únicamente los servicios y archivos a los que se le da acceso** (por ello se le dio el nombre de _jail_: cárcel en inglés, como si los procesos estuvieran en prisión). El administrador del sistema puede separar el sistema en varias celdas asignando un superusuario a cada una **sin perder el control** del sistema completo.

Para **mantener la seguridad** deseada se **desactivaron ciertas llamadas al sistema** (para evitar, por ejemplo, el _spoofing_), aunque esto impida utilizar ciertas instrucciones comunes como el comando `ping`.

Por último, **cada celda tiene sus propios UID y GID**: un mismo usuario en una celda puede corresponderse con otro usuario en distinta celda.

### Solaris Zones (2004)
Los creadores de esta tecnología no sólo buscaban **mantener la seguridad**, sino también **mejorar el uso de recursos a gran escala**. De hecho, se pretendía dar soporte a aplicaciones comerciales potencialmente escalables.

El objetivo era lograr que los **administradores** necesitaran pocos minutos para **configurar y lanzar una nueva _zone_**: el sistema se encargaría de crearla automáticamente añadiendo los límites en recursos compartidos que se consumían (inicialmente sólo CPU).

El administrador puede configurar las propiedades de las _zones_ en tiempo real o mediante scripts.

Existen dos tipos de _zones_:

* _Global zone_: es la _zone_ por defecto que tiene control sobre todos los procesos. Siempre existe aunque no se haya creado ninguna manualmente.
* _Non-global zone_: son las _zones_ configuradas desde la global.

Algo a tener en cuenta es que las Zones se preocupan por **mantener ciertas utilidades** sin dejar a un lado la seguridad (por ejemplo, permite utilizar el comando `ping` a diferencia de las _jails_).

### Cgroups y Namespaces (2006-2007)
En el 2006 se propuso un _framework_ que **agrupa procesos** y aprovecha los mecanismos de control existentes del _kernel_. El objetivo era que los **usuarios se centren en el controlador de recursos** y se abstraigan de cómo los procesos son gestionados y monitorizados.

Este mecanismo en un principio se llamaba _process containers_ pero se cambió el nombre a **_cgroups_** para diferenciarlo de los contenedores que conocemos hoy en día.

Por otro lado, en el 2002 se creó el primer **_namespace_**, el _mount namespace_. Aunque fue más adelante a partir del año 2007 cuando se comenzó a ver el potencial que podían tener junto a otras funcionalidades como los _cgroups_, así que se empezaron a desarrollar nuevos _namespaces_. En la actualidad existen ocho.

### LXC (2008)
Los LXC se podrían considerar los primeros contenedores tal y como conocemos el concepto hoy en día: fue la **primera tecnología en aplicar los dos elementos principales de los contenedores: _cgroups_ y _namespaces_**.

Permite al usuario comunicarse con las facilidades que ofrecen las funcionalidades del _kernel_ mediante una **interfaz en línea de comandos**. Sin embargo, no es una herramienta apropiada para realizar tareas de gestión de contenedores en un alto nivel.

### LMCTFY y Docker (2013)
**LMCTFY** (Let Me Containerize That For You) fue la versión open-source del **stack de Google** a los contenedores de Linux. Su desarrollo **finalizó en 2015** y Google comenzó a transferir parte de la implementación a _libcontainer_ ahora mismo es una parte fundamental del stack de Docker y forma parte de OCI.

Por otro lado, **Docker** supuso un despunte en la popularidad de los contenedores en el año de su salida y el crecimiento de ambos conceptos han ido de la mano desde entonces siendo **hoy en día líder en este ámbito**.

En sus inicios utilizaba LXC, pero más tarde lo sustituyó por su propia librería _libcontainer_.

### OCI (2015)
**OCI** (Open Container Initiative) es un proyecto de la Linux Foundation cuyo objetivo es **diseñar un estándar abierto para la virtualización basada en contenedores**. Fue establecida en 2015 por Docker y otros líderes de la industria.

Después del lanzamiento de Docker, surgió una comunidad alrededor de los contenedores. Sin embargo, con el paso del tiempo fueron apareciendo nuevas tecnologías y herramientas que satisfacían las neceseidades que iban surgiendo. Este fue el motivo principal por el que surgió este estándar.

Actualmente, OCI define dos especificaciones, aunque [hablaremos más tarde de ellas en este documento](#chroot-1979) *******CAMBIAAAAAAAAAR LIIIIINK********

### Otras tecnologías
A lo largo de las dos últimas décadas han ido surgiendo otras tecnologías de virtualización de sistema operativo, pero han sido menos importantes para el ecosistema de los contenedores que las anteriormente mencionadas.

Entre ellas podemos encontrar **Linux VServer** (2001), **OpenVZ** (2005), **Warden** (2011), **Singularity** (2015) o **Podman** (2018).

## Arquitectura de los contenedores
Dada la ambigüedad con la que muchos profesionales se refieren a cada una de las capas que forman esta arquitectura y, dado que dependiendo de la tecnología a usar, pueden cambiar ligeramente las funciones que realizan sus componentes, se utilizará como referencia el stack de Docker.


<figure>
<img src="./img/docker_architecture.png" width="600" />
<figcaption>Generalización de la arquitectura de los contenedores a partir de la arquitectura de Docker</figcaption>
</figure>

Comenzaremos el estudio por la parte inferior de la imágen, los componentes del Linux Kernel, subiendo hasta llegar al _Container Engine_.

### 1º Componentes del Linux Kernel
Existen numerosos componentes y herramientas del kernel que utilizan los contenedores para aislar los procesos. En concreto aquí estudiaremos cinco de ellos, siendo los más importantes los dos primeros:

#### Namespaces
!!! note "Definición de namespaces"
    Proporcionan el **aislamiento entre procesos** mediante la **encapsulación de ciertos recursos** del sistema. De esta forma, hacen creer a los procesos dentro de un contenedor que tienen su propia instancia del recurso, aunque realmente lo están compartiendo.

Esta herramienta impide que, mediante una vulnerabilidad en un contenedor, unos intrusos puedan acceder a la máquina completa comprometiéndola.

La API del kernel que facilita esta característica tiene tres llamadas principales:

| Llamada      | Descripción                          |
| ----------- | ------------------------------------ |
| `clone`       | Crea un proceso hijo al igual que `fork` pero proporcionando más control sobre qué partes del contexto se comparten.  |
| `unshare`       | Crea un nuevo _namespace_ y ejecuta un nuevo proceso dentro de éste. |
| `setns`    | Mueve el hilo actual dentro de un _namespace_ existente |
| `ioctl` | Ofrece información sobre _namespaces_  (la información se maneja desde `/proc/$PID/ns`)|

En la actualidad existen 8 _namespaces_ distintos, cada uno de ellos proporciona **distinto tipo de aislamiento**:



??? summary "Tipos de namespaces"

    === "MNT NS"

        **Mount Namespace**

        Fue el primer _namespace_ implementado. En aquel momento no se sabía que iban a crearse más tipos de _namespaces_ así que su _flag_ correspondiente en los comandos de `clone` y `unshare` es `CLONE_NEWNS`.

        !!! note "Función"
            Proporcionan **aislamiento** a las estructuras de datos que utiliza el sistema para gestionar los **puntos de montaje**. De esta forma, los procesos en distinto _mount namespace_ tienen una visión distinta de la jerarquía de los sistemas de archivos.

        Algunos usos que ofrece son:

        * Cada usuario puede tener su propio `/tmp` para aumentar la seguridad frente a un usuario malicioso.
        * Distintos procesos pueden tener un sistema de archivos raíz (es un concepto parecido a `chroot`).
        * Los puntos de montaje pueden ser privados o compartidos.


    === "UTS NS"

        **UNIX Time-Sharing**

        Desde la implementación del anterior _namespace_ hasta este pasaron 4 años. Este es el más sencillo.

        !!! note "Función"
            **Aísla el _hostname_ y el _domain-name_** del contenedor.

        La _flag_ correspondiente es `CLONE_NEWUTS`.
        


    === "IPC NS"

        **Inter Process Communication**

        !!! note "Función"
            **Aísla cierta la compartición de ciertos recursos que facilitan la comunicación** como la memoria, objetos System V IPC y colas de mensajes POSIX.

        Sin embargo, después de un `clone` se siguen compartiendo señales, sockets, descriptores de archivos o sondeos de memoria (habría que aislarlos con otros _namespaces_).

        Su flag correspondiente es `CLONE_NETIPC`.
    


    === "PID NS"

        **Process ID**

        Los procesos de un sistema pertenecen a un árbol de procesos global visible únicamente por el _host_.

        !!! note "Función"
            Este _namespace_ crea un nuevo **árbol de procesos** propio de cada contenedor.

        Los **procesos** que estén dentro de este _namespace_ **pertenecen**, por tanto, **a dos árboles**, el global y el propio del _namespace_. De esta forma un mismo proceso puede tener varios PIDs dependiendo de en qué árbol se consulte.

        <figure>
            <img src="./img/pid_ns.png" width="500" />
            <figcaption>Ejemplo generalizado de bridge</figcaption>
        </figure>

        El primer PID dentro del nuevo _namespace_ siempre es el 1. Este proceso debería tener unas características únicas para poder funcionar como _init_ y así poder tener más de un proceso en un contenedor.

        Uno de los principales beneficios que presenta es que los contenedores se pueden migrar de una máquina a otra manteniendo el árbol de procesos.

        Su _flag_ es `CLONE_NEWPID`.



    === "NET NS"

        **Network**

        !!! note "Función"
            Proporciona **aislamiento de los recursos de red**, donde cada contenedor puede tener su propio stack de recursos: tablas de enrutado, reglas de _iptables_, _sockets_, etc.

        En las tecnologías actuales de contenedores se suele utilizar el modelo _bridge_. En este modelo el puente, que está en el _namespace_ global, ofrece la conexión a los contenedores mediante una interfaz virtual. A la vez, cada contenedor, en su propio _namespace_, tiene una interfaz virtual conectada a la del puente.

        <figure>
            <img src="./img/net_ns.png" width="500" />
            <figcaption>Diferencia infraestructura: VM y contenedor</figcaption>
        </figure>

        La _flag_ correspondiente es `CLONE_NEWNET`.


    
    === "User NS"

        !!! note "Función"
            Permite que **los procesos puedan creer que están operando como root** dentro de un contenedor, pero fuera del contenedor tienen realmente los privilegios de un usuario común.

        Antes de que existiera este _namespace_ si un proceso tenía permisos de root en cualquier entorno de aislamiento, también lo tenía en el sistema global, lo que era un problema muy grande en lo que a la seguridad se refiere.

        El UID dentro del _namespace_ siempre se corresponderá con otro fuera de él en el host. Este mapeo se puede hacer manualmente y elegir qué UID se selecciona dentro del contenedor para un usuario global.

        Un ejemplo de mapeo entre las dos tablas de usuarios se puede observar en la siguiente imagen.

        <figure>
            <img src="./img/user_ns.png" width="500" />
            <figcaption>Ejemplo de mapeo entre un UID del host con un UID dentro del contenedor</figcaption>
        </figure>

        La _flag_ correspondiente es `CLONE_NETUSER`.


    
    === "Cgroup NS"

        Los _cgroups_ (rupos de control en español) son una funcionalidad del kernel de Linux que trataremos en la sección ********************OJOOOOOOOOOOOOOO************* que, básicamente, implementan monitorización y limitación de los recursos que consume un grupo de procesos.

        !!! note "Función"
            **Virtualiza la vista de los grupos de control** que ven los procesos. Sin esta restricción, un proceso podría observar los grupos de control globales y acceder a la información sobre la limitación de recursos de otros procesos.

        Tiene dos funciones principales:

        1. Evita que un grupo de control acceda a los límites de otro superior a él.
        2. Facilita la migración de contenedores al aislar unos grupos de otros. De esta forma podemos ahorrar la necesidad de replicar los límites según se van trasladando los contenedores.

        La _flag_ correspondiente es `CLONE_NEWCGROUP`.



    === "Time NS"

        !!! note "Función"
            **Virtualiza dos relojes del sistema** permitiendo que, cuando un contenedor se migra de un nodo a otro, estos **relojes son restaurados de forma consistente** partiendo siempre del tiempo que tenía el reloj antes de ser migrado.

        Estos dos relojes son el `CLOCK_MONOTONIC` y `CLOCK_BOOTTIME`.

        La _flag_ correspondiente es `CLONE_NEWTIME`.



#### Control Groups

Antes de que se crearan los grupos de control existían formas para **monitorizar y controlar procesos individuales**, pero **no había soporte** para aplicar esas mismas operaciones **a grupos de procesos**. Por eso se acabaron implementando los _Process Containers_, que luego se renombraron a _cgroups_.

!!! note "Definición de Cgroups"

    Son un **_framework_** basado en los mecanismos existentes del kernel para **proporcionar una interfaz y un alcance más global a las operaciones de control y monitorización** de procesos.

    Los grupos de control permiten **repartir y asignar recursos entre un grupo de tareas o procesos** del sistema.


Los grupos de control ofrecen cuatro características principales:

1. **Limitación de recursos**, como procesador, memoria, dispositivos E/S, etc.
2. **Priorización**, que es parecida a la anterior característica, pero no limita recursos a procesos sino que le da preferencia de consumo de recursos a un proceso seleccionado.
3. **Monitorización** de un grupo de procesos para obtener inforamción de qué procesos están consumiendo cuántos recursos.
3. **Freezing de procesos**. Esta es una herramienta que permite paralizar y retomar grupos de procesos. Es muy utilizada en el procesamiento por lotes.

Existen dos versiones de **Cgroups**, hoy en día se está intentando hacer la transición en la mayoría de herramientas a la segunda versión. Para comprender esta última versión necesitamos entender tres conceptos:

* **Cgroup**: es un grupo de tareas al que se le asocia uno o más subsistemas.
* **Subsistema o _resource controller_**: representa un único recurso, como la memoria o el tiempo de CPU.
* **Jerarquía**: es el conjunto de todos los _cgroups_ en forma de árbol. Cada proceso del sistema está en un _cgroup_ determinado. Cada nodo (o sea, cada _cgroup_) del árbol tiene asociados uno o más subsistemas. En la versión 2 de los _cgroups_ existe una sóla jerarquía.

En la siguiente imagen podemos observar un ejemplo de jerarquía, donde existen cuatro _cgroups_ distintos. Cada uno de ellos tiene asociados varios subsistemas. En el caso del Group_2 los subsistemas asociados son "memoria" y "pids", y sus grupos hijos heredan los mismos subsistemas, pero pueden asociar otros nuevos o disociar los heredados.

<figure>
    <img src="./img/cgroup_jerar.png" width="700" />
    <figcaption>Ejemplo de jerarquía de grupos de control.</figcaption>
</figure>

Cada grupo se añade creando un nuevo directorio dentro de la carpeta raíz `/sys/fs/cgroup`. En el momento de crear dicha carpeta, el sistema añade los archivos necesarios para gestionar los _cgroups_.

??? note "¿Cómo añadimos los grupos de control de este ejemplo?"
    Los subsistemas del siguiente ejemplo se añadirían creando en primer lugar cada _cgroup_ y más tarde escribiendo en los subsistemas en los archivos autogenerados:

    ```bash
    mkdir /sys/fs/cgroup/Cgroup_1
    mkdir /sys/fs/cgroup/Cgroup_2
    mkdir /sys/fs/cgroup/Cgroup_2.1
    mkdir /sys/fs/cgroup/Cgroup_2.2

    echo "+memory +io" > /sys/fs/cgroup/Group_1/cgroup.subtree_control
    echo "+memory +pids" > /sys/fs/cgroup/Group_2/cgroup.subtree_control
    echo "+io" > /sys/fs/cgroup/Group_2.1/cgroup.subtree_control
    echo "-pids" > /sys/fs/cgroup/Group_2.2/cgroup.subtree_control
    ```


#### Union Filesystem

Los sistemas de archivos por capas permiten compartir archivos en el disco, lo que supone un ahorro de espacio.

!!! note "Definición"
    Los **Union Filesystem** son un tipo de sistema de archivos por capas. Permite a los archivos y directorios de sistemas de archivos distintos estar **superpuestos formando un único sistema de archivos**.

##### OverlayFS

Docker, actualmente, utiliza un tipo de Union Filesystem llamado Overlay2. Esta tecnología utiliza tres capas:

1. **Base** (sólo lectura): es la capa donde van los archivos base. En terminología de Docker se correspondería con la imagen.
2. **Overlay**: es la capa donde opera el usuario. En un principio ofrece una vista de la capa base y permite operar sobre ella, aunque los cambios no se guardan en esta capa.
3. **Diff**: es la capa donde se guardan los cambios realizados en la capa anterior.

Una de las características fundamentales que usan los UnionFS es la **técnica de COW** (_Copy On Write_) que permite **reducir el consumo de copias sin modificar**. Si un fichero existe en la capa **base** y otra capa quiere leerlo, lee el original. Cuando una capa quiera modificar el archivo es cuando debe copiarla y modificarla en la propia capa que la necesite.

<figure>
    <img src="./img/overlay.png" width="700" />
    <figcaption>Las tres capas de Overlay2 junto a la técnica COW.</figcaption>
</figure>


#### Capabilities

El usuario con UID 0 es el root y tiene el control completo del sistema. Esto puede traer problemas recurrentes de seguridad (las Jails ya intentaron confinar el root omnipotente). Por eso, las _capabilities_ proponen una solución.

!!! note "Definición"
    Las **_capabilities_** dividen los permisos que tiene un usuario en varias particiones (cada una de ellas es una _capability_).

Un buen ejemplo de su uso podría ser la posibilidad de permitir a un binario poder crear un _raw socket_ como en el caso de la instrucción `ping` mediante la _capability_ `CAP_NET_RAW` sin necesidad de asignarle el resto de privilegios que ofrece el usuario root.

Existen dos tipos de _capabilities_, unas son las asociadas a procesos y otros asociadas a archivos. Las que están asociadas a archivos se unen con las asociadas a procesos cuando un proceso quiere ejecutar un archivo.

!!! hint ""
    En la práctica, los contenedores no necesitan todos los privilegios que ofrece el root. Así que, realmente, los usuarios root dentro de los contenedores tienen asignadas algunas _capabilities_ para ofrecer únicamente ciertos permisos y restringir otros potencialmente inseguros**.

#### Pivot_root

En el capítulo *********+OJOOOOOOOOOOO******** ya explicamos la importancia que tuvo `chroot` para aparentar que cambiaba el directorio raíz de un proceso y de sus hijos. El problema que presenta es que no es segura, así que es mejor utiliza la alternativa _pivot\_root_.

La función completa de **pivot_root** es la siguiente:

```bash
int pivot_root (const char *new_root, const char *put_old);
```

Pivot_root **mueve el punto de montaje raíz al directorio put_old** haciéndolo inaccesible para el proceso llamador como a todos sus hijos, mientras que convierte a **new_root en el nuevo punto de montaje raíz**.

!!! hint ""
    Esto soluciona las brechas de seguridad de `chroot` porque se aplica al _mount namespace_, así que está cambiando únicamente en el _namespace_ del contenedor la posibilidad de acceder a la antigua raíz.



### 2º Container Runtimes

!!! note "¿Qué es un container runtime?"
    Un _container runtime_ es la herramienta o capa responsable de que el contenedor se ejecute correctamente (sin incluir los procesos que están dentro del contenedor).

Este término puede tener diferentes significados dependiendo del proyecto o comunidad donde se consulte, principalmente debido a que el rango de tareas que realiza un _runtime_ no está plenamente definido. En este documento vamos a **separar los _container runtimes_ según las tareas que realicen** en _low-level_ y _high-level_.

Algunos ejemplos de **low-level** pueden ser **runc, crun, gvisor o kata-runtime**, mientras que entre los ejemplos de **high-level** podemos encontrarnos con **containerd o CRI-O**.

??? info "¿Y Docker? ¿Dónde lo metemos?"
    En muchas referencias se trata a Docker como un _container runtime_ (refiriéndose a Docker con su significado de _daemon_ más su CLI). Llegados a este punto podemos entender perfectamente que se puede considerar de esta forma, ya que su función es gestionar los contenedores y las imágenes (definición de _high-level container runtime_). Sin embargo, en este documento se va a colocar en un escalón superior, en los _container engines_.

    Esta decisión se ha tomado en base a dos razones: la primera es que una de las características principales que hacen a esta herramienta más popular que containerd o CRI-O entre los usuarios es su capacidad de facilitar la comunicación con el usuario. La segunda razón es que delega gran parte del trabajo en otro _high-level runtime_, que por defecto es containerd. Dicho esto, reiteramos que **podría ser perfectamente considerado un _container runtime_**.

    De la misma forma, CRI-O o containerd se podrían llegar a considerar _container engines_ si se utilizan a través de plugins CLI que permitan al usuario interactuar debidamente.

    **En resumen, la forma de referirse a estas tecnologías no es universal y, como hemos dicho, dependiendo de la comunidad y el proyecto, las etiquetas pueden cambiar**.

    Una forma de entender gráficamente el porqué de que estos términos sean tan ambiguos se puede observar en la siguiente figura, que muestra de forma muy general y subjetiva el nivel al que operan algunos ejemplos de _container runtimes_.


    <figure>
        <img src="./img/runtimes.png" width="700" />
        <figcaption>Nivel al que opera cada container runtime</figcaption>
    </figure>

#### Low-level container runtimes

En esta página nos referiremos a los **_low-level runtimes_** para hablar de las capas que proporcionan las utilidades básicas como crear _namespaces_ y comenzar el proceso de encapsular una aplicación en un contenedor, es decir a las herramientas que se **comunican directamente con el kernel**.

##### runc

!!! note "Definición"
    Runc es un CLI que se encarga de crear contenedores y ejecutarlos según la especificación que se le proporcione. Esta especificación sigue el estándar OCI.

La razón de que sea un _low-level container runtime_ cuando realmente sólo es un CLI es que antiguamente Docker utilizaba como _container runtime_ LXC. Esto cambió más adelante y comenzó a crear junto con Google una librería propia que sustituyera a LXC: **_libcontainer_**. Hoy en día, runc utiliza _libcontainer_, de hecho es su pieza básica y fundamental, por eso se pueden considerar que ambos son la misma capa.

??? info "¿Qué estándar OCI?"
    Los estándares OCI se introdujeron en ******OJOOOOOOOOO********* indicando la existencia de dos estándares: de imágen y _runtime_. Runc sólo implementa la segunda especificación ya que no entiende de imágenes, sino de _bundles_.

    Un _bundle_ es un **conjunto de archivos que contiene todos los datos necesarios** para que un _runtime_ pueda realizar todas las operaciones necesarias para **crear y ejecutar un contenedor**.

