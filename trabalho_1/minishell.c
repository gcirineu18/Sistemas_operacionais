#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>
#define MAX_CMD_LEN 256
#define MAX_ARGS 32

// Variáveis globais para controle de processos em background
pid_t bg_processes[10];  // Array para armazenar PIDs de processos em background
int bg_count = 0;        // Contador de processos em background ativos
pid_t last_child_pid = 0; // Armazena PID do último processo filho

// Função para fazer parsing da linha de comando digitada pelo usuário
void parse_command(char *input, char **args, int *background) {
    char* token = strtok(input, " ");
    int i = 0;
    // Separa a entrada em tokens (argumentos)
    while(token != NULL && i < MAX_ARGS - 1){
        args[i] = token;
        token = strtok(NULL, " ");
        i++;
    }

    // Verifica se o último argumento é "&" (execução em background)
    if(i > 0 && strcmp(args[i-1], "&") == 0){
        *background = 1;      // Marca como processo background
        args[i-1] = NULL;     // Remove o "&" dos argumentos
    }
    else{
        *background = 0;      // Processo foreground
        args[i] = NULL;       // Termina array de argumentos
    }

}

// Adiciona um processo à lista de background
static void add_bg_process(pid_t pid) {
    if (bg_count < 10) {
        bg_processes[bg_count++] = pid;
        printf("[%d] %d\n", bg_count, pid);  // Exibe número do job e PID
        fflush(stdout);
    } else {
        fprintf(stderr, "Limite de processos em background atingido.\n");
    }
}

void execute_command(char **args, int background) {
    if (args == NULL || args[0] == NULL) return;

    // Executa este método para gerar uma cópia do processo pai atual
    // fazendo que ambos processo compartilhem o mesmo contexto por enquanto...
    pid_t pid = fork();

    // Se retornar -1, significa que algum erro ocorreu, como limite de 
    //processos, por exemplo
    if (pid < 0) {
        perror("fork");
        return;
    }

    if (pid == 0) {
        // --- Filho: substitui a imagem do processo pelo comando externo ---
        execvp(args[0], args);

        // Se chegou aqui, execvp falhou:
        perror("execvp");

        _exit(127); // 127 é código padrão para "command not found"/erro de execução
    }

    // --- Pai ---
    last_child_pid = pid;

    if (background) {
        // Não bloqueia; apenas registra o processo em background
        add_bg_process(pid);
        // Não dá wait aqui (evita bloquear). A limpeza pode ser feita
        // periodicamente com waitpid(..., WNOHANG) em uma função auxiliar.
    } else {
        int status;
        // Foreground: espera especificamente por esse filho
        if (waitpid(pid, &status, 0) < 0) {
            perror("waitpid");
            return;
        }

    }
}

int is_internal_command(char **args) {
   
    if(args[0] == NULL) return 0;

    // Retorna 1 caso o primeiro argumento da entrada
    // seja igual a um desses comandos internos
    return strcmp(args[0], "pid") == 0  ||
           strcmp(args[0], "exit") == 0 ||  
           strcmp(args[0], "wait") == 0 ||
           strcmp(args[0], "jobs") == 0 ;
}

void clean_finished_processes() {
    int status;
    pid_t pid;
    // WNOHANG = não bloqueia se nenhum processo terminou
    while ((pid = waitpid(-1, &status, WNOHANG)) > 0) {
        // Remove da lista de background
        for (int i = 0; i < bg_count; i++) {
            if (bg_processes[i] == pid) {
                printf("[%d]+ Done\n", i+1);
                // Remove elemento da lista
                for (int j = i; j < bg_count - 1; j++) {
                    bg_processes[j] = bg_processes[j+1];
                }
                bg_count--;
                break;
            }
        }
    }
}

void handle_internal_command(char **args) {

    // Comando "exit": finaliza o shell
    if (strcmp(args[0], "exit") == 0) {
        fflush(stdout);
        exit(0);
    }

    // Comando "pid": exibe PID do processo pai (shell) e do último filho executado
    if (strcmp(args[0], "pid") == 0) printf("PID pai: %d\nPID filho: %d\n",
         getpid(), last_child_pid
    );

    // Comando "wait": aguarda todos os processos em background terminarem
    if (strcmp(args[0], "wait") == 0) {
        printf("Aguardando processos em background...\n");
        // Loop até que todos os processos background terminem
        while (bg_count > 0) {
            int status;
            pid_t pid = wait(&status); // Bloqueia até um processo terminar
            // Remove o processo terminado da lista de background
            for (int i = 0; i < bg_count; i++) {
                if (bg_processes[i] == pid) {
                    printf("[%d]+ Done\n", i+1);
                    // Compacta o array removendo o elemento
                    for (int j = i; j < bg_count - 1; j++) {
                        bg_processes[j] = bg_processes[j+1];
                    }
                    bg_count--;
                    break;
                }
            }
        }
        printf("Todos os processos terminaram\n");
    }

    // Comando "jobs": lista todos os processos atualmente em execução em background
    if (strcmp(args[0], "jobs") == 0) {
        // Limpa processos que já terminaram mas ainda estão na lista.
        // Colocamos aqui também para não correr o risco do retorno deste comando
        // estar desatualizado
        clean_finished_processes();
        
        // Verifica se há processos em background para mostrar
        if (bg_count == 0) {
            printf("Nenhum processo em background\n");
        } else {
            printf("Processos em background:\n");
            // Itera por todos os processos em background e os exibe | Formato: [número_do_job] PID Status
            for (int i = 0; i < bg_count; i++) {
                printf("[%d] %d Running\n", i+1, bg_processes[i]);
            }
        }
    }

}

int main() {
    char input[MAX_CMD_LEN];
    char *args[MAX_ARGS];
    int background;

    printf("Mini-Shell iniciado (PID: %d)\n", getpid());
    printf("Digite 'exit' para sair\n\n");

    while (1) {
        clean_finished_processes();

        printf("minishell> ");
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
        parse_command(input, args, &background);

        // Executar comando
        if (is_internal_command(args)) {
            handle_internal_command(args);
        } else {
            execute_command(args, background);
        }
    }
        printf("Shell encerrado!\n");
        return 0;
}