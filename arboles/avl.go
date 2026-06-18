package arboles

// Pelicula representa el registro de nuestro dataset
type Pelicula struct {
	MovieId int
	Title   string
	Rating  float64
}

// NodoAVL es la estructura base sugerida para el árbol
type NodoAVL struct {
	Clave float64
	Datos []Pelicula
	Alt   int
	Izq   *NodoAVL
	Der   *NodoAVL
}

// ArbolAVL es un contenedor para mantener las métricas fácilmente
type ArbolAVL struct {
	Raiz       *NodoAVL
	Nodos      int
	Rotaciones int
}

// Altura devuelve la altura del nodo de forma segura
func Altura(n *NodoAVL) int {
	if n == nil {
		return 0
	}
	return n.Alt
}

func obtenerBalance(n *NodoAVL) int {
	if n == nil {
		return 0
	}
	return Altura(n.Izq) - Altura(n.Der)
}

// Rotación LL
func rotarDer(y *NodoAVL, rotaciones *int) *NodoAVL {
	*rotaciones++
	x := y.Izq
	T2 := x.Der

	// Rotación
	x.Der = y
	y.Izq = T2

	// Actualizar alturas
	y.Alt = max(Altura(y.Izq), Altura(y.Der)) + 1
	x.Alt = max(Altura(x.Izq), Altura(x.Der)) + 1

	return x
}

// Rotación RR
func rotarIzq(x *NodoAVL, rotaciones *int) *NodoAVL {
	*rotaciones++
	y := x.Der
	T2 := y.Izq

	// Rotación
	y.Izq = x
	x.Der = T2

	// Actualizar alturas
	x.Alt = max(Altura(x.Izq), Altura(x.Der)) + 1
	y.Alt = max(Altura(y.Izq), Altura(y.Der)) + 1

	return y
}

// Insertar añade un nuevo registro al árbol y realiza el autobalanceo
func Insertar(raiz *NodoAVL, clave float64, dato Pelicula, arbol *ArbolAVL) *NodoAVL {
	// Inserción normal de BST
	if raiz == nil {
		arbol.Nodos++
		return &NodoAVL{
			Clave: clave,
			Datos: []Pelicula{dato},
			Alt:   1,
		}
	}

	if clave < raiz.Clave {
		raiz.Izq = Insertar(raiz.Izq, clave, dato, arbol)
	} else if clave > raiz.Clave {
		raiz.Der = Insertar(raiz.Der, clave, dato, arbol)
	} else {
		// Si hay claves iguales, se agrupan en el slice
		raiz.Datos = append(raiz.Datos, dato)
		return raiz
	}

	// Actualizar altura del nodo actual
	raiz.Alt = 1 + max(Altura(raiz.Izq), Altura(raiz.Der))

	// Obtener factor de balance
	balance := obtenerBalance(raiz)

	// Rotación LL
	if balance > 1 && clave < raiz.Izq.Clave {
		return rotarDer(raiz, &arbol.Rotaciones)
	}

	// Rotación RR
	if balance < -1 && clave > raiz.Der.Clave {
		return rotarIzq(raiz, &arbol.Rotaciones)
	}

	// Rotación LR
	if balance > 1 && clave > raiz.Izq.Clave {
		raiz.Izq = rotarIzq(raiz.Izq, &arbol.Rotaciones)
		return rotarDer(raiz, &arbol.Rotaciones)
	}

	// Rotación RL
	if balance < -1 && clave < raiz.Der.Clave {
		raiz.Der = rotarDer(raiz.Der, &arbol.Rotaciones)
		return rotarIzq(raiz, &arbol.Rotaciones)
	}

	return raiz
}

// ConsultaRango devuelve los registros dentro del intervalo [a, b]
func ConsultaRango(raiz *NodoAVL, a, b float64) []Pelicula {
	var resultado []Pelicula
	if raiz == nil {
		return resultado
	}

	// Si la clave actual es mayor que a, puede haber respuestas en el subárbol izquierdo
	if raiz.Clave > a {
		resultado = append(resultado, ConsultaRango(raiz.Izq, a, b)...)
	}

	// Si está en el rango, agregar todos los registros almacenados en el nodo
	if raiz.Clave >= a && raiz.Clave <= b {
		resultado = append(resultado, raiz.Datos...)
	}

	// Si la clave actual es menor que b, puede haber respuestas en el subárbol derecho
	if raiz.Clave < b {
		resultado = append(resultado, ConsultaRango(raiz.Der, a, b)...)
	}

	return resultado
}
