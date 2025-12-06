#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>
#define MAX_CMD_LEN 256
#define MAX_ARGS 32

// Cria variáveis globais
char *memoria_fisica = NULL;
int tamanho_memoria = 0;
int bloco_id_counter = 0;

// Função para fazer parsing da linha de comando digitada pelo usuário
void parse_command(char *input, char **args) {
    char* token = strtok(input, " ");
    int i = 0;
    // Separa a entrada em tokens (argumentos)
    while(token != NULL && i < MAX_ARGS - 1){
        args[i] = token;
        token = strtok(NULL, " ");
        i++;
    }
    args[i] = NULL;
}

void first_fit_allocate(int block_id, int size) {
    for (int i = 0; i <= tamanho_memoria - size; i++) {
        int j;
        // Verifica se há espaço suficiente a partir da posição i
        for (j = 0; j < size; j++) {
            if (memoria_fisica[i + j] != '.') {
                break; // Encontrou um bloco ocupado
            }
        }
        
        // Se encontrou espaço suficiente
        if (j == size) {
            // Aloca o bloco
            for (j = 0; j < size; j++) {
                memoria_fisica[i + j] = '0' + (block_id % 10); // marcar com dígito
            }
            printf("Alocado bloco %d de tamanho %d na posição %d\n", block_id, size, i);
            return;
        }
    }
    // Se não encontrou espaço suficiente
    printf("Falha ao alocar bloco %d de tamanho %d: memória insuficiente\n", block_id, size);
    bloco_id_counter--; // Reverte o incremento do ID do bloco  
}

void best_fit_allocate(int block_id, int size) {

}

void worst_fit_allocate(int block_id, int size) {

}

void freeid(int id) {
    // Libera toda a memória (marca todos os blocos como livres)
    for (int i = 0; i < tamanho_memoria; i++) {
        if (memoria_fisica[i] == '0' + (id % 10)) {
            memoria_fisica[i] = '.';
        }
    }
    printf("Memória liberada\n");
}

// Função que aloca memória de acordo com o espaço e o algoritmo escolhido
void allocate_memory(int size, char *alg) {
    if (strcmp(alg, "first") == 0) {
        first_fit_allocate(bloco_id_counter++, size);
    } else if (strcmp(alg, "best") == 0) {
        best_fit_allocate(bloco_id_counter++, size);
    } else if (strcmp(alg, "worst") == 0) {
        worst_fit_allocate(bloco_id_counter++, size);
    } else {
        printf("Algoritmo de alocação desconhecido: %s\n", alg);
    }
}


void execute_command(char **args) {
    // Comando "exit": finaliza o shell
    if (strcmp(args[0], "exit") == 0) {
        fflush(stdout);
        exit(0);
    }

    // Comando init X: Inicializa o vetor que simula a memória física e cria o primeiro bloco livre.
    if (strcmp(args[0], "init") == 0) {
        if (args[1] == NULL) {
            printf("Uso: init <tamanho_da_memoria>\n");
            return;
        }
        tamanho_memoria = atoi(args[1]);
        if (tamanho_memoria <= 0) {
            printf("Tamanho da memória deve ser um número positivo.\n");
            return;
        }

        memoria_fisica = malloc(tamanho_memoria * sizeof(char));
        if (memoria_fisica == NULL) {
            printf("Erro ao alocar memória.\n");
            return;
        }
        // Inicializa toda a memória como livre (representado por '.')
        for (int i = 0; i < tamanho_memoria; i++) {
            memoria_fisica[i] = '.';
        }
        return;
    } else {
        if (memoria_fisica == NULL) {
            printf("Memória não inicializada. Use o comando 'init <tamanho_da_memoria>' primeiro.\n");
            return;
        } else if (strcmp(args[0], "alloc") == 0) {
            if (args[1] == NULL || args[2] == NULL) {
                printf("Uso: alloc <tamanho_do_bloco> <algoritmo>\n");
                return;
            }
            int size = atoi(args[1]);
            char *alg = args[2];
            allocate_memory(size, alg);
        } else if (strcmp(args[0], "freeid") == 0) {
            if (args[1] == NULL) {
                printf("Uso: freeid <id_do_bloco>\n");
                return;
            }
            freeid(atoi(args[1]));
            // Libera o bloco com o ID especificado
        } else if (strcmp(args[0], "show") == 0) {
            printf("%s\n", memoria_fisica);
        } else {
            printf("Comando não reconhecido: %s\n", args[0]);
        }
    }
}

int main() {
    char input[MAX_CMD_LEN];
    char *args[MAX_ARGS];

    printf("Shell iniciado\n");
    printf("Digite 'exit' para sair\n\n");

    while (1) {
        printf("shell> ");
        fflush(stdout);
        // Ler entrada do usuário
        if (!fgets(input, sizeof(input), stdin)) {
            break;
        }
        
        // Remover quebra de linha
        input[strcspn(input, "\n")] = 0;

        // Ignorar linhas vazias
        if (strlen(input) == 0) {
            continue;
        }

        // Fazer parsing do comando
        parse_command(input, args);

        execute_command(args);
    }
    printf("Shell encerrado!\n");
    return 0;
}