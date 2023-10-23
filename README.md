# Grupo27-Laboratorio-2

# **Integrantes**

---

- Nicolas Barrera 201951552-7

- Daniel  Sepúlveda 201873065-3

- Javier Maturana 201873604-k

# **Instrucciones**

---

**Primero es necesario para que todas las VM esten en la 
ruta correcta es necesrio ejecutar el siguiente comando en
todas las VM:**


```
cd Grupo27-Laboratorio-2
```

1. Ejecutar en dist105:
```
make docker-OMS
```

2. En la dist107 y dist108 se debe ejecutar:

```
make docker-datanode
```

3. Despues solo en dist106 se debe ejecutar:

```
make docker-ONU
```

4. Finalmente en todas las VM se debe ejecutar: 

```
make docker-continentes
```

5. Para revisar los log de los dockers que se ejecutan en segundo plano:

```

En dist107 y 108:

    docker logs regional
    docker logs datanode

En dist106:

    docker logs regional

en dist105:

    docker logs regional
    docker logs oms
```
6. Para detener programa:
```

En dist107 y 108:

    docker stop regional
    docker stop datanode

En dist106:

    docker stop regional
    docker stop onu

en dist105:

    docker stop regional
    docker stop oms
```
7. Para limpiar archivos .txt
```
    make clean
```
