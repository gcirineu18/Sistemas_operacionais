package main

type LFU struct {
	p *Pager
}

func (alg *LFU) executar() {

	aux := len(alg.p.trace)
	frameCapacity := cap(alg.p.frame)
	listaCircular := &ListaCircular{}

	var ponteiro *No =  nil				// Ponteiro usado para encontrar a página a ser substituída (vitima)

	// Loop nas referências de páginas no trace
	for i := 0; i < aux; i++ {

		// Verifica se a página já está na lista circular
		no := listaCircular.encontraNo(alg.p.trace[i])

		// Se a página não está na lista e há espaço nos frames, insere a página na lista
		if no == nil && listaCircular.retornaTamanho() < frameCapacity {

			listaCircular.insereNoLFU(alg.p.trace[i])
			alg.p.pageFaults++

		} else{

			// Se a página já está presente nos frames, incrementa o contador de acessos dessa página
			if no != nil{
				listaCircular.incrementarAcesso(alg.p.trace[i])
				continue
			} else{
				
				// Se a página não está presente, encontra a página com menor frequência de acesso (LFU) para ser substituída
				ponteiro = listaCircular.encontraProximaVitimaLFU(alg.p.trace[i], ponteiro)
				alg.p.pageFaults++
				alg.p.evictions++
			}		
		}
	}

	novoNo := listaCircular.primeiro
	for {
		alg.p.frame = append(alg.p.frame, novoNo.page)
		novoNo = novoNo.proximo
		if novoNo == listaCircular.primeiro{
			break
		}
	}

}
func (alg *LFU) adicionarPaginasNovas() {

}