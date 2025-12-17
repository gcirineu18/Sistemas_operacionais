package main

type ListaCircular struct{
	primeiro *No
	ultimo *No
	tamanho int 
}

type No struct {
	bitRef bool
	bitMod bool 			// Usado no NRU
	page int
	accessCount int 		// Usado no LFU
	proximo *No
}


func (lc *ListaCircular) insereNo(bitRef bool, page int){

	no := &No{bitRef: bitRef, page: page}

    if lc.primeiro == nil{
		lc.primeiro = no
		lc.ultimo = no
		no.proximo = no		
	} else{
		lc.ultimo.proximo = no
		lc.ultimo = no
		no.proximo = lc.primeiro
		
	}
	lc.tamanho++

}

func (lc *ListaCircular) insereNoLFU(page int){

	no := &No{bitRef: true, bitMod: false, page: page, accessCount: 1}

    if lc.primeiro == nil{
		lc.primeiro = no
		lc.ultimo = no
		no.proximo = no		
	} else{
		lc.ultimo.proximo = no
		lc.ultimo = no
		no.proximo = lc.primeiro
		
	}
	lc.tamanho++

}

// Incrementar acesso do FLU
func (lc *ListaCircular) incrementarAcesso(page int) {
	no := lc.primeiro

	for {
		if no.page == page {
			no.accessCount++
			break
		}

		no = no.proximo
		if no == lc.primeiro {
			break
		}
	}
}

func (lc *ListaCircular) retornaTamanho() int{
	return lc.tamanho
}

func (lc *ListaCircular) substituiNo(pagina int, pageToAdd int) {

	if lc.primeiro == nil{
		return 
	} 

	no := lc.primeiro
	for {
		if(no.proximo.page == pagina){

			novoNo := &No{bitRef: true, page: pageToAdd}
			if(no.proximo == lc.primeiro){
				lc.primeiro = novoNo
			} else if(no.proximo == lc.ultimo){
				lc.ultimo = novoNo
			}
			aux := no.proximo.proximo
			no.proximo = novoNo
			novoNo.proximo = aux
			break
		}
		no = no.proximo	

		if(no == lc.primeiro){
			break
		}
	} 
}

func (lc *ListaCircular) substituiNoNRU(pagina int, pageToAdd int) {
	if lc.primeiro == nil {
		return
	}

	no := lc.primeiro
	for {
		if no.page == pagina {
			novoNo := &No{bitRef: true, bitMod: false, page: pageToAdd}
			if no == lc.primeiro {
				lc.primeiro = novoNo
			} else if no == lc.ultimo {
				lc.ultimo = novoNo
			}

			aux := no.proximo
			no.proximo = novoNo
			novoNo.proximo = aux
			break
		}

		no = no.proximo
		if no == lc.primeiro {
			break
		}
	}
}

func (lc *ListaCircular) encontraNo(pagina int) *No{

	if lc.primeiro == nil{
		return nil
	} 

	no := lc.primeiro
	for {

		if(no.page == pagina){
			return no
		}
		no = no.proximo	

		if(no == lc.primeiro){
			break
		}
	} 

	return nil
}

func (lc *ListaCircular) encontraProximaVitima(pagina int, ponteiro *No) *No{
	if lc.primeiro == nil{
		return nil
	} 

	var no *No
	if(ponteiro == nil){
		no = lc.primeiro
	}else{
		no = ponteiro
	}

	for {
		if(!no.bitRef){
			lc.substituiNo(no.page, pagina)	
			return no.proximo		
		}
		no.bitRef = false
		no = no.proximo	

	} 

}

func (lc *ListaCircular) encontraProximaVitimaNRU(pagina int, ponteiro *No) *No{
	if lc.primeiro == nil {
		return nil
	}

	var no *No
	if ponteiro == nil {
		no = lc.primeiro
	} else {
		no = ponteiro
	}

	for {
		// Nível 00 - Melhores candidatas a substituição
		if !no.bitRef && !no.bitMod {
			lc.substituiNoNRU(no.page, pagina)
			return no.proximo
		}

		// Qualquer outro nível não é uma boa escolha pois ainda pode ser referenciada recenemente ou precisa salvar as alterações antes de ser subsituída
		// Vamos então apenas definir o bit de referência para permitir uma nova chance.
		no.bitRef = false
		no = no.proximo
	}
}

func (lc *ListaCircular) encontraProximaVitimaLFU(pagina int, ponteiro *No) *No{
	if lc.primeiro == nil {
		return nil
	}

	var no *No
	if ponteiro == nil {
		no = lc.primeiro
	} else {
		no = ponteiro
	}

	// Encontrar a página com menor número de acessos
	var menorPage *No = nil
	var menorCount = int(^uint(0) >> 1)

	for {
		if no.accessCount < menorCount {
			menorCount = no.accessCount
			menorPage = no
		}

		no = no.proximo
		if no == lc.primeiro {
			break
		}
	}

	
	if menorPage != nil {
		lc.substituiNo(menorPage.page, pagina)
	}

	return menorPage.proximo
}