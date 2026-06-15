# Taller 3 — Algoritmos y Estructuras de Datos

Integrantes

- Anchi Cristobal Ernesto Alonso
- Cesar Gomez Chavez
- Contreras Salcedo Maximo Simon

### Videos Explicativos

- Ejercicio 1: pendiente
- Ejercicio 2: https://youtu.be/r3se1OOeAQ4
- Ejercicio 3: pendiente
- Ejercicio 4: pendiente

---

## Ejercicio 1 — Stock Span con Pila Monótona

### Objetivo

Implementar una pila genérica y calcular el **Stock Span** de una serie de precios de cierre bursátiles utilizando una pila monótona de índices.

El *stock span* de un día indica cuántos días consecutivos hacia atrás (incluido el actual) el precio de cierre fue menor o igual al precio del día actual.

### Dataset

**Huge Stock Market Dataset** — Kaggle. Autor: Boris Marjanovic

https://www.kaggle.com/datasets/borismarjanovic/price-volume-data-for-all-us-stocks-etfs

Cada archivo (p. ej. `aapl.us.txt`) contiene columnas `Date, Open, High, Low, Close, Volume, OpenInt`. Se utiliza la columna `Close` de una acción a elección.

Ubicación esperada: `data/<accion>.us.txt`

### Estructura del Proyecto

```
taller3/
│
├── go.mod
├── README.md
│
├── data/
│   └── <accion>.us.txt   ← dataset bursátil (no incluido, descargar de Kaggle)
│
└── pilas/
    ├── pila.go           ← Pila genérica + CalcularStockSpan + lector CSV
    ├── pila_test.go      ← pruebas unitarias y benchmarks
    │
    └── cmd/
        └── main.go       ← punto de entrada
```

### Implementación

> **Pendiente** — sección en desarrollo.

### Complejidad Temporal

| Operación          | Complejidad |
| ------------------ | ----------- |
| `Push`             | O(1)        |
| `Pop`              | O(1)        |
| `Peek`             | O(1)        |
| `IsEmpty`          | O(1)        |
| `CalcularStockSpan`| O(n)        |

**Justificación:** Cada índice se apila exactamente una vez y se desapila exactamente una vez, por lo que el costo total del algoritmo es O(n) pese a tener un bucle interno.

### Ejecución

```shell
go run ./pilas/cmd <ruta_archivo> <N>
```

**Parámetros:**

| Argumento       | Descripción                                      |
| --------------- | ------------------------------------------------ |
| `<ruta_archivo>`| Ruta al archivo `.txt` de la acción              |
| `<N>`           | Número de primeros/últimos resultados a mostrar  |

### Pruebas Unitarias

**Ejecución:**

```shell
go test ./pilas -v
```

**Resultados:**

> Pendiente.

**Pruebas previstas:**

| Test                     | Descripción                                            |
| ------------------------ | ------------------------------------------------------ |
| `TestPushPop`            | Verifica orden LIFO al extraer elementos               |
| `TestPopVacia`           | Pop en pila vacía devuelve `false`                     |
| `TestPeek`               | Peek no modifica el tamaño de la pila                  |
| `TestStockSpanBasico`    | Span correcto sobre serie de precios conocida          |
| `TestStockSpanCreciente` | Serie ordenada ascendentemente → todos los spans crecen|
| `TestStockSpanDecreciente`| Serie ordenada descendentemente → todos los spans = 1 |
| `TestLeerPrecios`        | Carga y parseo correcto del archivo CSV                |

### Benchmarks

**Ejecución:**

```shell
go test ./pilas -bench=Benchmark -benchmem
```

**Resultados:**

> Pendiente.

### Conclusiones

> Pendiente.

---

## Ejercicio 2 — Rate Limiter con Cola FIFO

### Objetivo

Implementar un sistema de limitación de peticiones (Rate Limiter) utilizando una estructura de datos Cola FIFO implementada manualmente mediante una lista enlazada.

El sistema procesa registros de acceso de un servidor web y restringe la cantidad de solicitudes permitidas por dirección IP dentro de una ventana de tiempo configurable.

### Dataset

**Web Server Access Logs** — Kaggle. Autor: eliasdabbas

https://www.kaggle.com/datasets/eliasdabbas/web-server-access-logs

El dataset contiene millones de registros de acceso de un servidor web en formato Apache Common Log (`IP - - [fecha] "método recurso protocolo" código bytes`).

Debido a su tamaño (~3.5 GB), el archivo **no se incluye** en el repositorio y debe descargarse manualmente.

Ubicación esperada: `access.log`

Para pruebas rápidas también se utilizó `access_sample.log`, generado a partir de las primeras líneas del dataset original.

### Estructura del Proyecto

```
taller3/
│
├── go.mod
├── README.md
│
├── access.log          ← dataset completo (no incluido, descargar de Kaggle)
├── access_sample.log   ← muestra reducida para pruebas
│
└── colas/
    ├── cola.go         ← Cola FIFO + Rate Limiter + parser de log
    ├── cola_test.go    ← pruebas unitarias y benchmarks
    │
    └── cmd/
        └── main.go     ← punto de entrada
```

### Implementación

#### Cola FIFO

Se implementó una cola mediante una lista enlazada simple. Cada nodo almacena únicamente un timestamp Unix:

```go
type nodo struct {
    ts   int64
    next *nodo
}
```

La estructura `Cola` mantiene punteros al frente (elemento más antiguo) y al final (elemento más nuevo):

```go
type Cola struct {
    head *nodo
    tail *nodo
    len  int
}
```

Operaciones implementadas con complejidad **O(1)**:

| Operación   | Descripción                                      |
| ----------- | ------------------------------------------------ |
| `Enqueue`   | Agrega un timestamp al final de la cola          |
| `Dequeue`   | Extrae el timestamp del frente                   |
| `Front`     | Consulta el frente sin modificar la cola         |
| `Len`       | Devuelve la cantidad de elementos                |

#### Funcionamiento del Rate Limiter

Para cada dirección IP se mantiene una cola independiente de timestamps en un mapa:

```go
map[string]*Cola
```

Al llegar una nueva petición con IP `ip` y timestamp `ts`:

1. Se eliminan del frente los timestamps fuera de la ventana `[ts - T, ts]`.
2. Se verifica cuántas peticiones permanecen activas en la cola.
3. Si `cola.Len() < M` → se acepta la petición y se registra el timestamp.
4. Si `cola.Len() >= M` → se rechaza la petición.

Esto implementa una **ventana deslizante** sin necesidad de recorrer la cola completa.

#### Parseo del Log

```go
func ParsearLinea(linea string) (ip string, ts int64, err error)
```

Extrae la IP y el timestamp de cada línea en formato Apache Common Log. Devuelve `error` para líneas malformadas, que se omiten silenciosamente durante el procesamiento.

### Complejidad Temporal

#### Cola

| Operación | Complejidad |
| --------- | ----------- |
| Enqueue   | O(1)        |
| Dequeue   | O(1)        |
| Front     | O(1)        |
| Len       | O(1)        |

#### Rate Limiter

Complejidad amortizada: **O(1)** por petición.

**Justificación:** Cada timestamp se inserta exactamente una vez (`Enqueue`) y se elimina exactamente una vez (`Dequeue`). El costo total sobre *n* peticiones es O(n), por lo que el costo amortizado por petición es O(1).

### Ejecución

Desde la raíz del proyecto:

```shell
go run ./colas/cmd access_sample.log 10 60
```

O con el dataset completo:

```shell
go run ./colas/cmd access.log 10 60
```

**Parámetros:**

| Argumento         | Descripción                              |
| ----------------- | ---------------------------------------- |
| `access.log`      | Ruta al archivo de log de entrada        |
| `10`              | Máximo de peticiones permitidas (M)      |
| `60`              | Ventana de tiempo en segundos (T)        |

**Salida esperada:** decisión (`ACEPTADA` / `RECHAZADA`) por petición (muestra) + resumen con total de rechazos y las 5 IPs con más rechazos.

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

| Test                      | Descripción                                            |
| ------------------------- | ------------------------------------------------------ |
| `TestColaFIFO`            | Verifica orden FIFO al extraer elementos               |
| `TestColaVacia`           | Dequeue y Front en cola vacía devuelven `false`        |
| `TestColaUnElemento`      | Caso límite: un solo elemento en la cola               |
| `TestFront`               | Front no modifica el tamaño de la cola                 |
| `TestPermitirPeticionLimite` | La petición M+1 es rechazada correctamente          |
| `TestVentanaDeslizante`   | Timestamps expirados liberan cupo en la ventana        |
| `TestIPsIndependientes`   | Cada IP mantiene su propia cola de forma independiente |
| `TestParsearLineaOK`      | Línea Apache válida se parsea correctamente            |
| `TestParsearLineaInvalida`| Línea malformada devuelve error                        |

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

BenchmarkEnqueueDequeue-12    37039779    27.72 ns/op    16 B/op    1 allocs/op
BenchmarkPermitirPeticion-12  20256036    54.91 ns/op    16 B/op    1 allocs/op
```

**Interpretación:**

- Las operaciones de cola presentan tiempos prácticamente constantes, consistentes con O(1).
- El Rate Limiter mantiene comportamiento eficiente incluso con millones de operaciones.
- El consumo de memoria por operación es mínimo: 1 allocación de 16 bytes por llamada.

### Conclusiones

- Se implementó exitosamente una Cola FIFO mediante lista enlazada sin usar librerías externas.
- El Rate Limiter procesa correctamente registros reales de acceso web con política de ventana deslizante.
- Las pruebas unitarias cubren casos normales, límite y de error.
- Los benchmarks confirman tiempos compatibles con una complejidad O(1) amortizada.
- La solución escala adecuadamente para archivos de gran tamaño al procesar el log línea por línea.

---

## Ejercicio 3 — Caché LRU con Lista Doblemente Enlazada

### Objetivo

Implementar una caché **LRU (Least Recently Used)** utilizando una lista doblemente enlazada combinada con un mapa hash para lograr operaciones `Get` y `Put` en tiempo **O(1)**.

La caché se evalúa sobre una traza real de accesos a películas, calculando el **hit ratio** para distintos tamaños de caché.

### Dataset

**MovieLens** — GroupLens Research

https://grouplens.org/datasets/movielens/

Se utiliza el archivo `ratings.csv` con columnas `userId, movieId, rating, timestamp`. Los registros se ordenan por `timestamp` y se toma la columna `movieId` como la secuencia de accesos a simular.

El archivo **no se incluye** en el repositorio por su tamaño. Debe descargarse y colocarse en:

```
data/ratings.csv
```

Se recomienda la versión **100K** o **1M** de MovieLens.

### Estructura del Proyecto

```
taller3/
│
├── go.mod
├── README.md
│
├── data/
│   └── ratings.csv     ← dataset MovieLens (no incluido, descargar de GroupLens)
│
└── listas/
    ├── lru.go          ← Lista doblemente enlazada + LRU + carga del CSV
    ├── lru_test.go     ← pruebas unitarias y benchmarks
    │
    └── cmd/
        └── main.go     ← punto de entrada
```

### Implementación

#### Lista Doblemente Enlazada

Cada nodo almacena la clave, el valor y punteros al nodo anterior y siguiente:

```go
type Nodo struct {
    clave int
    valor int
    prev  *Nodo
    next  *Nodo
}
```

La lista mantiene punteros `head` (más recientemente usado) y `tail` (menos recientemente usado). Las operaciones internas de la lista son:

- `eliminar(node)` — desvincula un nodo en O(1).
- `agregarAlFrente(node)` — inserta un nodo como cabeza en O(1).

#### Caché LRU

```go
type LRU struct {
    cap  int
    mapa map[int]*Nodo
    head *Nodo
    tail *Nodo
}
```

El mapa `mapa` permite acceso directo a cualquier nodo en O(1), mientras que la lista doblemente enlazada mantiene el orden de uso.

**`Get(clave int) (int, bool)`**

1. Si la clave no existe en el mapa → devuelve `(0, false)` (cache miss).
2. Si existe → mueve el nodo al frente de la lista (marcándolo como recién usado) y devuelve el valor.

**`Put(clave, valor int)`**

1. Si la clave ya existe → actualiza el valor y mueve el nodo al frente.
2. Si es nueva y la caché está llena → elimina el nodo `tail` (menos recientemente usado) del mapa y de la lista, luego inserta el nuevo nodo al frente.
3. Si es nueva y hay espacio → inserta el nuevo nodo al frente.

#### Carga del Dataset

```go
func CargarSecuencia(ruta string) ([]int, error)
```

Lee `ratings.csv` con la librería estándar `encoding/csv`, carga todos los registros, los ordena por `timestamp` ascendente y devuelve la secuencia de `movieId` en ese orden.

### Complejidad Temporal

| Operación          | Complejidad |
| ------------------ | ----------- |
| `Get`              | O(1)        |
| `Put`              | O(1)        |
| `eliminar`         | O(1)        |
| `agregarAlFrente`  | O(1)        |
| `CargarSecuencia`  | O(n log n)  |

**Justificación:** El mapa hash garantiza acceso O(1) a cualquier nodo. La lista doblemente enlazada permite reordenar nodos en O(1) al tener punteros directos `prev` y `next`. La evicción del nodo `tail` también es O(1).

### Ejecución

Desde la raíz del proyecto:

```shell
go run ./listas/cmd/main.go ./data/ratings.csv 50 100 500 1000
```

**Parámetros:**

| Argumento            | Descripción                                        |
| -------------------- | -------------------------------------------------- |
| `./data/ratings.csv` | Ruta al archivo CSV de MovieLens                   |
| `50 100 500 1000`    | Tamaños de caché a evaluar (uno o más valores)     |

**Salida esperada:**

```
=======================================================
Capacidad       Total Accesos   Hits       Hit Ratio
=======================================================
50              100000          12345      12.35%
100             100000          18900      18.90%
500             100000          45321      45.32%
1000            100000          61200      61.20%
```

### Pruebas Unitarias

**Ejecución:**

```shell
go test ./listas -v
```

**Resultados:**

```
=== RUN   TestNewLRU
--- PASS: TestNewLRU (0.00s)
=== RUN   TestCacheHit
--- PASS: TestCacheHit (0.00s)
=== RUN   TestUpdateExistingKey
--- PASS: TestUpdateExistingKey (0.00s)
=== RUN   TestEviccionLRU
--- PASS: TestEviccionLRU (0.00s)
=== RUN   TestCacheMiss
--- PASS: TestCacheMiss (0.00s)
=== RUN   TestPutUpdateKeyNoHead
--- PASS: TestPutUpdateKeyNoHead (0.00s)
=== RUN   TestGetItemIsHead
--- PASS: TestGetItemIsHead (0.00s)
=== RUN   TestCargarSecuencia
--- PASS: TestCargarSecuencia (0.00s)
PASS
ok      taller3/listas  0.002s
```

**Pruebas realizadas:**

| Test                      | Descripción                                                  |
| ------------------------- | ------------------------------------------------------------ |
| `TestNewLRU`              | Verifica la inicialización correcta de la estructura         |
| `TestCacheHit`            | Get devuelve el valor correcto y actualiza el orden          |
| `TestUpdateExistingKey`   | Actualizar una clave existente modifica el valor correctamente |
| `TestEviccionLRU`         | Al superar la capacidad, el elemento LRU es expulsado        |
| `TestCacheMiss`           | Get de clave inexistente devuelve `(0, false)`               |
| `TestPutUpdateKeyNoHead`  | Actualizar una clave no-cabeza la mueve al frente            |
| `TestGetItemIsHead`       | Get de la cabeza no altera la estructura                     |
| `TestCargarSecuencia`     | La función ordena los IDs por timestamp correctamente        |

Todas las pruebas fueron superadas satisfactoriamente.

### Benchmarks

**Ejecución:**

```shell
go test ./listas -bench=Benchmark -benchmem
```

**Resultados obtenidos:**

```
goos: linux
goarch: amd64
pkg: taller3/listas
cpu: AMD Ryzen 5 PRO 8540U w/ Radeon 740M Graphics

BenchmarkLRUPut-12    13626411    74.91 ns/op    32 B/op    1 allocs/op
BenchmarkLRUGet-12    129905391    9.368 ns/op    0 B/op    0 allocs/op
PASS
ok      taller3/listas  3.263s
```

Los benchmarks `BenchmarkLRUPut` y `BenchmarkLRUGet` miden respectivamente la velocidad de inserción con evicción y la velocidad de lectura con cache hit, ambos sobre una caché de capacidad 500.

**Interpretación:**

- Tanto `Get` como `Put` exhiben tiempos constantes independientemente del número de operaciones.
- El acceso O(1) garantizado por el mapa hash evita recorridos lineales.
- La evicción es igualmente eficiente al eliminar directamente el nodo `tail`.

### Conclusiones

- Se implementó exitosamente una caché LRU combinando lista doblemente enlazada y mapa hash, sin librerías externas.
- Las operaciones `Get` y `Put` son O(1), incluyendo la evicción del elemento menos recientemente usado.
- La simulación sobre MovieLens muestra que el hit ratio crece al aumentar la capacidad de la caché, tendencia esperada para distribuciones de acceso con localidad temporal.
- Las pruebas unitarias cubren los casos normales, límite y de error exigidos por la rúbrica.
- Los benchmarks confirman el comportamiento O(1) empíricamente.

---

## Ejercicio 4 — Índice AVL con Consultas por Rango

### Objetivo

Implementar un árbol **AVL** (árbol binario de búsqueda autobalanceado) que garantice altura O(log n) e indexe películas de MovieLens por clave numérica, permitiendo consultas por rango `[a, b]` eficientes.

### Dataset

**MovieLens** — GroupLens Research

https://grouplens.org/datasets/movielens/

Se utilizan `movies.csv` y `ratings.csv`. La clave de indexación es el **rating promedio** por película (calculado agregando `ratings.csv`) o el **año** extraído del título en `movies.csv`.

Ubicación esperada: `data/movies.csv` y `data/ratings.csv`

### Estructura del Proyecto

```
taller3/
│
├── go.mod
├── README.md
│
├── data/
│   ├── movies.csv        ← dataset MovieLens (no incluido, descargar de GroupLens)
│   └── ratings.csv       ← dataset MovieLens (no incluido, descargar de GroupLens)
│
└── arboles/
    ├── avl.go            ← NodoAVL + rotaciones + Insertar + ConsultaRango
    ├── avl_test.go       ← pruebas unitarias y benchmarks
    │
    └── cmd/
        └── main.go       ← punto de entrada
```

### Implementación

> **Pendiente** — sección en desarrollo.

### Complejidad Temporal

| Operación        | Complejidad       |
| ---------------- | ----------------- |
| `Insertar`       | O(log n)          |
| `ConsultaRango`  | O(log n + k)      |
| Rotaciones (×4)  | O(1) c/u          |

**Justificación:** El balanceo automático (factor |FB| ≤ 1 en cada nodo) garantiza que la altura del árbol sea O(log n) incluso si los datos se insertan en orden. La consulta por rango poda subárboles completos fuera del intervalo `[a, b]`, logrando O(log n + k) donde k es el número de resultados.

### Ejecución

```shell
go run ./arboles/cmd <ruta_movies> <ruta_ratings> <a> <b>
```

**Parámetros:**

| Argumento        | Descripción                                    |
| ---------------- | ---------------------------------------------- |
| `<ruta_movies>`  | Ruta a `movies.csv`                            |
| `<ruta_ratings>` | Ruta a `ratings.csv`                           |
| `<a>`            | Límite inferior del rango de consulta          |
| `<b>`            | Límite superior del rango de consulta          |

### Pruebas Unitarias

**Ejecución:**

```shell
go test ./arboles -v
```

**Resultados:**

> Pendiente.

**Pruebas previstas:**

| Test                        | Descripción                                              |
| --------------------------- | -------------------------------------------------------- |
| `TestInsercionSimple`       | Inserción de nodos y verificación del BST               |
| `TestBalanceoLL`            | Rotación LL restaura el balance                          |
| `TestBalanceoRR`            | Rotación RR restaura el balance                          |
| `TestBalanceoLR`            | Rotación LR restaura el balance                          |
| `TestBalanceoRL`            | Rotación RL restaura el balance                          |
| `TestDatosOrdenados`        | Árbol permanece balanceado con inserciones ordenadas     |
| `TestConsultaRango`         | Devuelve exactamente los nodos dentro del rango `[a, b]` |
| `TestConsultaRangoVacio`    | Rango sin resultados devuelve lista vacía                |

### Benchmarks

**Ejecución:**

```shell
go test ./arboles -bench=Benchmark -benchmem
```

**Resultados:**

> Pendiente.

### Conclusiones

> Pendiente.
