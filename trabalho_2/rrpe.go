package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
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

// lerArquivo lê o arquivo de entrada e cria os processos
func lerArquivo(nomeArquivo string) ([]*Processo, error) {
	arquivo, err := os.Open(nomeArquivo)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo: %v", err)
	}
	defer arquivo.Close()

	var processos []*Processo
	scanner := bufio.NewScanner(arquivo)
	id := 1

	// Lê linha por linha do arquivo
	for scanner.Scan() {
		linha := strings.TrimSpace(scanner.Text())
		if linha == "" {
			continue // Ignora linhas vazias
		}

		// Separa os valores da linha por espaços
		campos := strings.Fields(linha)
		if len(campos) != 3 {
			return nil, fmt.Errorf("formato inválido na linha: %s", linha)
		}

		// Converte as strings para números inteiros
		instanteCriacao, _ := strconv.Atoi(campos[0])
		duracao, _ := strconv.Atoi(campos[1])
		prioridade, _ := strconv.Atoi(campos[2])

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
		}
		processos = append(processos, processo)
		id++
	}

	// Ordena os processos por instante de criação
	sort.Slice(processos, func(i, j int) bool {
		return processos[i].instanteCriacao < processos[j].instanteCriacao
	})

	return processos, scanner.Err()
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

// adicionarProcessosNovos verifica se há processos novos chegando neste instante
func (s *Simulador) adicionarProcessosNovos() {
	for _, p := range s.processos {
		// Se o processo chegou agora e ainda tem tempo restante (primeira vez na fila)
		if p.instanteCriacao == s.tempoAtual && p.tempoRestante == p.duracao {
			s.filaDeExecucao = append(s.filaDeExecucao, p)
		}
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
func (s *Simulador) aplicarEnvelhecimento() {
	for _, p := range s.filaDeExecucao {
		p.quantunsEsperando++
		// A cada quantum esperando, aumenta o número da prioridade
		p.prioridadeAtual++
	}
}

// verificarSeTerminou verifica se todos os processos foram finalizados
func (s *Simulador) verificarSeTerminou() bool {
	for _, p := range s.processos {
		if p.tempoRestante > 0 {
			return false // Ainda há processos não finalizados
		}
	}
	return true // Todos os processos terminaram
}

// executar roda a simulação completa do escalonamento
func (s *Simulador) executar() {
	// Loop principal da simulação
	for {
		// Adiciona processos que chegaram neste momento
		s.adicionarProcessosNovos()

		// Verifica se todos os processos já terminaram
		if len(s.filaDeExecucao) == 0 && s.verificarSeTerminou() {
			break // Todos os processos foram finalizados
		}

		// Se não há processos na fila, mas ainda há processos pendentes, avança o tempo
		if len(s.filaDeExecucao) == 0 {
			s.registrarDiagrama(nil)
			s.tempoAtual++
			continue
		}

		// Ordena a fila por prioridade antes de escolher o próximo processo
		s.ordenarFilaPorPrioridade()

		// Pega o processo de maior prioridade (primeiro da fila ordenada)
		processoAtual := s.filaDeExecucao[0]
		s.filaDeExecucao = s.filaDeExecucao[1:] // Remove da fila

		// Marca quando o processo iniciou pela primeira vez
		if processoAtual.tempoInicio == -1 {
			processoAtual.tempoInicio = s.tempoAtual
		}

		// Conta troca de contexto (quando muda de um processo para outro)
		if s.processoAnterior != nil && s.processoAnterior != processoAtual {
			s.trocasContexto++
		}
		s.processoAnterior = processoAtual

		// Reseta o contador de quantums esperando (o processo vai executar agora)
		processoAtual.quantunsEsperando = 0

		// Calcula quanto tempo o processo vai executar
		tempoExecucao := s.quantum
		if processoAtual.tempoRestante < tempoExecucao {
			tempoExecucao = processoAtual.tempoRestante
		}

		// Executa o processo por tempoExecucao unidades de tempo
		for i := 0; i < tempoExecucao; i++ {
			s.registrarDiagrama(processoAtual)
			s.tempoAtual++
			processoAtual.tempoRestante--

			// Durante a execução, podem chegar novos processos
			s.adicionarProcessosNovos()

			// Se o processo terminou
			if processoAtual.tempoRestante == 0 {
				processoAtual.tempoTermino = s.tempoAtual
				break
			}
		}

		// Se o processo NÃO terminou no quantum
		if processoAtual.tempoRestante > 0 {
			// Restaura a prioridade original
			processoAtual.prioridadeAtual = processoAtual.prioridadeOriginal
			// Reinsere o processo na fila
			s.filaDeExecucao = append(s.filaDeExecucao, processoAtual)
		}

		// Aplica envelhecimento aos processos que ficaram esperando
		s.aplicarEnvelhecimento()
	}
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
func (s *Simulador) imprimirResultados() {
	tempoMedioVida, tempoMedioEspera := s.calcularEstatisticas()

	fmt.Printf("Tempo médio de vida (turnaround): %.2f\n", tempoMedioVida)
	fmt.Printf("Tempo médio de espera: %.2f\n", tempoMedioEspera)
	fmt.Printf("Número de trocas de contexto: %d\n", s.trocasContexto)
	fmt.Println("\nDiagrama de tempo:")

	// Cabeçalho do diagrama
	fmt.Print("tempo ")
	for _, p := range s.processos {
		fmt.Printf("P%d ", p.id)
	}
	fmt.Println()

	// Linhas do diagrama
	for i, linha := range s.diagramaTempo {
		fmt.Printf("%2d-%2d ", i, i+1)
		for _, estado := range linha {
			fmt.Printf("%s ", estado)
		}
		fmt.Println()
	}
}

func main() {
	// Verifica se os argumentos foram fornecidos
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run main.go <arquivo_entrada> <quantum>")
		fmt.Println("Exemplo: go run main.go processos.txt 2")
		os.Exit(1)
	}

	nomeArquivo := os.Args[1]
	quantum, err := strconv.Atoi(os.Args[2])
	if err != nil || quantum <= 0 {
		fmt.Println("Quantum deve ser um número inteiro positivo")
		os.Exit(1)
	}

	// Lê os processos do arquivo
	processos, err := lerArquivo(nomeArquivo)
	if err != nil {
		fmt.Printf("Erro ao ler arquivo: %v\n", err)
		os.Exit(1)
	}

	// Cria e executa o simulador
	simulador := novoSimulador(processos, quantum)
	simulador.executar()
	simulador.imprimirResultados()
}