package main

import (
	"fmt"
	"math/rand/v2"
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
	
	// ajuste do número de iterações para compensar
	iteracoesPorCPU := (simToRun / numCPUs) / 12
	
	var sucessosGlobais int64 = int64(sucessos)

	for c := 0; c < numCPUs; c++ {
		wg.Add(1)
		go func(n int, id int) {
			defer wg.Done()
			
			pcg := rand.NewPCG(uint64(time.Now().UnixNano()), uint64(id+1))
			r := rand.New(pcg)
			
			localSucessos := 0
			for i := 0; i < n; i++ {
				// Geramos 64 bits de uma vez
				val := r.Uint64()
				
				// Extraímos 12 experimentos de 5 bits cada de um único número
				// reduz a carga sobre o gerador aleatório
				if bits.OnesCount64(val & mask) == alvoCaras { localSucessos++ }
				if bits.OnesCount64((val >> 5) & mask) == alvoCaras { localSucessos++ }
				if bits.OnesCount64((val >> 10) & mask) == alvoCaras { localSucessos++ }
				if bits.OnesCount64((val >> 15) & mask) == alvoCaras { localSucessos++ }
				if bits.OnesCount64((val >> 20) & mask) == alvoCaras { localSucessos++ }
				if bits.OnesCount64((val >> 25) & mask) == alvoCaras { localSucessos++ }
				if bits.OnesCount64((val >> 30) & mask) == alvoCaras { localSucessos++ }
				if bits.OnesCount64((val >> 35) & mask) == alvoCaras { localSucessos++ }
				if bits.OnesCount64((val >> 40) & mask) == alvoCaras { localSucessos++ }
				if bits.OnesCount64((val >> 45) & mask) == alvoCaras { localSucessos++ }
				if bits.OnesCount64((val >> 50) & mask) == alvoCaras { localSucessos++ }
				if bits.OnesCount64((val >> 55) & mask) == alvoCaras { localSucessos++ }
			}
			atomic.AddInt64(&sucessosGlobais, int64(localSucessos))
		}(iteracoesPorCPU, c)
	}
	wg.Wait()
	return int(sucessosGlobais), int(totalSimulacoes+simToRun)
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
