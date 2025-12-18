package main

import (
	"fmt"
	"math/rand"
	"time"
    "os"
    "strconv"
    "log"
)

// 0= coroa, 1=cara 
func lancaMoeda() int{
    return rand.Intn(2)
}
func main() {

	const lancamentosPorVez = 5
	const alvoCaras = 3

	totalSimulacoes, err := strconv.Atoi(os.Args[1])
    if err != nil{
        log.Fatal(err)
    }

	// Seed de aleatoriedade 
	rand.Seed(time.Now().UnixNano())

	sucessos := 0

    // para cada simulacao
	for i := 0; i < totalSimulacoes; i++ {
		carasNoExperimento := 0

        // lancar a moeda 5 vezes.	
		for j := 0; j < lancamentosPorVez; j++ {
			if lancaMoeda() == 1 {
				carasNoExperimento++
			}
		}

		// Se o experimento resultou em exatamente 3 caras, contamos como sucesso
		if carasNoExperimento == alvoCaras {
			sucessos++
		}
	}

	// Calculando a probabilidade estimada
	probabilidadeEstimada := float64(sucessos) / float64(totalSimulacoes)

	fmt.Printf("Simulou %d rodadas de %d lançamentos...\n", totalSimulacoes, lancamentosPorVez)
	fmt.Printf("Casos com exatamente %d caras: %d\n", alvoCaras, sucessos)
	fmt.Printf("Probabilidade Estimada: %.4f (%.2f%%)\n", probabilidadeEstimada, probabilidadeEstimada*100)
	fmt.Println("Probabilidade Teórica: 0.3125 (31.25%)")
}
