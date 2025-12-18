package main

type SC struct {
	p *Pager
}

func (alg *SC) executar() {

	aux := len(alg.p.trace)
	frameCapacity := cap(alg.p.frame)
	listaCircular := &ListaCircular{}

	var ponteiro *No =  nil				// Ponteiro usado para encontrar a página a ser substituída

	// Loop nas referênias das páginas
	for i := 0; i < aux; i++ {

		no := listaCircular.encontraNo(alg.p.trace[i])

		// Se a página não está na lista e há espaço nos frames, insere a página com o bit de referência como verdadeiro
		if no == nil && listaCircular.retornaTamanho() < frameCapacity {

			listaCircular.insereNo(true, alg.p.trace[i])
			alg.p.pageFaults++

		} else{

			// Se a página já está presente nos frames, apenas atualiza o bit de referência para 1
			if no != nil{
				no.bitRef = true
				continue
			} else{
				// Se a página não está nos frames, encontra a página com "Second Chance" para ser substituída
				ponteiro = listaCircular.encontraProximaVitima(alg.p.trace[i], ponteiro)
				alg.p.pageFaults++
				alg.p.evictions++
			}		
		}
	}

	// Atualiza os frames com as páginas da lista circular
	novoNo := listaCircular.primeiro
	for {
		alg.p.frame = append(alg.p.frame, novoNo.page)
		novoNo = novoNo.proximo
		if novoNo == listaCircular.primeiro{
			break
		}
	}

}
func (alg *SC) adicionarPaginasNovas() {

}