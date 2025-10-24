package main

import (

	"fmt"
	"sort"

)

// Processo representa uma tarefa a ser executada
type Processo struct {
	id                 int // Identificador único do processo
	instanteCriacao    int // Momento em que o processo chega ao sistema
	duracao            int // Tempo total que o processo precisa para executar
	prioridadeOriginal int // Prioridade estática original (menor número = maior prioridade)
	prioridadeAtual    int // Prioridade dinâmica que muda com o envelhecimento
	tempoRestante      int // Quanto tempo ainda falta para o processo terminar
	tempoInicio        int // Momento em que o processo começou a executar pela primeira vez
	tempoTermino       int // Momento em que o processo terminou completamente
	quantunsEsperando  int // Quantos quantums o processo passou esperando desde a última execução
}

// Simulador gerencia toda a execução do escalonamento
type Simulador struct {
	processos        []*Processo   // Lista de todos os processos
	filaDeExecucao   []*Processo   // Fila de processos prontos para executar
	quantum          int           // Tamanho do quantum (tempo que cada processo pode executar)
	tempoAtual       int           // Relógio do simulador
	trocasContexto   int           // Contador de trocas de contexto
	diagramaTempo    [][]string    // Matriz para armazenar o diagrama de execução
	processoAnterior *Processo     // Guarda o último processo que executou
}

type Escalonador interface{
	executar()
	adicionarProcessosNovos()
}

// Dependendo do tipo escolhido, cria o escalonador e simulador correspondente
// detalhe para o rrpe que também passamos o valor do aging
func novoEscalonador(tipo string, s *Simulador, aging int) (Escalonador, error) {
    switch tipo{
	case "rrpe":
		return &RRPE{s, aging}, nil
	case "srtf":
		return &SRTF{s}, nil
	case "rr":
		return  &RR{s}, nil
	case "sjf":
		return &SJF{s}, nil	
	case "fcfs":
		return &FCFS{s}, nil
	case "psp":
		return &PSP{s}, nil		
	case "pcpp":
		return &PCPP{s}, nil		
	default:
		return  nil, fmt.Errorf("algoritimo inválido")	
	}
}

// lerArquivo lê o arquivo de entrada e cria os processos
func lerEntradas(body ContextBody) ([]*Processo, error) {
	
	var processos []*Processo

	id := 1

	// Lê linha por linha do arquivo
	for i := range body.Input {

		// Converte as strings para números inteiros
		instanteCriacao:= body.Input[i].Begin
		duracao:= body.Input[i].Duration
		prioridade:= body.Input[i].Priority 

		// Cria um novo processo com os dados lidos
		processo := &Processo{
			id:                 id,
			instanteCriacao:    instanteCriacao,
			duracao:            duracao,
			prioridadeOriginal: prioridade,
			prioridadeAtual:    prioridade, // Inicialmente, prioridade atual = original
			tempoRestante:      duracao,
			tempoInicio:        -1, // -1 indica que ainda não começou
			quantunsEsperando:  0,
			tempoTermino: -1,
		}
		processos = append(processos, processo)
		id++
	}

	if(processos== nil){
		return nil, fmt.Errorf("arquivo vazio")
	}

	// Ordena os processos por instante de criação
	sort.Slice(processos, func(i, j int) bool {
		return processos[i].instanteCriacao < processos[j].instanteCriacao
	})

	return processos, nil
}

// novoSimulador cria um novo simulador com os processos e quantum fornecidos
func novoSimulador(processos []*Processo, quantum int) *Simulador {
	return &Simulador{
		processos:      processos,
		filaDeExecucao: make([]*Processo, 0),
		quantum:        quantum,
		tempoAtual:     0,
		diagramaTempo:  make([][]string, 0),
	}
}


// ordenarFilaPorPrioridade ordena a fila de execução pela prioridade atual
// Maior número = maior prioridade (prioridade 5 é mais importante que prioridade 1)
// Em caso de empate na prioridade, mantém a ordem de chegada na fila (FIFO)
func (s *Simulador) ordenarFilaPorPrioridade() {
	sort.SliceStable(s.filaDeExecucao, func(i, j int) bool {
		// Ordena por prioridade atual (maior número = maior prioridade)
		return s.filaDeExecucao[i].prioridadeAtual > s.filaDeExecucao[j].prioridadeAtual
	})
}

// aplicarEnvelhecimento aumenta a prioridade dos processos que estão esperando
// A cada quantum de espera, a prioridade aumenta em 1 (número maior = mais prioritário)
func (s *Simulador) aplicarEnvelhecimento(aging int) {
	for _, p := range s.filaDeExecucao {
		p.quantunsEsperando++
		// A cada quantum esperando, aumenta o número da prioridade
		p.prioridadeAtual += aging
	}
}



// verificarSeTerminou verifica se todos os processos foram finalizados
func (s *Simulador) verificarSeTerminou() bool {
	// Verifica se ainda existem processos com tempo restante
	for _, p := range s.processos {
		if p.tempoRestante > 0 {
			return false // Ainda há processos não finalizados
		}
	}
	return true // Todos os processos terminaram
}


// registrarDiagrama registra no diagrama qual processo executou neste segundo
func (s *Simulador) registrarDiagrama(processoAtual *Processo) {
	linha := make([]string, len(s.processos))
	
	for i, p := range s.processos {
		if processoAtual != nil && p.id == processoAtual.id {
			linha[i] = "##" // Processo está executando
		} else if p.instanteCriacao <= s.tempoAtual && p.tempoRestante > 0 {
			linha[i] = "--" // Processo está esperando
		} else if p.tempoRestante == 0 && p.tempoTermino <= s.tempoAtual {
			linha[i] = "  " // Processo já terminou
		} else {
			linha[i] = "  " // Processo ainda não chegou
		}
	}
	
	s.diagramaTempo = append(s.diagramaTempo, linha)
}

// calcularEstatisticas calcula as métricas finais do escalonamento
func (s *Simulador) calcularEstatisticas() (float64, float64) {
	var somaTempoVida, somaTempoEspera float64

	for _, p := range s.processos {
		// Verifica se o processo realmente foi executado
		if p.tempoTermino > 0 && p.tempoInicio >= 0 {
			// Tempo de vida (turnaround) = tempo de término - instante de criação
			tempoVida := p.tempoTermino - p.instanteCriacao
			somaTempoVida += float64(tempoVida)

			// Tempo de espera = tempo de vida - duração de execução
			tempoEspera := tempoVida - p.duracao
			somaTempoEspera += float64(tempoEspera)
		}
	}

	numProcessos := float64(len(s.processos))
	return somaTempoVida / numProcessos, somaTempoEspera / numProcessos
}

// imprimirResultados exibe todos os resultados da simulação
func (s *Simulador) imprimirResultados() (float64, float64, int, [][]string, []string){
	tempoMedioVida, tempoMedioEspera := s.calcularEstatisticas()

	ordemProcess := make([]string, len(s.processos))
	for i, p := range s.processos {
		ordemProcess[i] =  fmt.Sprintf("P%d ", p.id)
	}

	return tempoMedioVida, tempoMedioEspera, s.trocasContexto, s.diagramaTempo, ordemProcess
}

func processScheduler(body ContextBody) (float64, float64, int, [][]string, []string){

	algoritmo := body.Alg
	quantum := body.Quantum
	aging:= body.Aging

	// Lê os processos do arquivo
	processos, err := lerEntradas(body)
	if err != nil {
		return 0, 0, 0, nil, nil
	}

	// Cria e executa o simulador
	simulador := novoSimulador(processos, quantum)
	scheduler, err:= novoEscalonador(algoritmo, simulador, aging)

	if err != nil {
		fmt.Printf("Erro ao criar escalonador: %v\n", err)
		return 0, 0, 0, nil,nil
	}

	scheduler.executar()
	return simulador.imprimirResultados()
}

