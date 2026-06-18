package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"taller3/arboles"
)

type RatingAggregation struct {
	Suma     float64
	Cantidad int
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso correcto: go run main.go <ruta_ratings.csv> <ruta_movies.csv>")
		fmt.Println("Ejemplo: go run main.go ./data/ratings_sample.csv ./data/movies.csv")
		os.Exit(1)
	}

	ratingsPath := os.Args[1]
	moviesPath := os.Args[2]

	moviesMap := cargarPeliculas(moviesPath)
	ratingsMap := calcularPromedios(ratingsPath)
	arbol := &arboles.ArbolAVL{}

	for movieID, agg := range ratingsMap {
		avgRating := agg.Suma / float64(agg.Cantidad)

		title := "Película sin título mapeado"
		if t, ok := moviesMap[movieID]; ok {
			title = t
		}

		pelicula := arboles.Pelicula{
			MovieId: movieID,
			Title:   title,
			Rating:  avgRating,
		}

		arbol.Raiz = arboles.Insertar(arbol.Raiz, avgRating, pelicula, arbol)
	}

	fmt.Println("\n================ METRICAS DEL ÁRBOL AVL ================")
	fmt.Printf("Altura del árbol:              %d\n", arboles.Altura(arbol.Raiz))
	fmt.Printf("Número total de nodos:         %d\n", arbol.Nodos)
	fmt.Printf("Rotaciones totales realizadas: %d\n", arbol.Rotaciones)

	balanceRaiz := arboles.Altura(arbol.Raiz.Izq) - arboles.Altura(arbol.Raiz.Der)
	fmt.Printf("Factor de balance en la raíz:  %d\n", balanceRaiz)
	fmt.Println("========================================================")
	fmt.Println()

	var a, b float64
	fmt.Print("Límite inferior del rango (a): ")
	if _, err := fmt.Scan(&a); err != nil {
		log.Fatalf("Error al leer el límite inferior: %v", err)
	}
	fmt.Print("Límite superior del rango (b): ")
	if _, err := fmt.Scan(&b); err != nil {
		log.Fatalf("Error al leer el límite superior: %v", err)
	}

	fmt.Printf("\nConsultando en intervalo [%.2f, %.2f]...\n", a, b)
	resultados := arboles.ConsultaRango(arbol.Raiz, a, b)

	fmt.Printf("\nSe encontraron %d registros dentro del rango:\n", len(resultados))
	for _, p := range resultados {
		fmt.Printf("(ID: %7d) (★ %.2f) %s\n", p.MovieId, p.Rating, p.Title)
	}
}

func cargarPeliculas(path string) map[int]string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error al abrir el archivo de películas: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Omitir la fila de encabezados
	_, _ = reader.Read()

	moviesMap := make(map[int]string)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error leyendo movies.csv: %v", err)
		}
		if len(row) < 2 {
			continue
		}

		id, err1 := strconv.Atoi(row[0])
		if err1 != nil {
			continue
		}
		moviesMap[id] = row[1]
	}
	return moviesMap
}

func calcularPromedios(path string) map[int]*RatingAggregation {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error al abrir el archivo de ratings: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Omitir la fila de encabezados
	_, _ = reader.Read()

	ratingsMap := make(map[int]*RatingAggregation)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error leyendo ratings.csv: %v", err)
		}
		if len(row) < 3 {
			continue
		}

		movieID, err1 := strconv.Atoi(row[1])
		rating, err2 := strconv.ParseFloat(row[2], 64)
		if err1 != nil || err2 != nil {
			continue
		}

		if _, existe := ratingsMap[movieID]; !existe {
			ratingsMap[movieID] = &RatingAggregation{}
		}
		ratingsMap[movieID].Suma += rating
		ratingsMap[movieID].Cantidad++
	}
	return ratingsMap
}
