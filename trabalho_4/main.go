package main

import
(
  "fmt"
  "os"
  "strconv"
  "bufio"
   "log"
)

type Pager struct{
	alg string
	trace []int
	frame int
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

	alg:=getAlgo(os.Args[argsMap["algo"]])
	frame:=getFrame(os.Args[argsMap["frames"]])
	trace:=getTrace(os.Args[argsMap["trace"]])

	return &Pager{alg: alg, frame: frame, trace: trace}
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


func getFrame(frame string) int{

  num, err := strconv.Atoi(frame)	
  if err != nil || num <= 0 {
	fmt.Println("O número de frames é inválido.")
	os.Exit(1)
  }
  return num
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

	linha := scanner.Text() // linha inteira como string

    num, err := strconv.Atoi(linha)
    if err != nil {
        log.Printf("%v", err)
        continue
    }

    fmt.Printf("num: %d, tipo: %[1]T\n", num)

	 pagesInput = append(pagesInput, num)

  }

  if err := scanner.Err(); err != nil{
	log.Fatal(err)
  }

  return pagesInput

}

func main(){

	if len(os.Args) < 7  || len(os.Args) > 8 {
		fmt.Printf("Número de argumentos inválidos. O Padrão deve ser:\n" +
		"go run main.go --algo <ALGO> --frames <N> --trace <arquivo> [--verbose]\n")
		os.Exit(1)
	}

	parseCommands()
}