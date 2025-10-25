package main


 type RR struct{
	s *Simulador
 }



// adicionarProcessosNovos verifica se há processos novos chegando neste instante
func (alg *RR) adicionarProcessosNovos() {
	for _, p := range alg.s.processos {
		// Se o processo chegou agora e ainda tem tempo restante (primeira vez na fila)
		if p.instanteCriacao == alg.s.tempoAtual && p.tempoRestante == p.duracao {
			alg.s.filaDeExecucao = append(alg.s.filaDeExecucao, p)
		}
	}

}

// executar roda a simulação completa do escalonamento
func (alg *RR) executar() {
	// Loop principal da simulação
	// Continua enquanto houver processos na fila OU processos ainda não finalizados
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

		// Calcula quanto tempo o processo vai executar (quantum ou o que resta)
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

		// Se o processo ainda tem tempo restante, reinsere na fila
		if processoAtual.tempoRestante > 0 {
			alg.s.filaDeExecucao = append(alg.s.filaDeExecucao, processoAtual)
		}
	}
}

