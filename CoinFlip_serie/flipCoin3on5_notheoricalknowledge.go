package main

import (
	"fmt"
	"math/rand"
    "math"
	"time"
)

// 0= coroa, 1=cara 
func lancaMoeda() int{
    return rand.Intn(2)
}

func rodada( numSimulations int) float64 {

	sucessos := 0
	for i := 0; i < numSimulations ; i++ {
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
	return float64(sucessos) / float64(numSimulations)
}

const lancamentosPorVez = 5
const alvoCaras = 3

func main() {

	rand.Seed(time.Now().UnixNano())
    // delta aceitavel na estimativa 0.0001
    variacaoAceita := 0.00014
	totalSimulacoes := 10000000
    stepSize := 1000000
    nivelConfianca := 1.96 // = 95% de confiança.
    

    probabilidadeEstimada:= 0.0
    for {
        fmt.Printf("\n----------\n Tamanho Amostragem: %d\n", totalSimulacoes)
        probabilidadeEstimada = rodada(totalSimulacoes)
        erroPadrao := math.Sqrt( (probabilidadeEstimada * ( 1- probabilidadeEstimada)) / float64(totalSimulacoes))
        margemErroAtual := nivelConfianca * erroPadrao
            fmt.Printf("Estimativa : %f | Margem de erro : %.4f \n", probabilidadeEstimada, margemErroAtual )
            if  margemErroAtual > variacaoAceita {
                totalSimulacoes = totalSimulacoes + stepSize
                fmt.Printf("[NOK] Margem de Erro %f > %f.\n", margemErroAtual , variacaoAceita)
            }else {
                fmt.Printf("[OK] Margem de Erro %f < %f.\n", margemErroAtual,variacaoAceita)
                break
            }
    }
	fmt.Printf("Simulou %d rodadas de %d lançamentos...\n", totalSimulacoes, lancamentosPorVez)
	fmt.Printf("Probabilidade Estimada: %.4f (%.2f%%)\n", probabilidadeEstimada, probabilidadeEstimada*100)
}
