package main

import
(
  "fmt"
  "os"
)

type Pager struct{
	alg string
	trace  string
	frames uint8
}


func parseCommands(){
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

	getAlgo(os.Args[argsMap["algo"]])
	getFrame(os.Args[argsMap["frames"]])
	getTrace(os.Args[argsMap["trace"]])
}

func getAlgo(alg string){

}
func getFrame(frame string){

}
func getTrace(trace string){

}

func main(){

	if len(os.Args) < 7  || len(os.Args) > 8 {
		fmt.Printf("Número de argumentos inválidos. O Padrão deve ser:\n" +
		"go run --algo <ALGO> --frames <N> --trace <arquivo> [--verbose]\n")
	}

	parseCommands()
}