#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/types.h>
#include <sys/wait.h>
#define MAX_CMD_LEN 256
#define MAX_ARGS 32

// Array para armazenar PIDs de processos em background
pid_t bg_processes[10];
int bg_count = 0;
pid_t last_child_pid = 0; // Armazena PID do último processo filho

void parse_command(char *input, char **args, int *background) {
    char* token = strtok(input, " ");
    int i = 0;
    while(token != NULL && i < MAX_ARGS - 1){
        args[i] = token;
        token = strtok(NULL, " ");
        i++;
    }

    if(i > 0 && strcmp(args[i-1], "&") == 0){
        *background = 1;
        args[i-1] = NULL;
    }
    else{
        *background = 0; 
        args[i] = NULL;
    }

}

void execute_command(char **args, int background) {
    // TODO: Implementar execução
    // Gerenciar background se necessário

    if(args[0] == NULL) return;
        
    int retval = 0;
    retval = fork();

    if(retval < 0){
        perror("Error: ");
        exit(1);
    } 

    else if(!background && retval > 0 ){
        last_child_pid = retval;
        wait(0);
    }

    else if(background && retval > 0 ){
        // TODO --> Processo em background
    }

    else{
        if (execvp(args[0], args) == -1) {
            perror("Erro"); 
            exit(0);
        }
        else{
           execvp(args[0], args); 
           return;
        }      
    }

}

int is_internal_command(char **args) {
   
    if(args[0] == NULL) return 0;

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

    // clean all processes and close mini-shell
    if (strcmp(args[0], "exit") == 0) {
        fflush(stdout);
        exit(0);
    }

    if (strcmp(args[0], "pid") == 0) printf("PID pai: %d\nPID filho: %d\n",
         getpid(), last_child_pid
    );

    // TODO: tratar para os outros comandos

    if (strcmp(args[0], "wait") == 0) {
        printf("Aguardando processos em background...\n");
        while (bg_count > 0) {
            int status;
            pid_t pid = wait(&status); // Bloqueia até um processo terminar
            // Remove da lista (código similar ao clean_finished_processes)
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
        printf("Todos os processos terminaram\n");
    }

    if (strcmp(args[0], "jobs") == 0) {
        if (bg_count == 0) {
            printf("Nenhum processo em background\n");
        } else {
            printf("Processos em background:\n");
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