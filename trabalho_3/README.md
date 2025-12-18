# Alocação de memória
Trabalho Prático da cadeira de Sistemas Operacionais (CK0234) da Universidade Federal Ceará, 2025.
### Equipe:
- Francisco Leonan Marques Carneiro - 539000
- Guilherme Garbosa Cirineu - 520404
- Isaac Mosiah Bandeira Maia de Maria - 497447
- Kauan Guilherme de Brito Soares - 537063

## Instruções para Execução:

### 1. Na pasta do trabalho, compilar o código:
```bash 
gcc main.c -o <nome_do_executavel>
```

### 2. Executar:
```bash 
./<nome_do_executavel>
```


### 3. Simular:
Comece iniciando a área de memória com:

```bash 
init <valor>
```
Após isso, você pode alocar novos blocos de memória com o algoritmo selecionado:

```bash 
alloc <valor> first
alloc <valor> best
alloc <valor> worst
```
Quando quiser, poderá verificar a situação da memória com o seguinte comando:

```bash
show
```
Além disso, é possível ver status de alocação com o comando:

```bash
stats
```