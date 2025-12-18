#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>
#define MAX_CMD_LEN 256
#define MAX_ARGS 32

typedef struct {
    int id;           // ID único do bloco
    int start;        // Posição inicial na memória
    int size;         // Tamanho do bloco em bytes
    int allocated;    // 1 = alocado, 0 = livre
} MemoryBlock;

// Cria variáveis globais
MemoryBlock *blocks = NULL;         // Array de structs
int blocks_count = 0;               // Quantidade de blocos
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

// Inicializa a memória simulada e o bloco livre inicial
void init_memory(char **args) {
    // Verifica se o tamanho da memória foi fornecido
    if (args[1] == NULL) {
        printf("Uso: init <tamanho_da_memoria>\n");
        return;
    }
    tamanho_memoria = atoi(args[1]);
    // Verifica se o tamanho é válido
    if (tamanho_memoria <= 0) {
        printf("Tamanho da memória deve ser um número positivo.\n");
        return;
    }

    // Aloca espaço para um único bloco inicial (representa toda a memória como livre)
    blocks = malloc(1 * sizeof(MemoryBlock));
    if (blocks == NULL) {
        printf("Erro ao alocar memória.\n");
        return;
    }

    // Inicializa o bloco único representando toda a memória como livre
    blocks[0].id = -1;              // ID -1 indica espaço livre
    blocks[0].start = 0;            // Começa na posição 0
    blocks[0].size = tamanho_memoria; // Tamanho é toda a memória
    blocks[0].allocated = 0;        // 0 = livre, 1 = alocado
    
    blocks_count = 1;              // Começamos com 1 bloco (o espaço livre)
    
    printf("Memória inicializada com %d bytes\n", tamanho_memoria);
}

// Insere um novo bloco na posição especificada, deslocando os blocos seguintes
void insert_block_at(int pos, MemoryBlock new_block) {
    blocks = realloc(blocks, (blocks_count + 1) * sizeof(MemoryBlock));
    memmove(&blocks[pos + 1], &blocks[pos], (blocks_count - pos) * sizeof(MemoryBlock));
    blocks[pos] = new_block;
    blocks_count++;
}

// Junta blocos livres adjacentes para reduzir fragmentação externa
static void coalesce_free_blocks() {
    if (blocks_count <= 1) return;
    int write = 0;
    // Percorre os blocos e mescla os livres contíguos
    for (int read = 0; read < blocks_count; read++) {
        if (write == 0) {
            blocks[write] = blocks[read];
            write++;
            continue;
        }
        MemoryBlock *previous = &blocks[write - 1];
        MemoryBlock *current = &blocks[read];
        // Se ambos são livres e contíguos, mesclar
        if (!previous->allocated && !current->allocated && (previous->start + previous->size == current->start)) {
            previous->size += current->size;
        } else {
            blocks[write] = *current;
            write++;
        }
    }
    // Ajusta o tamanho lógico; opcionalmente pode reduzir com realloc
    blocks_count = write;
}

// =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-= Algoritmos de alocação =-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

// Aloca memória usando o algoritmo First Fit
void first_fit_allocate(int block_id, int size) {
    for (int i = 0; i < blocks_count; i++) {
        // Verifica se o bloco está livre e é grande o suficiente
        if (!blocks[i].allocated && blocks[i].size >= size) {
            // Se o bloco é exatamente do tamanho necessário
            if (blocks[i].size == size){
                blocks[i].allocated = 1;
                blocks[i].id = block_id;
            } else { // Se o bloco é maior, divide-o
                MemoryBlock free_tail = {
                    .id = -1,
                    .start = blocks[i].start + size,
                    .size = blocks[i].size - size,
                    .allocated = 0
                };
                blocks[i].allocated = 1;
                blocks[i].id = block_id;
                blocks[i].size = size;
                insert_block_at(i + 1, free_tail);
            }
            printf("Alocado bloco %d de tamanho %d na posição %d\n", block_id, size, blocks[i].start);
            return;
        }
    }
    // Se não encontrou espaço suficiente
    printf("Falha ao alocar bloco %d de tamanho %d: memória insuficiente\n", block_id, size);
}

// Aloca memória usando o algoritmo Best Fit
void best_fit_allocate(int block_id, int size) {
    int best = -1;
    int best_size = 1e9; // Inicializa com um valor maior que qualquer buraco
    for (int i = 0; i < blocks_count; i++) { // Percorre todos os blocos
        if (!blocks[i].allocated && blocks[i].size >= size && blocks[i].size < best_size) {
            best = i;
            best_size = blocks[i].size;
        }
    } 
    if (best == -1) { // Nenhum bloco adequado encontrado
        printf("Falha ao alocar bloco %d de tamanho %d: memória insuficiente\n", block_id, size);
        return;
    }
    if (blocks[best].size == size) {
        blocks[best].allocated = 1;
        blocks[best].id = block_id;
    } else {
        MemoryBlock free_tail = {
            .id = -1,
            .start = blocks[best].start + size,
            .size = blocks[best].size - size,
            .allocated = 0
        };
        blocks[best].allocated = 1;
        blocks[best].id = block_id;
        blocks[best].size = size;
        insert_block_at(best + 1, free_tail);
    }
    printf("Alocado bloco %d de tamanho %d na posição %d\n", block_id, size, blocks[best].start);
    return;
}

void worst_fit_allocate(int block_id, int size) {
    int worst = -1;
    int worst_size = -1; // Inicializa com um valor maior que qualquer buraco
    for (int i = 0; i < blocks_count; i++) {
        if (!blocks[i].allocated && blocks[i].size >= size && blocks[i].size > worst_size) {
            worst = i;
            worst_size = blocks[i].size;
        }
    }
    if (worst == -1) { // Nenhum bloco adequado encontrado
        printf("Falha ao alocar bloco %d de tamanho %d: memória insuficiente\n", block_id, size);
        return;
    }
    // Aloca o bloco no espaço encontrado
    if (blocks[worst].size == size) {
        blocks[worst].allocated = 1;
        blocks[worst].id = block_id;
    } else { // Se o bloco é maior, divide-o
        MemoryBlock free_tail = {
            .id = -1,
            .start = blocks[worst].start + size,
            .size = blocks[worst].size - size,
            .allocated = 0
        };
        blocks[worst].allocated = 1;
        blocks[worst].id = block_id;
        blocks[worst].size = size;
        insert_block_at(worst + 1, free_tail);
    }
    printf("Alocado bloco %d de tamanho %d na posição %d\n", block_id, size, blocks[worst].start);
    return;
}

// Libera o bloco com o ID especificado
void freeid(int id) {
    int freed = 0;
    for (int i = 0; i < blocks_count; i++) {
        if (blocks[i].allocated && blocks[i].id == id) {
            blocks[i].allocated = 0;
            blocks[i].id = -1; // marca como livre
            freed += blocks[i].size;
        }
    }
    // Coalesce blocos livres após a liberação
    if (freed > 0) {
        coalesce_free_blocks();
        printf("Liberado id=%d (%d bytes) e coalescido blocos livres\n", id, freed);
    } else {
        printf("Nenhum bloco encontrado com id=%d\n", id);
    }
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

// Função auxiliar para imprimir uma linha repetida
void print_repeat(char c, int n) {
    for (int i = 0; i < n; i++) putchar(c);
    putchar('\n');
}

// Exibe o mapa de memória no formato solicitado
void show_memory() {
    printf("Mapa de Memória (%d bytes):\n", tamanho_memoria);
    print_repeat('-', tamanho_memoria + 2);
    
    // Reconstrói a visualização a partir das structs
    char *fisic_memory = malloc(tamanho_memoria + 1);
    char *allocated_ids = malloc(tamanho_memoria + 1);
    if (fisic_memory == NULL || allocated_ids == NULL) {
        printf("Erro ao alocar memória para visualização\n");
        return;
    }
    
    // Inicializa tudo como livre
    for (int i = 0; i < tamanho_memoria; i++) {
        fisic_memory[i] = '.';
        allocated_ids[i] = '.';
    }
    
    // Marca os blocos alocados
    for (int i = 0; i < blocks_count; i++) {
        if (blocks[i].allocated) {
            char marker = '0' + (blocks[i].id % 10);
            for (int j = 0; j < blocks[i].size; j++) {
                fisic_memory[blocks[i].start + j] = '#';
                allocated_ids[blocks[i].start + j] = marker;
            }
        }
    }
    
    // Finaliza as strings
    fisic_memory[tamanho_memoria] = '\0';
    allocated_ids[tamanho_memoria] = '\0';
    printf("[%s]\n", fisic_memory);
    printf("[%s]\n", allocated_ids);
    print_repeat('-', tamanho_memoria + 2);
    
    free(fisic_memory);
    free(allocated_ids);
}

// Exibe estatísticas e blocos ativos no formato solicitado
void show_stats() {
    int occupied = 0;
    int holes = 0;
    int allocated_blocks_count = 0;

    // Usa diretamente o array global blocks
    for (int i = 0; i < blocks_count; i++) {
        if (blocks[i].allocated) {
            occupied += blocks[i].size;
            allocated_blocks_count++;
        } else {
            holes++;  // Bloco livre = buraco
        }
    }

    // Calcular fragmentação interna (não temos metadados de pedido vs alocado -> 0)
    int internal_frag = 0;

    // Imprimir blocos ativos na mesma linha, separados por ' | '
    printf("Blocos ativos: ");
    int printed = 0;
    for (int i = 0; i < blocks_count; i++) {
        if (blocks[i].allocated) {
            if (printed > 0) printf(" | ");
            printf("[id=%d] @%d +%dB (usado=%dB)", blocks[i].id, blocks[i].start, blocks[i].size, blocks[i].size);
            printed++;
        }
    }
    if (printed == 0) printf("(nenhum)");
    printf("\n");

    // Estatísticas
    printf("== Estatísticas ==\n");
    printf("Tamanho total: %d bytes\n", tamanho_memoria);
    printf("Ocupado: %d bytes | Livre: %d bytes\n", occupied, tamanho_memoria - occupied);
    printf("Buracos (fragmentação externa): %d\n", holes);
    printf("Fragmentação interna: %d bytes\n", internal_frag);
    double uso = 0.0;
    if (tamanho_memoria > 0) uso = ((double)occupied / (double)tamanho_memoria) * 100.0;
    printf("Uso efetivo: %.2f%%\n", uso);
}

// Executa o comando baseado nos argumentos fornecidos
void execute_command(char **args) {
    // Comando "exit": finaliza o shell
    if (strcmp(args[0], "exit") == 0) {
        fflush(stdout);
        exit(0);
    }

    // Comando init X: Inicializa o vetor que simula a memória física e cria o primeiro bloco livre.
    if (strcmp(args[0], "init") == 0) {
        init_memory(args);
    } else {
        if (blocks == NULL) {
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
            show_memory();
        } else if (strcmp(args[0], "stats") == 0) {
            show_stats();
        } else {
            printf("Comando não reconhecido: %s\n", args[0]);
        }
    }
}

// Função principal do shell
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