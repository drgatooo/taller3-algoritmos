// Package colas implementa una cola FIFO con lista enlazada simple
// y un rate limiter basado en ventana deslizante de timestamps.
package colas

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// -------------------------------------------------------
// 1. NODO — elemento interno de la lista enlazada
// -------------------------------------------------------

type nodo struct {
	ts   int64 // timestamp Unix (segundos)
	next *nodo
}

// -------------------------------------------------------
// 2. COLA FIFO
// -------------------------------------------------------

// Cola es una cola FIFO de timestamps implementada con lista enlazada.
// head apunta al elemento más antiguo (frente); tail al más nuevo.
// Complejidad: Enqueue O(1), Dequeue O(1), Front O(1), Len O(1).
type Cola struct {
	head *nodo
	tail *nodo
	len  int
}

// Enqueue agrega el timestamp ts al final de la cola. O(1).
func (c *Cola) Enqueue(ts int64) {
	n := &nodo{ts: ts}
	if c.tail != nil {
		c.tail.next = n
	}
	c.tail = n
	if c.head == nil {
		c.head = n
	}
	c.len++
}

// Dequeue extrae el timestamp del frente de la cola. O(1).
// Devuelve (0, false) si la cola está vacía.
func (c *Cola) Dequeue() (int64, bool) {
	if c.head == nil {
		return 0, false
	}
	ts := c.head.ts
	c.head = c.head.next
	if c.head == nil {
		c.tail = nil
	}
	c.len--
	return ts, true
}

// Front devuelve el timestamp del frente sin eliminarlo. O(1).
// Devuelve (0, false) si la cola está vacía.
func (c *Cola) Front() (int64, bool) {
	if c.head == nil {
		return 0, false
	}
	return c.head.ts, true
}

// Len devuelve la cantidad de elementos en la cola. O(1).
func (c *Cola) Len() int {
	return c.len
}

// -------------------------------------------------------
// 3. RATE LIMITER — ventana deslizante
// -------------------------------------------------------

// PermitirPeticion decide si la petición de ip en el instante ts
// se acepta bajo la política: máximo M peticiones en los últimos T segundos.
//
// Complejidad amortizada O(1) por llamada:
// cada timestamp se inserta una vez (Enqueue) y se elimina una vez (Dequeue),
// por lo que el costo total sobre n peticiones es O(n).
func PermitirPeticion(colas map[string]*Cola, ip string, ts int64, M int, T int64) bool {
	c, existe := colas[ip]
	if !existe {
		c = &Cola{}
		colas[ip] = c
	}

	// Descartar del frente los timestamps fuera de la ventana [ts-T, ts]
	limite := ts - T
	for {
		frente, ok := c.Front()
		if !ok || frente > limite {
			break
		}
		c.Dequeue()
	}

	// Verificar cupo
	if c.Len() >= M {
		return false // RECHAZADA
	}

	// Registrar esta petición y aceptar
	c.Enqueue(ts)
	return true
}

// -------------------------------------------------------
// 4. PARSEO DE LÍNEA — formato Apache Common Log
// -------------------------------------------------------

// Registro representa una línea del log ya parseada.
type Registro struct {
	IP string
	TS int64
}

// ParsearLinea extrae la IP y el timestamp de una línea en formato Apache:
// 127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] "GET /..." 200 2326
// Devuelve error si la línea no tiene el formato esperado.
func ParsearLinea(linea string) (ip string, ts int64, err error) {
	partes := strings.Fields(linea)
	if len(partes) < 4 {
		return "", 0, fmt.Errorf("línea inválida: %q", linea)
	}

	ip = partes[0]

	// El timestamp está en partes[3], e.g. [10/Oct/2000:13:55:36
	fechaRaw := strings.TrimPrefix(partes[3], "[")
	t, errT := time.Parse("02/Jan/2006:15:04:05", fechaRaw)
	if errT != nil {
		return "", 0, fmt.Errorf("timestamp inválido %q: %w", fechaRaw, errT)
	}

	return ip, t.Unix(), nil
}

// -------------------------------------------------------
// 5. PROCESAMIENTO DEL LOG COMPLETO
// -------------------------------------------------------

// Resultado almacena el resumen global tras procesar el log.
type Resultado struct {
	TotalPeticiones int
	TotalRechazos   int
	RechazosPorIP   map[string]int
}

// ProcesarLog lee el archivo en ruta, aplica el rate limiter con parámetros M y T,
// imprime las primeras muestra decisiones y devuelve el resumen global.
// Complejidad total: O(n) siendo n el número de líneas del log.
func ProcesarLog(ruta string, M int, T int64, muestra int) (Resultado, error) {
	f, err := os.Open(ruta)
	if err != nil {
		return Resultado{}, fmt.Errorf("no se pudo abrir el log: %w", err)
	}
	defer f.Close()

	colasIP := make(map[string]*Cola)
	rechazosPorIP := make(map[string]int)
	totalPeticiones, totalRechazos, mostradas := 0, 0, 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		linea := scanner.Text()
		ip, ts, err := ParsearLinea(linea)
		if err != nil {
			continue // saltar líneas mal formadas
		}

		totalPeticiones++
		aceptada := PermitirPeticion(colasIP, ip, ts, M, T)

		if !aceptada {
			totalRechazos++
			rechazosPorIP[ip]++
		}

		if mostradas < muestra {
			estado := "ACEPTADA "
			if !aceptada {
				estado = "RECHAZADA"
			}
			fmt.Printf("IP: %-15s  ts: %d  → %s\n", ip, ts, estado)
			mostradas++
		}
	}

	if err := scanner.Err(); err != nil {
		return Resultado{}, fmt.Errorf("error leyendo log: %w", err)
	}

	return Resultado{
		TotalPeticiones: totalPeticiones,
		TotalRechazos:   totalRechazos,
		RechazosPorIP:   rechazosPorIP,
	}, nil
}
