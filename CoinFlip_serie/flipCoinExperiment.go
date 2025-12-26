package main

import (
	"fmt"
	"math/rand"
    "math/bits"
    "math"
	"time"
    "runtime"
    "sync"
    "sync/atomic"
)

/*
change numSimulations for simtoRun to say how much simulations need to run
*/
func rodada( simtoRun int, totalSimulacoes int, sucessos int, r *rand.Rand) (int, int){
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

func rodadaParalela(simToRun int, totalSimulacoes int, sucessos int) (int, int) {
	var wg sync.WaitGroup
	numCPUs := runtime.NumCPU()
	chunks := simToRun / numCPUs
	var sucessosRodada int64 = int64(sucessos)

	for c := 0; c < numCPUs; c++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
            // esqueci e tentei compartilhar o gerador mas e cai no caso do mutex. 
            // Nao compartilhem geradores, amiguinhos. Só se for um com mutex :)
            // melhor que cada uma das goroutines tenha seu proprio gerador.
            r := rand.New(rand.NewSource(time.Now().UnixNano()))
			localSucessos := 0
			
			for i := 0; i < n; i++ {
				res := r.Uint32() & mask
				
				if bits.OnesCount32(res) == alvoCaras {
					localSucessos++
				}
			}
			atomic.AddInt64(&sucessosRodada, int64(localSucessos))
		}(chunks)
	}
	wg.Wait()
	return int(sucessosRodada), (totalSimulacoes + simToRun)
}

const lancamentosPorVez = 5
const alvoCaras = 3
const mask = (1 << lancamentosPorVez) -1
func main() {
    variacaoAceita := 0.00014
    stepSize := 1000000
    nivelConfianca := 1.96 // = 95% de confiança.
    sucessos := 0
	totalSimulacoes := 10000000

    probabilidadeEstimada:= 0.0

    sucessos, totalSimulacoes = rodadaParalela(totalSimulacoes,0,sucessos)
    // start simulations 
    for {
            probabilidadeEstimada = float64(sucessos) / float64(totalSimulacoes)
            erroPadrao := math.Sqrt( (probabilidadeEstimada * ( 1- probabilidadeEstimada)) / float64(totalSimulacoes))
            margemErroAtual := nivelConfianca * erroPadrao

            // start by checking the stop rules
            if  margemErroAtual > variacaoAceita {
            }else {
                break
            }

            sucessos, totalSimulacoes = rodadaParalela(stepSize, totalSimulacoes, sucessos)
    }

	fmt.Printf("Simulou %d rodadas de %d lançamentos: sucessos %d...\n", totalSimulacoes, lancamentosPorVez, sucessos)
	fmt.Printf("Probabilidade Estimada: %.4f (%.2f%%)\n", probabilidadeEstimada, probabilidadeEstimada*100)
}
