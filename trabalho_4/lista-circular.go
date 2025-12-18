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

// Novo no de acordo com o valor do bitMod
func (lc *ListaCircular) insereNoNRU(bitRef bool, bitMod bool, page int){

	no := &No{bitRef: bitRef, bitMod: bitMod, page: page}

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

// Novo no de acordo com o valor do accessCount
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

// Única alteração para o insereNoLFU é o valor do accessCount
func (lc *ListaCircular) insereNoMFU(page int){

	no := &No{bitRef: true, bitMod: false, page: page, accessCount: 0}

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

// Incrementar acesso do FLU e MFU
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

// Vai percorrer a lista para realizar a substituição pagina -> pageToAdd
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
			// O novoNo é integrado e ligado ao seu no anterior e posterior
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
    bitMod := false
	no := lc.primeiro
	for {
		if no.proximo.page == pagina {
			// Páginas pares terão o bit mod = 1
				if (pageToAdd % 2) == 0{
					bitMod = true
				} else{
					bitMod = false
				}
			novoNo := &No{bitRef: true, bitMod: bitMod, page: pageToAdd}
			if no.proximo == lc.primeiro {
				lc.primeiro = novoNo
			} else if no.proximo == lc.ultimo {
				lc.ultimo = novoNo
			}

			aux := no.proximo.proximo
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

// Busca o No correspondente a pagina recebida do Pager
// -- Essa função serve de base para as outras funções encontraNo dos demais algoritmos --
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

// Percorre a lista até encontrar uma página com bit de referência 0
// As páginas com bit 1 são atualizadas para 0 como uma segunda chance
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
	aux := no 
	count := 1				
	// count vai funcionar como uma variável de "prioridade" - onde 1 é a maior prioridade (Nível 00) e 4 é a menor (Nível 11)
	// Caso o algoritmo percorreu toda a lista e ainda não encontrou uma página para substituir ele passa para o critério com menor prioridade

	for {
		
		
		if !no.bitRef && !no.bitMod && count == 1 {
			// Nível 00 - Melhores candidatas a substituição
			lc.substituiNoNRU(no.page, pagina)
			return no.proximo
		} else if !no.bitRef && no.bitMod && count == 2{
			// Nível 01	
			lc.substituiNoNRU(no.page, pagina)
			return no.proximo
		}	else if no.bitRef && !no.bitMod && count == 3{
			// Nível 10
			lc.substituiNoNRU(no.page, pagina)
			return no.proximo
		} else if no.bitRef && no.bitMod && count == 4 {
			// Nível 11
			lc.substituiNoNRU(no.page, pagina)
			return no.proximo
		}		

		no.bitRef = false
		no = no.proximo
		
		// Toda a lista já foi percorrida, a prioridade é alterada
		if no == aux && count < 4{
			count++
		} else if count == 4{
			no = aux
			count = 1
		}
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


// A função segue a mesma lógica da LFU porém agora escolhendo a página com maior número de acessos
func (lc *ListaCircular) encontraProximaVitimaMFU(pagina int, ponteiro *No) *No{
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
	var maiorPage *No = nil
	var maiorCount = -1

	for {
		if no.accessCount > maiorCount {
			maiorCount = no.accessCount
			maiorPage = no
		}

		no = no.proximo
		if no == lc.primeiro {
			break
		}
	}

	
	if maiorPage != nil {
		lc.substituiNo(maiorPage.page, pagina)
	}

	return maiorPage.proximo
}