 package main


 type RRPE struct{
	s *Simulador
 }

 // adicionarProcessosNovos verifica se há processos novos chegando neste instante
func (alg *RRPE) adicionarProcessosNovos() {
	for _, p := range alg.s.processos {
		// Se o processo chegou agora e ainda tem tempo restante (primeira vez na fila)
		if p.instanteCriacao == alg.s.tempoAtual && p.tempoRestante == p.duracao {
			alg.s.filaDeExecucao = append(alg.s.filaDeExecucao, p)
		}
	}
}


// executar roda a simulação completa do escalonamento
func (alg *RRPE) executar() {
	// Loop principal da simulação
	for {
		// Adiciona processos que chegaram neste momento
		alg.adicionarProcessosNovos()

		// Verifica se todos os processos já terminaram
		if len(alg.s.filaDeExecucao) == 0 && alg.s.verificarSeTerminou() {
			break // Todos os processos foram finalizados
		}

		// Se não há processos na fila, mas ainda há processos pendentes, avança o tempo
		if len(alg.s.filaDeExecucao) == 0 {
			alg.s.registrarDiagrama(nil)
			alg.s.tempoAtual++
			continue
		}

		// Ordena a fila por prioridade antes de escolher o próximo processo
		alg.s.ordenarFilaPorPrioridade()

		// Pega o processo de maior prioridade (primeiro da fila ordenada)
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

		// Reseta o contador de quantums esperando (o processo vai executar agora)
		processoAtual.quantunsEsperando = 0

		// Calcula quanto tempo o processo vai executar
		tempoExecucao := alg.s.quantum
		if processoAtual.tempoRestante < tempoExecucao {
			tempoExecucao = processoAtual.tempoRestante
		}

		// Executa o processo por tempoExecucao unidades de tempo
		for i := 0; i < tempoExecucao; i++ {
			alg.s.registrarDiagrama(processoAtual)
			alg.s.tempoAtual++
			processoAtual.tempoRestante--

			// Durante a execução, podem chegar novos processos
			alg.adicionarProcessosNovos()

			// Se o processo terminou
			if processoAtual.tempoRestante == 0 {
				processoAtual.tempoTermino = alg.s.tempoAtual
				break
			}
		}

		// Se o processo NÃO terminou no quantum
		if processoAtual.tempoRestante > 0 {
			// Restaura a prioridade original
			processoAtual.prioridadeAtual = processoAtual.prioridadeOriginal
			// Reinsere o processo na fila
			alg.s.filaDeExecucao = append(alg.s.filaDeExecucao, processoAtual)
		}

		// Aplica envelhecimento aos processos que ficaram esperando
		alg.s.aplicarEnvelhecimento()
	}
}
