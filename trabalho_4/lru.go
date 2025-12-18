package main

import (
	"slices"
)

type LRU struct {
	p *Pager
}

func (alg *LRU) executar() {

	aux := len(alg.p.trace)
	frameCapacity := cap(alg.p.frame)
	var frequencyArr []int				// Array para armazenar a ordem de acesso das páginas

	// Loop nas referências de páginas no trace
	for i := 0; i < aux; i++ {

		// Se a página não está nos frames e ainda há espaço, adiciona a página
		if !slices.Contains(alg.p.frame, alg.p.trace[i]) &&  len(alg.p.frame) < frameCapacity {
			alg.p.frame = append(alg.p.frame, alg.p.trace[i])
			alg.p.pageFaults++

		} else {

			// Se a página já está nos frames, apenas atualiza sua posição
			if slices.Contains(alg.p.frame, alg.p.trace[i]) {
				calculaLRU(alg.p.trace[i], &frequencyArr)
				continue
			} else {

				// Se a página não está nos frames, encontra a página menos recentemente usada (LRU)
				leastRecent := frequencyArr[0]
				index := slices.Index(alg.p.frame, leastRecent)
				alg.p.frame[index] = alg.p.trace[i]

				// Remove a página substituída da lista de frequência e atualiza com a nova página
				frequencyArr = frequencyArr[1:]

				alg.p.pageFaults++
				alg.p.evictions++
			}
		}

		calculaLRU(alg.p.trace[i], &frequencyArr)
	}
}

// Função para atualizar a ordem de acesso das páginas
func calculaLRU(mostRecent int, frequencyArr *[]int){

	// Se a página não estiver na lista de frequência, adiciona ela no final
	if !slices.Contains(*frequencyArr, mostRecent){
 		*frequencyArr= append(*frequencyArr, mostRecent)
	} else{

		// Se a página já estiver na lista, remove ela e coloca no final
		index := slices.Index(*frequencyArr, mostRecent)
		*frequencyArr = append((*frequencyArr)[:index],(*frequencyArr)[index+1:]...)
		*frequencyArr = append(*frequencyArr, mostRecent)
	}
	
}

func (alg *LRU) adicionarPaginasNovas() {

}