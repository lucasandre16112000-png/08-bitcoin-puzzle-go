# 🧩 Bitcoin Puzzle Solver - Versão em Go

Um solver de alta performance para Bitcoin puzzles escrito em Go, capaz de testar de 60.000 a 120.000 combinações de mnemônicos BIP39 por segundo usando processamento paralelo.

## ✨ Funcionalidades

- **Alta Performance**: 60.000-120.000 combinações/segundo.
- **Paralelismo Massivo**: Utilização de CPU de 12 núcleos com Goroutines.
- **Validação BIP39**: Validação completa de mnemônicos BIP39.
- **Sistema de Checkpoint**: Tolerante a falhas com salvamento de progresso a cada 5 minutos.
- **Capacidade de Resumo**: Pause e retome a execução sem perda de dados.

## 🛠️ Tecnologias

- **Go 1.19+**: Linguagem de programação de sistemas de alta performance.
- **Goroutines**: Modelo de concorrência leve.
- **Bibliotecas de Criptografia**: Pacotes de criptografia nativos do Go.

## 📋 Guia de Instalação e Execução (Para Qualquer Pessoa)

### Pré-requisitos

1.  **Git**: [**Download aqui**](https://git-scm.com/downloads)
2.  **Go**: [**Download aqui**](https://go.dev/dl/) (versão 1.19+)

### Passo 1: Baixar o Projeto

```bash
git clone https://github.com/lucasandre16112000-png/08-bitcoin-puzzle-go.git
cd 08-bitcoin-puzzle-go
```

### Passo 2: Compilar o Projeto

```bash
# No Windows
go build -o bitcoin_puzzle.exe main.go

# No macOS ou Linux
go build -o bitcoin_puzzle main.go
chmod +x bitcoin_puzzle
```

### Passo 3: Executar o Solver

```bash
# No Windows
.\bitcoin_puzzle.exe

# No macOS ou Linux
./bitcoin_puzzle
```

### Passo 4: Monitorar o Progresso

- O terminal exibirá o progresso em tempo real.
- Se uma solução for encontrada, um arquivo `ENCONTRADO_*.txt` será criado.

## 👨‍💻 Autor

Lucas André S - [GitHub](https://github.com/lucasandre16112000-png)
