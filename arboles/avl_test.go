package arboles

import (
	"testing"
)

// Verifica que el cálculo de altura maneje correctamente los nodos nulos e inicializados.
func TestAltura(t *testing.T) {
	if Altura(nil) != 0 {
		t.Errorf("La altura de un nodo nil debería ser 0")
	}
	n := &NodoAVL{Alt: 5}
	if Altura(n) != 5 {
		t.Errorf("Se esperaba altura 5, se obtuvo %d", Altura(n))
	}
}

// Verifica que el cálculo del balance retorne 0 cuando se evalúa un nodo inexistente.
func TestObtenerBalance(t *testing.T) {
	if obtenerBalance(nil) != 0 {
		t.Errorf("El balance de un nodo nil debería ser 0")
	}
}

// Comprueba la correcta ejecución de las cuatro rotaciones de autobalanceo del árbol AVL.
func TestRotacionesAVL(t *testing.T) {
	arbolLL := &ArbolAVL{}
	arbolLL.Raiz = Insertar(arbolLL.Raiz, 30.0, Pelicula{MovieId: 1}, arbolLL)
	arbolLL.Raiz = Insertar(arbolLL.Raiz, 20.0, Pelicula{MovieId: 2}, arbolLL)
	arbolLL.Raiz = Insertar(arbolLL.Raiz, 10.0, Pelicula{MovieId: 3}, arbolLL)
	if arbolLL.Raiz.Clave != 20.0 {
		t.Errorf("Rotación LL fallida. Raíz esperada: 20.0, obtenida: %.1f", arbolLL.Raiz.Clave)
	}

	arbolRR := &ArbolAVL{}
	arbolRR.Raiz = Insertar(arbolRR.Raiz, 10.0, Pelicula{MovieId: 1}, arbolRR)
	arbolRR.Raiz = Insertar(arbolRR.Raiz, 20.0, Pelicula{MovieId: 2}, arbolRR)
	arbolRR.Raiz = Insertar(arbolRR.Raiz, 30.0, Pelicula{MovieId: 3}, arbolRR)
	if arbolRR.Raiz.Clave != 20.0 {
		t.Errorf("Rotación RR fallida. Raíz esperada: 20.0, obtenida: %.1f", arbolRR.Raiz.Clave)
	}

	arbolLR := &ArbolAVL{}
	arbolLR.Raiz = Insertar(arbolLR.Raiz, 30.0, Pelicula{MovieId: 1}, arbolLR)
	arbolLR.Raiz = Insertar(arbolLR.Raiz, 10.0, Pelicula{MovieId: 2}, arbolLR)
	arbolLR.Raiz = Insertar(arbolLR.Raiz, 20.0, Pelicula{MovieId: 3}, arbolLR)
	if arbolLR.Raiz.Clave != 20.0 {
		t.Errorf("Rotación LR fallida. Raíz esperada: 20.0, obtenida: %.1f", arbolLR.Raiz.Clave)
	}

	arbolRL := &ArbolAVL{}
	arbolRL.Raiz = Insertar(arbolRL.Raiz, 10.0, Pelicula{MovieId: 1}, arbolRL)
	arbolRL.Raiz = Insertar(arbolRL.Raiz, 30.0, Pelicula{MovieId: 2}, arbolRL)
	arbolRL.Raiz = Insertar(arbolRL.Raiz, 20.0, Pelicula{MovieId: 3}, arbolRL)
	if arbolRL.Raiz.Clave != 20.0 {
		t.Errorf("Rotación RL fallida. Raíz esperada: 20.0, obtenida: %.1f", arbolRL.Raiz.Clave)
	}
}

// Valida que el árbol agrupe múltiples películas con el mismo rating numérico en un solo nodo.
func TestInsertarDuplicados(t *testing.T) {
	arbol := &ArbolAVL{}
	arbol.Raiz = Insertar(arbol.Raiz, 4.5, Pelicula{MovieId: 1}, arbol)
	arbol.Raiz = Insertar(arbol.Raiz, 4.5, Pelicula{MovieId: 2}, arbol)

	if len(arbol.Raiz.Datos) != 2 {
		t.Errorf("Se esperaban 2 elementos agrupados, se obtuvieron %d", len(arbol.Raiz.Datos))
	}
}

// Evalúa que la consulta por rango retorne únicamente los elementos en el intervalo exacto y maneje casos vacíos.
func TestConsultaRango(t *testing.T) {
	arbol := &ArbolAVL{}
	claves := []float64{3.0, 4.0, 5.0, 2.0, 4.5, 3.5}
	for i, c := range claves {
		arbol.Raiz = Insertar(arbol.Raiz, c, Pelicula{MovieId: i}, arbol)
	}

	resultados := ConsultaRango(arbol.Raiz, 3.5, 4.5)
	if len(resultados) != 3 {
		t.Errorf("Se esperaban 3 resultados, se obtuvieron %d", len(resultados))
	}

	resultadosFueraRango := ConsultaRango(arbol.Raiz, 8.0, 9.0)
	if len(resultadosFueraRango) != 0 {
		t.Errorf("Se esperaban 0 resultados, se obtuvieron %d", len(resultadosFueraRango))
	}

	resultadosArbolVacio := ConsultaRango(nil, 1.0, 5.0)
	if len(resultadosArbolVacio) != 0 {
		t.Errorf("Se esperaban 0 resultados para un árbol nulo, se obtuvieron %d", len(resultadosArbolVacio))
	}
}

// Mide el rendimiento de inserción de datos crecientes, forzando rebalanceos constantes O(log n).
func BenchmarkInsertar(b *testing.B) {
	arbol := &ArbolAVL{}
	for i := 0; i < b.N; i++ {
		arbol.Raiz = Insertar(arbol.Raiz, float64(i), Pelicula{MovieId: i}, arbol)
	}
}

// Evalúa la velocidad de ejecución de las consultas de rango O(log n + k) sobre un árbol precargado.
func BenchmarkConsultaRango(b *testing.B) {
	arbol := &ArbolAVL{}
	for i := 0; i < 10000; i++ {
		arbol.Raiz = Insertar(arbol.Raiz, float64(i%10), Pelicula{MovieId: i}, arbol)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ConsultaRango(arbol.Raiz, 3.0, 6.0)
	}
}
