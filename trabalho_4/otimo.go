package main

import (
	"slices"
)

type Otimo struct {
	p *Pager
}

func (alg *Otimo) executar() {

	worst := -1						// Variável para encontrar a página mais distante no futuro
	aux := len(alg.p.trace)
	frameCapacity := cap(alg.p.frame)
	
	// Loop nas referências de páginas no trace
	for i := 0; i < aux; i++ {
		
		// Verifica se a página está nos Frames
		inFrame:=slices.Contains(alg.p.frame, alg.p.trace[i])

		// Se a página não está nos frames e há espaço, adiciona a página
		if !inFrame &&  len(alg.p.frame) < frameCapacity {
			alg.p.frame = append(alg.p.frame, alg.p.trace[i])
			alg.p.pageFaults++

		} else {
			// Se a página já está nos frames, continua para a próxima referência
			if inFrame {				
				continue
			} else {
				// Se a página não está nos frames
				calculaOtimo(alg.p.frame, i, &alg.p.trace, &worst)

				// Substitui a página mais distante no futuro pela nova página
				index := slices.Index(alg.p.frame, worst)
				alg.p.frame[index] = alg.p.trace[i]

				alg.p.pageFaults++
				alg.p.evictions++
			}
		}

		
	}
}

// Função que calcula a página mais distante no futuro
func calculaOtimo(optimumArr []int, linha int,  trace *[]int, worst *int){

	tamanhoEntrada := len(*trace)
	aux := make(map [int]int)			// Mapeia as páginas para a próxima ocorrência no futuro
	var assigned bool
	
	for i:= 0 ; i < len(optimumArr); i++ {
		assigned = false

		// Verifica a próxima ocorrência da página no trace a partir da linha atual
		for j := linha + 1; j < tamanhoEntrada; j++{
           if (optimumArr)[i] == (*trace)[j]{
				aux[(optimumArr)[i]] = j
				assigned = true
				break
		   }
		}

		// Se a página não ocorre mais no futuro, atribui um valor grande para sua próxima ocorrência
		if(!assigned){
			aux[(optimumArr)[i]] = tamanhoEntrada + 1
		}
	}

	// Encontra a página com a maior distância no futuro
	*worst = -1
	maior := -1	
	for key, value := range aux{	
		
		// Se a página tem a maior distância no futuro, ela é a "mais distante"	
		if value > maior{
			maior = value
			*worst = key		
		}	
	}
		
}

func (alg *Otimo) adicionarPaginasNovas() {

}