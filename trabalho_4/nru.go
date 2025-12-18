package main

type NRU struct {
	p *Pager
}

func (alg *NRU) executar() {

	aux := len(alg.p.trace)
	frameCapacity := cap(alg.p.frame)
	listaCircular := &ListaCircular{}
	mod:= false
	var ponteiro *No =  nil
	for i := 0; i < aux; i++ {

		no := listaCircular.encontraNo(alg.p.trace[i])


		if no == nil && listaCircular.retornaTamanho() < frameCapacity {
			// Páginas pares terão o bit mod = 1
				if (alg.p.trace[i] % 2) == 0{
					mod = true
				} else{
					mod = false
				}

			listaCircular.insereNoNRU(true, mod, alg.p.trace[i])

			alg.p.pageFaults++

		} else{

			if no != nil{
				no.bitRef = true				
				continue
			} else{
				ponteiro = listaCircular.encontraProximaVitimaNRU(alg.p.trace[i], ponteiro)
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
func (alg *NRU) adicionarPaginasNovas() {

}