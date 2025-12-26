package main

import (
	"fmt"
	"math/rand"
	"time"
    "os"
    "strconv"
    "log"
    "math"
    "sync"
    "sync/atomic"
)

// 0= coroa, 1=cara 
func lancaMoeda( sucessos *atomic.Int64, wg *sync.WaitGroup){
    defer wg.Done()
    if rand.Intn(2) == 1 {
        sucessos.Add(1)
    }
}

func rodada(totalSimulacoes int64, sucessos *atomic.Int64, wg *sync.WaitGroup){
    defer wg.Done()

	for i := int64(0); i < totalSimulacoes; i++ {
        // lancar a moeda 5 vezes.	
		for j := 0; j < lancamentosPorVez; j++ {
            wg.Add(1)
			go lancaMoeda(sucessos, wg)
		}
	}
}

const lancamentosPorVez = 5
const alvoCaras = 3


func main() {
    var wg sync.WaitGroup

	rand.Seed(time.Now().UnixNano())

	totalSimulacoes, err := strconv.ParseInt(os.Args[1], 10, 64)
    if err != nil{
        log.Fatal(err)
    }

	var sucessos atomic.Int64
    var split int64 = 1000
    partTotalSimulacoes, fracpartTotalSimulacoes := math.Modf(float64(totalSimulacoes)/float64(split))

    fmt.Printf("Executando %d rodadas de %.1f lancamentos + 1 rodada de %.1f lancamentos", split, partTotalSimulacoes, fracpartTotalSimulacoes)
    for i:= int64(1); i<= split; i++{
        wg.Add(1)
        go rodada( int64(partTotalSimulacoes), &sucessos, &wg)
    }
    if fracpartTotalSimulacoes > 0.0 {
        wg.Add(1)
        go rodada( int64(fracpartTotalSimulacoes), &sucessos, &wg)
    }

    wg.Wait()

	// Calculando a probabilidade estimada
	probabilidadeEstimada := float64(sucessos.Load()) / float64(totalSimulacoes)

	fmt.Printf("Simulou %d rodadas de %d lançamentos...\n", totalSimulacoes, lancamentosPorVez)
	fmt.Printf("Casos com exatamente %d caras: %d\n", alvoCaras, sucessos.Load())
	fmt.Printf("Probabilidade Estimada: %.8f (%.2f%%)\n", probabilidadeEstimada, probabilidadeEstimada*100)
	fmt.Println("Probabilidade Teórica: 0.3125 (31.25%)")
}
