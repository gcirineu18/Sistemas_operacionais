package main

import (
	"slices"
)

type FIFO struct {
	p *Pager
}

func (alg *FIFO) executar() {

	aux := len(alg.p.trace)
	j:= 0
	frameCapacity := cap(alg.p.frame)

	for i := 0; i < aux; i++ {
		
		if !slices.Contains(alg.p.frame, alg.p.trace[i]) && len(alg.p.frame) < frameCapacity {
			alg.p.frame = append(alg.p.frame, alg.p.trace[i])
			alg.p.pageFaults++
		} else{
			if j > frameCapacity - 1{
				j = 0
			}

			if slices.Contains(alg.p.frame, alg.p.trace[i]){
				continue
			} else{
				alg.p.frame = alg.p.frame[1:] 
				alg.p.frame = append(alg.p.frame, alg.p.trace[i])
				j++
				alg.p.pageFaults++
				alg.p.evictions++
			}
			
		}

	}

}

func (alg *FIFO) adicionarPaginasNovas() {

}