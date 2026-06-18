# Taller 3 — Algoritmos y Estructuras de Datos

Integrantes

- Anchi Cristobal Ernesto Alonso
- Cesar Gomez Chavez
- Contreras Salcedo Maximo Simon

### Videos Explicativos

- Ejercicio 2: https://youtu.be/r3se1OOeAQ4
- Ejercicio 3: https://youtu.be/ACihE3TRW40
- Ejercicio 4: pendiente

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

Implementar un árbol **AVL** (árbol binario de búsqueda autobalanceado) que garantice altura $O(\log n)$ e indexe películas de MovieLens por clave numérica, permitiendo consultas por rango `[a, b]` eficientes.

### Dataset

**MovieLens** — GroupLens Research

https://grouplens.org/datasets/movielens/

Se utilizan `movies.csv` y `ratings.csv`. La clave de indexación es el **rating promedio** por película (calculado agregando `ratings.csv`) o el **año** extraído del título en `movies.csv`.

Ubicación esperada: `data/movies.csv` y `data/ratings.csv`

### Estructura del Proyecto

```text
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

El Índice AVL ha sido diseñado para maximizar el rendimiento al tratar con colisiones (películas con el mismo rating o del mismo año) y para optimizar las consultas masivas mediante poda geométrica de subárboles.

#### Estructura del Nodo y Manejo de Colisiones
En un dataset real como MovieLens, muchas películas comparten exactamente la misma calificación promedio (ej. 4.0). En lugar de insertar estos duplicados como nodos hijos (lo que aumentaría la profundidad del árbol innecesariamente), se implementó un `NodoAVL` que agrupa las colisiones en un *slice* continuo en memoria `Datos []Pelicula`.

```go
type Pelicula struct {
	MovieID int
	Title   string
	Rating  float64
}

type NodoAVL struct {
	Clave float64
	Datos []Pelicula // Agrupación de registros con la misma clave O(1)
	Alt   int
	Izq   *NodoAVL
	Der   *NodoAVL
}
```

#### Inserción y Autobalanceo O(log n)

La función `Insertar` agrega el nuevo registro siguiendo las reglas de un Árbol Binario de Búsqueda. En el retorno de la recursión, actualiza las alturas y evalúa el **factor de balance**. Si el subárbol se desequilibra ($|FB| > 1$), aplica las rotaciones necesarias (LL, RR, LR, RL) reasignando los punteros en tiempo constante.

```go
func Insertar(raiz *NodoAVL, clave float64, dato Pelicula, arbol *ArbolAVL) *NodoAVL
```

#### Consultas por Rango con Poda O(log n + k)

Para garantizar una alta eficiencia, la función `ConsultaRango` implementa un algoritmo de **descarte inteligente**. En lugar de recorrer todo el árbol, compara la clave del nodo actual con los límites `[a, b]` para decidir si es matemáticamente posible encontrar resultados en sus subárboles, podando ramas enteras de forma anticipada.

```go
func ConsultaRango(raiz *NodoAVL, a, b float64) []Pelicula
```

### Complejidad Temporal

| Operación       | Complejidad     |
| --------------- | --------------- |
| `Insertar`      | $O(\log n)$     |
| `ConsultaRango` | $O(\log n + k)$ |
| Rotaciones (×4) | $O(1)$ c/u      |

**Justificación:** El balanceo automático (factor $|FB| \le 1$ en cada nodo) garantiza que la altura del árbol sea estrictamente $O(\log n)$ incluso si los datos se insertan en orden creciente o decreciente. La consulta por rango poda subárboles completos fuera del intervalo `[a, b]`, logrando la meta $O(\log n + k)$ donde $k$ es el número total de películas retornadas.

### Ejecución

```shell
go run ./arboles/cmd/main.go <ruta_ratings> <ruta_movies>
```

**Parámetros:**

| Argumento | Descripción |
| --- | --- |
| `<ruta_ratings>` | Ruta a `ratings.csv` |
| `<ruta_movies>` | Ruta a `movies.csv` |

### Pruebas Unitarias

**Ejecución:**

```shell
go test ./arboles -v -cover
```

**Resultados:**

```text
=== RUN   TestAltura
--- PASS: TestAltura (0.00s)
=== RUN   TestObtenerBalance
--- PASS: TestObtenerBalance (0.00s)
=== RUN   TestRotacionesAVL
--- PASS: TestRotacionesAVL (0.00s)
=== RUN   TestInsertarDuplicados
--- PASS: TestInsertarDuplicados (0.00s)
=== RUN   TestConsultaRango
--- PASS: TestConsultaRango (0.00s)
PASS
coverage: 100.0% of statements
ok      taller3/arboles 0.002s  coverage: 100.0% of statements
```

Las pruebas validan que todas las rotaciones funcionen independientemente, aseguran que las agrupaciones en *slices* operen sin sobrescribir datos y comprueban que la lógica de poda del intervalo extraiga los datos exactos.

### Benchmarks

**Ejecución:**

```shell
go test ./arboles -bench=. -benchmem
```

**Resultados:**

```text
goos: linux
goarch: amd64
pkg: taller3/arboles
cpu: AMD Ryzen 7 5700G with Radeon Graphics         
BenchmarkInsertar-16             6999496               194.0 ns/op            96 B/op          2 allocs/op
BenchmarkConsultaRango-16          19891               65664 ns/op        548866 B/op          8 allocs/op
PASS
ok      taller3/arboles 3.473s
```

### Conclusiones

- A diferencia de un *Binary Search Tree* (BST) estándar que degradaría su eficiencia a $O(n)$ si recibe datos ya ordenados o secuenciales, el árbol implementado activa sus mecanismos de rebalanceo constantes para mantener la forma plana e indexada.
- La decisión de crear arreglos anidados dentro del nodo fue clave. En bases de datos reales como MovieLens, muchos registros comparten exactamente la misma calificación (ej. 4.0 o 5.0). Agruparlos redujo enormemente la profundidad estructural total del árbol.
- La poda implementada en `ConsultaRango` evita iterar todo el dataset (como pasaría en una consulta tradicional de array), permitiendo descartar miles de nodos en memoria en microsegundos y cumpliendo exitosamente el requisito de $O(\log n + k)$.
