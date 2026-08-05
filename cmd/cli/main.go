package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rianeiromiron/code-challenge/internal/api"
)

func runCLI(pais string, limit int) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resultado, err := api.ObtenerClima(ctx, pais, limit)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("\n=== %s (%.4f, %.4f) ===\n", resultado.Pais, resultado.Latitud, resultado.Longitud)
	for i, r := range resultado.Resultados {
		fmt.Printf("\n--- #%d ---\n", i+1)
		fmt.Printf("Hora: %s\n", r.FechaHora)
		fmt.Printf("Temp: %.1f°C | Presión: %.1f hPa | Humedad: %.1f%%\n", r.Temperatura, r.Presion, r.Humedad)
		fmt.Printf("Viento: %.1f m/s (%.0f°) | Nubes: %.1f%%\n", r.Viento, r.VientoDir, r.Nubes)
	}
}

func runCLIInteractive() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Ingrese el país a buscar: ")
	pais, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error leyendo input:", err)
		return
	}
	pais = strings.TrimSpace(pais)

	fmt.Print("Ingrese el número de resultados a buscar: ")
	limitStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error leyendo input:", err)
		return
	}
	limitStr = strings.TrimSpace(limitStr)
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		fmt.Println("Error: ingrese un número válido")
		return
	}

	runCLI(pais, limit)
}

func main() {
	var pais string
	var limit int

	flag.StringVar(&pais, "pais", "", "País a buscar")
	flag.IntVar(&limit, "limit", 1, "Número de resultados")
	flag.Parse()

	if pais != "" {
		runCLI(pais, limit)
	} else {
		runCLIInteractive()
	}
}
