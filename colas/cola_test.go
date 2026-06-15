package colas

import (
	"testing"
)

// -------------------------------------------------------
// Tests de la Cola
// -------------------------------------------------------

// TestColaFIFO verifica que los elementos salen en orden FIFO.
func TestColaFIFO(t *testing.T) {
	c := &Cola{}
	c.Enqueue(100)
	c.Enqueue(200)
	c.Enqueue(300)

	for _, esperado := range []int64{100, 200, 300} {
		v, ok := c.Dequeue()
		if !ok || v != esperado {
			t.Errorf("esperaba %d, obtuve %d (ok=%v)", esperado, v, ok)
		}
	}
}

// TestColaVacia verifica que Dequeue en cola vacía devuelve false.
func TestColaVacia(t *testing.T) {
	c := &Cola{}
	_, ok := c.Dequeue()
	if ok {
		t.Error("Dequeue en cola vacía debería devolver false")
	}
	_, ok = c.Front()
	if ok {
		t.Error("Front en cola vacía debería devolver false")
	}
}

// TestColaUnElemento verifica el caso límite de un solo elemento.
func TestColaUnElemento(t *testing.T) {
	c := &Cola{}
	c.Enqueue(42)
	if c.Len() != 1 {
		t.Errorf("Len esperado 1, obtenido %d", c.Len())
	}
	v, ok := c.Dequeue()
	if !ok || v != 42 {
		t.Errorf("esperaba 42, obtuve %d", v)
	}
	if c.Len() != 0 {
		t.Errorf("Len esperado 0 tras Dequeue, obtenido %d", c.Len())
	}
}

// TestFront verifica que Front no modifica la cola.
func TestFront(t *testing.T) {
	c := &Cola{}
	c.Enqueue(10)
	c.Enqueue(20)

	f, ok := c.Front()
	if !ok || f != 10 {
		t.Errorf("Front esperaba 10, obtuvo %d", f)
	}
	if c.Len() != 2 {
		t.Error("Front no debería modificar el tamaño de la cola")
	}
}

// -------------------------------------------------------
// Tests del Rate Limiter
// -------------------------------------------------------

// TestPermitirPeticionLimite verifica que la petición M+1 se rechaza.
func TestPermitirPeticionLimite(t *testing.T) {
	colasIP := make(map[string]*Cola)
	ip := "1.2.3.4"
	M, T := 3, int64(60)
	ts := int64(1000)

	for i := 0; i < M; i++ {
		if !PermitirPeticion(colasIP, ip, ts, M, T) {
			t.Errorf("petición %d debería aceptarse", i+1)
		}
	}
	if PermitirPeticion(colasIP, ip, ts, M, T) {
		t.Error("la petición M+1 debería rechazarse")
	}
}

// TestVentanaDeslizante verifica que los timestamps viejos expiran y se vuelve a aceptar.
func TestVentanaDeslizante(t *testing.T) {
	colasIP := make(map[string]*Cola)
	ip := "5.6.7.8"
	M, T := 2, int64(10)

	// Llena la ventana en t=0
	PermitirPeticion(colasIP, ip, 0, M, T)
	PermitirPeticion(colasIP, ip, 0, M, T)

	// En t=5 sigue llena
	if PermitirPeticion(colasIP, ip, 5, M, T) {
		t.Error("debería rechazarse en t=5 (ventana llena)")
	}

	// En t=11 los timestamps de t=0 expiraron → se acepta
	if !PermitirPeticion(colasIP, ip, 11, M, T) {
		t.Error("debería aceptarse en t=11 (timestamps expirados)")
	}
}

// TestIPsIndependientes verifica que cada IP tiene su propia cola.
func TestIPsIndependientes(t *testing.T) {
	colasIP := make(map[string]*Cola)
	M, T := 1, int64(60)
	ts := int64(1000)

	if !PermitirPeticion(colasIP, "a.a.a.a", ts, M, T) {
		t.Error("primera petición de IP A debería aceptarse")
	}
	if !PermitirPeticion(colasIP, "b.b.b.b", ts, M, T) {
		t.Error("primera petición de IP B debería aceptarse")
	}
	if PermitirPeticion(colasIP, "a.a.a.a", ts, M, T) {
		t.Error("segunda petición de IP A debería rechazarse")
	}
}

// -------------------------------------------------------
// Tests del Parser
// -------------------------------------------------------

// TestParsearLineaOK verifica una línea Apache válida.
func TestParsearLineaOK(t *testing.T) {
	linea := `83.149.9.216 - - [17/May/2015:10:05:03 +0000] "GET /presentations HTTP/1.1" 200 5678`
	ip, ts, err := ParsearLinea(linea)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if ip != "83.149.9.216" {
		t.Errorf("IP esperada 83.149.9.216, obtenida %s", ip)
	}
	if ts <= 0 {
		t.Errorf("timestamp debe ser positivo, obtenido %d", ts)
	}
}

// TestParsearLineaInvalida verifica que una línea malformada devuelve error.
func TestParsearLineaInvalida(t *testing.T) {
	_, _, err := ParsearLinea("esto no es un log")
	if err == nil {
		t.Error("se esperaba error con línea inválida")
	}
}

// -------------------------------------------------------
// Benchmarks
// -------------------------------------------------------

// BenchmarkEnqueueDequeue mide el costo de una operación Enqueue+Dequeue. Esperado: O(1).
func BenchmarkEnqueueDequeue(b *testing.B) {
	c := &Cola{}
	for i := 0; i < b.N; i++ {
		c.Enqueue(int64(i))
		c.Dequeue()
	}
}

// BenchmarkPermitirPeticion mide el costo del rate limiter con ventana activa. Esperado: O(1) amortizado.
func BenchmarkPermitirPeticion(b *testing.B) {
	colasIP := make(map[string]*Cola)
	M, T := 100, int64(60)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PermitirPeticion(colasIP, "1.1.1.1", int64(i), M, T)
	}
}
