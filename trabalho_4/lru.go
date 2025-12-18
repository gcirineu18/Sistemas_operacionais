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
	var frequencyArr []int

	for i := 0; i < aux; i++ {

		if !slices.Contains(alg.p.frame, alg.p.trace[i]) &&  len(alg.p.frame) < frameCapacity {
			alg.p.frame = append(alg.p.frame, alg.p.trace[i])
			alg.p.pageFaults++

		} else {

			if slices.Contains(alg.p.frame, alg.p.trace[i]) {
				calculaLRU(alg.p.trace[i], &frequencyArr)
				continue
			} else {
				leastRecent := frequencyArr[0]
				index := slices.Index(alg.p.frame, leastRecent)
				alg.p.frame[index] = alg.p.trace[i]

				frequencyArr = frequencyArr[1:]

				alg.p.pageFaults++
				alg.p.evictions++
			}
		}

		calculaLRU(alg.p.trace[i], &frequencyArr)
	}
}

func calculaLRU(mostRecent int, frequencyArr *[]int){
	if !slices.Contains(*frequencyArr, mostRecent){
 		*frequencyArr= append(*frequencyArr, mostRecent)
	} else{
		index := slices.Index(*frequencyArr, mostRecent)
		*frequencyArr = append((*frequencyArr)[:index],(*frequencyArr)[index+1:]...)
		*frequencyArr = append(*frequencyArr, mostRecent)
	}
	
}

func (alg *LRU) adicionarPaginasNovas() {

}