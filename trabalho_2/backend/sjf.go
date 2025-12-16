package main

import (
	"sort"
)

type SJF struct{
	s *Simulador
}

// adicionarProcessosNovos verifica se há processos novos chegando neste instante, se houver um novo processo, ordena os processos que não executaram ainda pelo tempo de duracao
func (alg *SJF) adicionarProcessosNovos() {
	for _, p := range alg.s.processos {
		// Se o processo chegou agora
		if p.instanteCriacao == alg.s.tempoAtual {
			alg.s.filaDeExecucao = append(alg.s.filaDeExecucao, p)
		}
	}

	// Ordena a fila de execução pelo tempo restante
	sort.Slice(alg.s.filaDeExecucao, func(i, j int) bool {
		return alg.s.filaDeExecucao[i].tempoRestante < alg.s.filaDeExecucao[j].tempoRestante
	})
}


// executar roda a simulação completa do escalonamento
func (alg *SJF) executar() {
	// Loop principal da simulação
	// Continua enquanto houver processos na fila
	// Adiciona processos que chegaram neste momento
	alg.adicionarProcessosNovos()
	for {	
		// Verifica se todos os processos já terminaram
		if len(alg.s.filaDeExecucao) == 0 && alg.s.verificarSeTerminou() {
			break // Todos os processos foram finalizados, podemos parar
		}

		// Se não há processos na fila, mas ainda há processos pendentes, avança o tempo
		if len(alg.s.filaDeExecucao) == 0 {
			// Registra tempo ocioso no diagrama
			alg.s.registrarDiagrama(nil)
			alg.s.tempoAtual++
			alg.adicionarProcessosNovos()
			continue
		}
		
		// Pega o primeiro processo da fila
		processoAtual := alg.s.filaDeExecucao[0]

		alg.s.filaDeExecucao = alg.s.filaDeExecucao[1:] // Remove da fila

		// Marca quando o processo iniciou pela primeira vez
		if processoAtual.tempoInicio == -1 {
			processoAtual.tempoInicio = alg.s.tempoAtual
		}

		// Conta troca de contexto (quando muda de um processo para outro)
		if alg.s.processoAnterior != nil && alg.s.processoAnterior != processoAtual {
			alg.s.trocasContexto++
		}
		alg.s.processoAnterior = processoAtual

		// Calcula quanto tempo o processo vai executar
		tempoExecucao := processoAtual.duracao

		// Executa o processo por tempoExecucao unidades de tempo
		for i := 0; i < tempoExecucao; i++ {
			
			alg.s.registrarDiagrama(processoAtual)
			alg.s.tempoAtual++
			processoAtual.tempoRestante--

			// // Durante a execução, podem chegar novos processos
			alg.adicionarProcessosNovos()

			// Se o processo terminou
			if processoAtual.tempoRestante == 0 {
				processoAtual.tempoTermino = alg.s.tempoAtual
				break
			}
		}
	}
}



