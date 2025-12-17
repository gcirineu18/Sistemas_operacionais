package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)


const(
	fifo = "fifo"
	lru = "lru"
	otimo = "otimo"
	sc = "sc"
	clock = "clock" 
	nru = "nru"
	lfu = "lfu"
	mfu = "mfu"
)

var algName = map[string]string{
	fifo: "FIFO",
	lru: "LRU",
	otimo: "Ótimo",
	sc: "Second Chance",
	clock: "Clock",
	nru: "NRU",
	lfu: "LFU",
	mfu: "MFU" ,
}

//Faz o mapeamento
func  String(algorithm string) string{
	return algName[algorithm]
}


type Pager struct{
	alg string
	trace []int
	frame []int
	evictions int 
	pageFaults int
}

type Simulador interface {
	executar()
	adicionarPaginasNovas()
}


func parseCommands() *Pager{
	argsMap := make(map[string]int)
	argsMap["algo"] = 0
	argsMap["frames"] = 0
	argsMap["trace"] = 0

	for i := range len(os.Args){
	  switch os.Args[i] {
		case "--algo":
			argsMap["algo"] = i + 1
		case "--frames":
			argsMap["frames"] = i + 1
		case "--trace":
			argsMap["trace"] = i + 1
		} 		
	}

	for key, val := range argsMap{
		if(val == 0){
			fmt.Printf("Flag %s não encontrada, por favor, tente novamente.\n", key)
			os.Exit(1)
		}
	}

	alg:= getAlgo(os.Args[argsMap["algo"]])
	frame:= getFrame(os.Args[argsMap["frames"]])
	trace:= getTrace(os.Args[argsMap["trace"]])

	return &Pager{alg: alg, frame: frame, trace: trace, evictions: 0, pageFaults: 0}
}

func getAlgo(alg string) string{
	aux:= 0 
	switch alg {
		case "fifo", "lru", "otimo", "sc", "clock", "nru", "lfu", "mfu":
			aux = 1
	}
		
	if aux == 0 {
		fmt.Printf("Nome de algoritmo inválido: %s\n", alg)
		os.Exit(1)
	}		
	return alg	
}


func getFrame(frame string) []int{

  num, err := strconv.Atoi(frame)	
  if err != nil || num <= 0 {
	fmt.Println("O número de frames é inválido. Deve ser maior que 0.")
	os.Exit(1)
  }
  return make([]int, 0, num)
}

func getTrace(trace string) []int{
  file, err := os.Open(trace)

  if  err != nil {
	log.Fatal(err)
  }

  defer file.Close()
  
  var pagesInput []int
  scanner:= bufio.NewScanner(file)
  for scanner.Scan(){

	linha := scanner.Text()

    num, err := strconv.Atoi(linha)
    if err != nil {
        log.Printf("%v", err)
        continue
    }

	pagesInput = append(pagesInput, num)

  }

  if err := scanner.Err(); err != nil{
	log.Fatal(err)
  }

  return pagesInput

}

func novoSimulador(p *Pager) (Simulador, error){
	
	switch p.alg {
	case "fifo": 
	      return &FIFO{p}, nil;
	case "lru":
		  return &LRU{p}, nil;  
	case "otimo":
		return &Otimo{p}, nil;
	case "sc":
		return &SC{p}, nil;
	case "nru":
		return &NRU{p}, nil;
	case "lfu":
		return &LFU{p}, nil;
	case "mfu":
		return &MFU{p}, nil;
	default:
		return nil,	fmt.Errorf("algoritimo inválido")
	}
}

func calculaEstatisticas(pager *Pager){

	taxa := float32(pager.pageFaults) * 100  / float32(len(pager.trace))
	fmt.Println("Algoritmo: ", String(pager.alg))
	fmt.Println("Frames: ", len(pager.frame))
	fmt.Println("Referências: ", len(pager.trace))
	fmt.Println("Faltas de página: ", pager.pageFaults)
	fmt.Printf("Taxa de faltas: %.2f%% \n", taxa)
	fmt.Println("Evicções: ", pager.evictions)
	fmt.Println("Conjunto residente final:")
	fmt.Printf("frame_ids:")
	for i:= 0; i < len(pager.frame); i++{
		fmt.Printf("  %d", i)
	}

	fmt.Printf("\npage_ids: ")
	for i:= 0; i < len(pager.frame); i++{
		fmt.Printf("  %d", pager.frame[i])
	}
	fmt.Println()
	
}

func main(){

	start:= time.Now()
	if len(os.Args) < 7  || len(os.Args) > 8 {
		fmt.Printf("Número de argumentos inválidos. O Padrão deve ser:\n" +
		"go run . --algo <ALGO> --frames <N> --trace <arquivo> [--verbose]\n")
		os.Exit(1)
	}

	pager := parseCommands()

	sim, err:= novoSimulador(pager)
	
	if err != nil {
	    log.Fatal("Erro ao instanciar o Simulador")
	}

	fmt.Println(pager.trace)

	sim.executar()

	calculaEstatisticas(pager)

	elapsed := time.Since(start)

	fmt.Println("Tempo de execução: ", elapsed)


}