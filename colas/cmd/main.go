// Ejecutable principal del ejercicio 4.2 — Rate Limiter
// Uso: go run ./colas/cmd/ <ruta_log> <M> <T_segundos>
package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"taller3/colas"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Uso: go run ./colas/cmd/ <ruta_log> <M_peticiones> <T_segundos>")
		fmt.Println("Ejemplo: go run ./colas/cmd/ access.log 10 60")
		os.Exit(1)
	}

	ruta := os.Args[1]
	M, errM := strconv.Atoi(os.Args[2])
	T, errT := strconv.ParseInt(os.Args[3], 10, 64)
	if errM != nil || errT != nil || M <= 0 || T <= 0 {
		fmt.Fprintln(os.Stderr, "Error: M y T deben ser enteros positivos")
		os.Exit(1)
	}

	fmt.Printf("=== Rate Limiter | M=%d peticiones / T=%d segundos ===\n\n", M, T)
	fmt.Println("--- Muestra de las primeras 20 decisiones ---")

	res, err := colas.ProcesarLog(ruta, M, T, 20)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	// Resumen global
	fmt.Printf("\n--- Resumen global ---\n")
	fmt.Printf("Total peticiones : %d\n", res.TotalPeticiones)
	fmt.Printf("Total rechazos   : %d (%.1f%%)\n",
		res.TotalRechazos,
		100*float64(res.TotalRechazos)/float64(res.TotalPeticiones))

	// Top 5 IPs con más rechazos
	type ipCount struct {
		ip    string
		count int
	}
	ranking := make([]ipCount, 0, len(res.RechazosPorIP))
	for ip, n := range res.RechazosPorIP {
		ranking = append(ranking, ipCount{ip, n})
	}
	sort.Slice(ranking, func(i, j int) bool {
		return ranking[i].count > ranking[j].count
	})

	fmt.Println("\n--- Top 5 IPs con más rechazos ---")
	for i, r := range ranking {
		if i == 5 {
			break
		}
		fmt.Printf("%d. %-15s → %d rechazos\n", i+1, r.ip, r.count)
	}
}
