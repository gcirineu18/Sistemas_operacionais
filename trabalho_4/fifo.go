package main

import (
	"slices"
)

type FIFO struct {
	p *Pager
}

func (alg *FIFO) executar() {

	aux := len(alg.p.trace)
	j:= 0 								// Índice auxiliar para gerenciar a substituição das páginas
	frameCapacity := cap(alg.p.frame)

	// Loop nas referências de páginas no trace
	for i := 0; i < aux; i++ {

		// Se a página não está presente nos frames e ainda há espaço, adiciona a página
		if !slices.Contains(alg.p.frame, alg.p.trace[i]) && len(alg.p.frame) < frameCapacity {
			alg.p.frame = append(alg.p.frame, alg.p.trace[i])
			alg.p.pageFaults++
		} else{

			// Se o índice de substituição ultrapassar a capacidade, reinicia
			if j > frameCapacity - 1{
				j = 0
			}

			// Se a página já está nos frames, ignora e continua para a próxima referência
			if slices.Contains(alg.p.frame, alg.p.trace[i]){
				continue
			} else{

				// Caso a página não esteja nos frames, substitui a página no índice j
				alg.p.frame[j] = alg.p.trace[i]			
				j++
				alg.p.pageFaults++
				alg.p.evictions++
			}		
		}
	}
}

func (alg *FIFO) adicionarPaginasNovas() {

}