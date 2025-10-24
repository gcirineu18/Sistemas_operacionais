let endpoint = ""

function handleOnSubmit() {
    const quantum = document.getElementById('quantum').value;
    const aging = document.getElementById('aging').value;
    
    const processData = document.getElementById('processData').value.trim();
    
    const processos = processData.split('\n').map(line => {
        const [begin, duration, priority] = line.split(' ').map(Number);
        return { begin, duration, priority };
    });

    const algoritmo = getSelectedAlgorithms();

    if (!algoritmo) {
        alert("Por favor escolha um algoritmo");
        return;
    }

    const data = {
        alg: algoritmo,
        quantum: quantum,
        aging: aging,
        inputs: {
            processo: processos
        }
    };

    console.log('Dados JSON a ser enviado:', JSON.stringify(data, null, 2));

    enviarDados(data);
};

function getSelectedAlgorithms() {
    const selectedAlg = document.querySelector('input[name="alg"]:checked');
    return selectedAlg ? selectedAlg.value : null; 
}

function enviarDados(data) {
    fetch(`http://localhost:8080/${endpoint}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(result => {
        exibirResultados(result);
    })
    .catch(error => {
        console.error('Erro ao enviar os dados:', error);
    });
}

function exibirResultados(result) {
    document.getElementById('avgTurnaround').textContent = `Tempo médio de vida (turnaround): ${result.tempoMedioVida}`;
    document.getElementById('avgWait').textContent = `Tempo médio de espera: ${result.tempoMedioEspera}`;
    document.getElementById('contextSwitches').textContent = `Número de trocas de contexto: ${result.trocasContexto}`;

    let timeDiagram = '';
    result.diagramaTempo.forEach((linha, index) => {
        timeDiagram += `${index}-${index + 1} ${linha.join(' ')}\n`;
    });
    document.getElementById('timeDiagram').textContent = timeDiagram;

    document.getElementById('results').style.display = 'block';
}