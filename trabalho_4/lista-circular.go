package main

type ListaCircular struct{
	primeiro *No
	ultimo *No
	tamanho int 
}

type No struct {
	bitRef bool
	page int
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
