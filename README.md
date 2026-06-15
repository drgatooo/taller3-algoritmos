# Taller 3 — Algoritmos y Estructuras de Datos

Integrantes

- Anchi Cristobal Ernesto Alonso
- Cesar Gomez Chavez
- Contreras Salcedo Maximo Simon

### Video Explicativo

Enlace al video: https://youtu.be/r3se1OOeAQ4

## Ejercicio 4.2 — Rate Limiter con Cola FIFO

### Objetivo

Implementar un sistema de limitación de peticiones (Rate Limiter) utilizando una estructura de datos Cola FIFO implementada manualmente mediante una lista enlazada.

El sistema procesa registros de acceso de un servidor web y restringe la cantidad de solicitudes permitidas por dirección IP dentro de una ventana de tiempo configurable.

**Dataset Utilizado:**

Web Server Access Logs (Kaggle). Autor: eliasdabbas

https://www.kaggle.com/datasets/eliasdabbas/web-server-access-logs

El dataset contiene millones de registros de acceso de un servidor web en formato Apache Log.

Debido a su tamaño (aproximadamente 3.5 GB), el archivo no se incluye en el repositorio y debe descargarse manualmente.

Ubicación esperada: `access.log`

Para pruebas rápidas también se utilizó: `access_sample.log`
generado a partir de las primeras líneas del dataset original.

### Estructura del Proyecto

```
taller3/
│
├── go.mod
├── README.md
│
├── access.log
├── access_sample.log
│
└── colas/
    ├── cola.go
    ├── cola_test.go
    │
    └── cmd/
        └── main.go
```

### Tecnologías Utilizadas

- Go 1.26.2
- PowerShell
- Visual Studio Code
- Dataset Kaggle

### Implementación

#### Cola FIFO

Se implementó una cola utilizando una lista enlazada simple. Cada nodo almacena:

```go
type nodo struct {
    ts int64
    next *nodo
}
```

La estructura Cola mantiene referencias al frente y al final:

```go
type Cola struct {
    head *nodo
    tail *nodo
    len  int
}
```

Operaciones implementadas:

- Enqueue()
- Dequeue()
- Front()
- Len()

Todas con complejidad O(1).

#### Funcionamiento del Rate Limiter

Para cada dirección IP se mantiene una cola independiente de timestamps.

```go
map[string]*Cola
```

Al llegar una nueva petición:

1. Se eliminan timestamps fuera de la ventana de tiempo.
2. Se verifica cuántas peticiones permanecen activas.
3. Si existen menos de M peticiones se acepta.
4. Si existen M o más peticiones se rechaza.

Esto permite implementar una ventana deslizante eficiente.

### Complejidad Temporal

#### Cola

| Operación | Complejidad |
| --------- | ----------- |
| Enqueue   | O(1)        |
| Dequeue   | O(1)        |
| Front     | O(1)        |
| Len       | O(1)        |

#### Rate Limiter

Complejidad amortizada: `O(1)`

**Justificación:**

Cada timestamp:

- se inserta exactamente una vez
- se elimina exactamente una vez

Por lo tanto el costo total para n peticiones es O(n).

### Ejecución

Desde la raíz del proyecto:

```shell
go run ./colas/cmd access_sample.log 10 60
```

o utilizando el dataset completo:

```shell
go run ./colas/cmd access.log 10 60
```

Parámetros:

- `access.log` → archivo de entrada
- `10` → máximo de peticiones permitidas
- `60` → ventana de tiempo en segundos

### Pruebas Unitarias

**Ejecución:**

```shell
go test ./colas -v
```

**Resultados:**

```
PASS
ok      taller3/colas
```

**Pruebas realizadas:**

- TestColaFIFO
- TestColaVacia
- TestColaUnElemento
- TestFront
- TestPermitirPeticionLimite
- TestVentanaDeslizante
- TestIPsIndependientes
- TestParsearLineaOK
- TestParsearLineaInvalida

Todas las pruebas fueron superadas satisfactoriamente.

### Benchmarks

**Ejecución:**

```shell
go test ./colas -bench=Benchmark -benchmem
```

**Resultados obtenidos:**

```
goos: windows
goarch: amd64
pkg: taller3/colas
cpu: AMD Ryzen 5 2600 Six-Core Processor

BenchmarkEnqueueDequeue-12
37039779
27.72 ns/op
16 B/op
1 allocs/op

BenchmarkPermitirPeticion-12
20256036
54.91 ns/op
16 B/op
1 allocs/op
```

**Interpretación:**

- Las operaciones de cola presentan tiempos constantes.
- El Rate Limiter mantiene comportamiento eficiente incluso con millones de operaciones.
- El consumo de memoria por operación es mínimo.

### Conclusiones

- Se implementó exitosamente una Cola FIFO mediante lista enlazada.
- El sistema de Rate Limiter procesa correctamente registros reales de acceso web.
- Las pruebas unitarias validan el correcto funcionamiento de la solución.
- Los benchmarks muestran tiempos de ejecución compatibles con una complejidad O(1).
- La solución escala adecuadamente para archivos de gran tamaño gracias al procesamiento secuencial línea por línea.
