package main

import (
	"fmt"
	"math/rand"
    "math/bits"
    "math"
	"time"
)

/*
change numSimulations for simtoRun to say how much simulations need to run
*/
func rodada( simtoRun int, totalSimulacoes int, sucessos int, r * rand.Rand) (int, int){
    fmt.Printf("Rodando %d simulacoes a mais...\n", simtoRun)
	for i := 0; i < simtoRun; i++ {
		carasNoExperimento := 0
        // lancar a moeda 5 vezes.	
        res := r.Uint32() & mask
        carasNoExperimento = bits.OnesCount32(res)
		if carasNoExperimento == alvoCaras {
			sucessos++
		}
	}
	return sucessos, (totalSimulacoes + simtoRun)
}

const lancamentosPorVez = 5
const alvoCaras = 3
const mask = (1 << lancamentosPorVez) -1
func main() {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    // delta aceitavel na estimativa 0.0001

    variacaoAceita := 0.00014
    stepSize := 1000000
    nivelConfianca := 1.96 // = 95% de confiança.
    sucessos := 0
	totalSimulacoes := 10000000

    probabilidadeEstimada:= 0.0

    sucessos, totalSimulacoes = rodada(totalSimulacoes,0,sucessos, r)
    // start simulations 
    for {
            probabilidadeEstimada = float64(sucessos) / float64(totalSimulacoes)
            erroPadrao := math.Sqrt( (probabilidadeEstimada * ( 1- probabilidadeEstimada)) / float64(totalSimulacoes))
            margemErroAtual := nivelConfianca * erroPadrao

            // start by checking the stop rules
            if  margemErroAtual > variacaoAceita {
                fmt.Printf("[NOK] Margem de Erro %f > %f.\n", margemErroAtual , variacaoAceita)
            }else {
                fmt.Printf("[OK] Margem de Erro %f < %f.\n", margemErroAtual,variacaoAceita)
                break
            }

            fmt.Printf("\n----------\n Tamanho Amostragem/sucessos: %d / %d \n", totalSimulacoes, sucessos)
            fmt.Printf("Estimativa : %f | Margem de erro : %.4f \n", probabilidadeEstimada, margemErroAtual )

            sucessos, totalSimulacoes = rodada(stepSize, totalSimulacoes, sucessos, r)
    }

	fmt.Printf("Simulou %d rodadas de %d lançamentos: sucessos %d...\n", totalSimulacoes, lancamentosPorVez, sucessos)
	fmt.Printf("Probabilidade Estimada: %.4f (%.2f%%)\n", probabilidadeEstimada, probabilidadeEstimada*100)
}
