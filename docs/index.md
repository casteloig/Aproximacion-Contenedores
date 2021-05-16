# ¿Qué es un contenedor?

La virtualización es un proceso mediante el cual un software es usado para crear una abstracción sobre unos recursos, dando la sensación de que los elementos hardware se dividen en varias computadores virtuales.

Existen dos tipos técnicas principales de virtualización:

1. Máquinas virtuales.
2. Contenedores (virtualización ligera).

!!! note "Definición"
    Un **contenedor** es un límite lógico que se crea dentro de un sistema operativo proporcionado por el aislamiento de recursos hardware.


    La característica principal es que en esta técnica se utilizan herramientas que proporciona el Linux Kernel (como _cgroups_ y _namespaces_).

## Diferencia entre MV y contenedores

*****IMAGEN*****

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

        1. **Rapidez y ligereza**.
        2. **Portabilidad e independencia de plataformas**.
        3. **Escalabilidad**.



Esto es una prueba del actions2
