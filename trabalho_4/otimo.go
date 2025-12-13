package main

import (
	"slices"
)

type Otimo struct {
	p *Pager
}

func (alg *Otimo) executar() {

	worst := -1
	aux := len(alg.p.trace)
	frameCapacity := cap(alg.p.frame)
	

	for i := 0; i < aux; i++ {
		inFrame:=slices.Contains(alg.p.frame, alg.p.trace[i])

		if !inFrame &&  len(alg.p.frame) < frameCapacity {
			alg.p.frame = append(alg.p.frame, alg.p.trace[i])
			alg.p.pageFaults++

		} else {

			if inFrame {				
				continue
			} else {
				calculaOtimo(alg.p.frame, i, &alg.p.trace, &worst)

				index := slices.Index(alg.p.frame, worst)
				alg.p.frame[index] = alg.p.trace[i]

				alg.p.pageFaults++
				alg.p.evictions++
			}
		}

		
	}
}

func calculaOtimo(optimumArr []int, linha int,  trace *[]int, worst *int){

	tamanhoEntrada := len(*trace)
	aux := make(map [int]int)
	var assigned bool
	
	for i:= 0 ; i < len(optimumArr); i++ {
		assigned = false
		for j := linha + 1; j < tamanhoEntrada; j++{
           if (optimumArr)[i] == (*trace)[j]{
				aux[(optimumArr)[i]] = j
				assigned = true
				break
		   }
		}
		if(!assigned){
			aux[(optimumArr)[i]] = tamanhoEntrada + 1
		}
	}
	*worst = -1
	maior := -1	
	for key, value := range aux{		
		if value > maior{
			maior = value
			*worst = key		
		}	
	}
		
}

func (alg *Otimo) adicionarPaginasNovas() {

}