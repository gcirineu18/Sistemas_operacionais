

// Função para quando clicar no botão
function handleOnSubmit() {
    const quantum = Number(document.getElementById('quantum').value)  ;
    const aging =  Number(document.getElementById('aging').value );
    
    const processData = document.getElementById('processData').value.trim();
    
    // ALgumas checagens de entradas
    if (!quantum || quantum <= 0) {
        alert("Por favor escolha um valor justo para o quantum");
        return;
    }

    if (!processData) {
        alert("Por favor forneça os dados no formato especificado");
        return;
    }

    // Tratar dados do textarea
    const processos = processData.split('\n').map(line => {
        const [begin, duration, priority] = line.split(' ')
        .map(v =>{ 
            console.log("-----------")
            if( isNaN(v) || v === '' || v== null || v == undefined) return -1
            return Number(v)
         } )

        return { begin, duration, priority };
    });

    const algoritmo = getSelectedAlgorithms();

    // Verificação se um algoritmo foi escolhido
    if (!algoritmo) {
        alert("Por favor escolha um algoritmo");
        return;
    }

    // Detalhe para escolha de aging quando algoritmo Round-Robin com envelhecimento for escolhido
    if ((!aging || aging <= 0) && algoritmo == "rrpe") {
        alert("Esse algoritmo necessita de um valor positivo para o aging");
        return;
    }

    const data = {
        alg: algoritmo,
        quantum: quantum,
        aging: aging,
        input: processos
    };

    // Verificar se os dados foram capturados corretamente
    //console.log('Dados JSON a ser enviado:', JSON.stringify(data, null, 2));

    enviarDados(data);
};

// Função para selecionar o algoritmo de acordo com os radius no html
function getSelectedAlgorithms() {
    const selectedAlg = document.querySelector('input[name="alg"]:checked');
    return selectedAlg ? selectedAlg.value : null; 
}

// Requisição para endpoint com dados processados anteriormente
function enviarDados(data) {
    fetch(`http://localhost:8081/processes`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    
    .then(response => {
        if(response.ok) return response.json()
        return response.json().then(response => {throw new Error(response.error, response.process)})    
    })
    .then(result => {
        exibirResultados(result);      
    })
    .catch((error) => {
       
        Swal.fire({
            icon: 'error',
            title: 'Erro ao enviar dados',
            text: `${error.message}`
        }) 
    });
}

// Função para tratar JSON de resposta
function exibirResultados(result) {
    // Mostrar valores das estatísticas no elemento html correspondente
    document.getElementById('avgTurnaround').textContent = `Tempo médio de vida (turnaround): ${result.tempoMedioVida}`;
    document.getElementById('avgWait').textContent = `Tempo médio de espera: ${result.tempoMedioEspera}`;
    document.getElementById('contextSwitches').textContent = `Número de trocas de contexto: ${result.trocasContexto}`;

    // Preparar para exibir diagrama de tempo
    // Vamos fazer um cabeçalho com os ids dos processos de acordo com os dados recebidos do backend
    // Utilizamos o padStart para formatar o alinhamento entre cabeçalho, tempo e status dos procesos
    let timeDiagram = '';

    const numProcessos = result.diagramaTempo[0].length;
    let header = ' '.padStart(6);
    for (let i = 0; i < numProcessos; i++) {
        header += result.ordemProcessos[i];
    }

    timeDiagram += `${header}\n`

    result.diagramaTempo.forEach((linha, index) => {
        timeDiagram += `${index}-${index + 1} `.padStart(6)
        timeDiagram +=  `${linha.join(' ')}\n`;
    });

    // Verificar resultado do Diagrama
    //console.log(timeDiagram)

    // Adicionar diagrama no elemento html e mudar o display para blokc - "fazendo aparecer"
    document.getElementById('timeDiagram').textContent = timeDiagram;

    document.getElementById('results').style.display = 'block';
}