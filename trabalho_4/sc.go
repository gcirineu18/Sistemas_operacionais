package main

type SC struct {
	p *Pager
}

func (alg *SC) executar() {

	aux := len(alg.p.trace)
	frameCapacity := cap(alg.p.frame)
	listaCircular := &ListaCircular{}

	var ponteiro *No =  nil
	for i := 0; i < aux; i++ {

		no := listaCircular.encontraNo(alg.p.trace[i])

		if no == nil && listaCircular.retornaTamanho() < frameCapacity {

			listaCircular.insereNo(true, alg.p.trace[i])
			alg.p.pageFaults++

		} else{

			if no != nil{
				no.bitRef = true
				continue
			} else{
				ponteiro = listaCircular.encontraProximaVitima(alg.p.trace[i], ponteiro)
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
func (alg *SC) adicionarPaginasNovas() {

}