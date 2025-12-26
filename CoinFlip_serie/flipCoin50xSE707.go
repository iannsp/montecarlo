package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	// 1. Configurações
	rand.Seed(time.Now().UnixNano())
	n := 50.0             // Lançamentos de moeda por experimento
	p := 0.5              // Probabilidade teórica de cara (50%)
	numExperimentos := 100000 // Quantas vezes vamos repetir o "balde de 50 moedas"

	// 2. Cálculo Teórico do Erro Padrão (Régua de Confiança)
	// SE = sqrt( p * (1-p) / n )
	seTeorico := math.Sqrt(p * (1 - p) / n)

	fmt.Printf("--- Teoria ---\n")
	fmt.Printf("Para n=50, o Erro Padrão teórico é: %.4f (%.2f%%)\n\n", seTeorico, seTeorico*100)

	// 3. Simulação
	proporcoes := make([]float64, numExperimentos)
	somaProporcoes := 0.0

	for i := 0; i < numExperimentos; i++ {
		caras := 0
		for j := 0; j < int(n); j++ {
			if rand.Float64() < p {
				caras++
			}
		}
		// Guardamos a proporção de caras deste experimento (ex: 0.48, 0.54...)
		proporcao := float64(caras) / n
		proporcoes[i] = proporcao
		somaProporcoes += proporcao
	}

	mediaDasProporcoes := somaProporcoes / float64(numExperimentos)

	// 4. Cálculo do Desvio Padrão dos Experimentos (O SE na prática)
	varianciaAmostral := 0.0
	for _, prop := range proporcoes {
		varianciaAmostral += math.Pow(prop-mediaDasProporcoes, 2)
	}
	desvioPadraoAmostral := math.Sqrt(varianciaAmostral / float64(numExperimentos))

	// 5. Resultado
	fmt.Printf("--- Prática (Simulação de Monte Carlo) ---\n")
	fmt.Printf("Média de todas as proporções: %.4f\n", mediaDasProporcoes)
	fmt.Printf("Desvio Padrão dos experimentos: %.4f (%.2f%%)\n", desvioPadraoAmostral, desvioPadraoAmostral*100)
	
	fmt.Printf("\nInsight: A variação real entre os experimentos foi de %.2f%%, ", desvioPadraoAmostral*100)
	fmt.Printf("provando o Erro Padrão de %.2f%%.\n", seTeorico*100)
}
