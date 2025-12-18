package main

type LFU struct {
	p *Pager
}

func (alg *LFU) executar() {

	aux := len(alg.p.trace)
	frameCapacity := cap(alg.p.frame)
	listaCircular := &ListaCircular{}

	var ponteiro *No =  nil
	for i := 0; i < aux; i++ {

		no := listaCircular.encontraNo(alg.p.trace[i])

		if no == nil && listaCircular.retornaTamanho() < frameCapacity {

			listaCircular.insereNoLFU(alg.p.trace[i])
			alg.p.pageFaults++

		} else{

			if no != nil{
				// Aumenta o acesso da página
				listaCircular.incrementarAcesso(alg.p.trace[i])
				continue
			} else{
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