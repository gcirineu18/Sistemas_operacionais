package main

type MFU struct {
	p *Pager
}

func (alg *MFU) executar() {

	aux := len(alg.p.trace)
	frameCapacity := cap(alg.p.frame)
	listaCircular := &ListaCircular{}

	var ponteiro *No =  nil
	for i := 0; i < aux; i++ {

		no := listaCircular.encontraNo(alg.p.trace[i])

		if no == nil && listaCircular.retornaTamanho() < frameCapacity {

			listaCircular.insereNoMFU(alg.p.trace[i])
			alg.p.pageFaults++

		} else{

			if no != nil{
				// Aumenta o acesso da página
				listaCircular.incrementarAcesso(alg.p.trace[i])
				continue
			} else{
				ponteiro = listaCircular.encontraProximaVitimaMFU(alg.p.trace[i], ponteiro)
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
func (alg *MFU) adicionarPaginasNovas() {

}