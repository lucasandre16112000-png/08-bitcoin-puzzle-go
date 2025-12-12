# Bitcoin Puzzle Solver - Go Version

A high-performance Bitcoin puzzle solver written in Go, capable of testing 60,000-120,000 BIP39 mnemonic combinations per second using parallel processing.

## ✨ Features

- **High Performance**: 60,000-120,000 combinations/second
- **Massive Parallelism**: 12-core CPU utilization with Goroutines
- **BIP39 Validation**: Full BIP39 mnemonic validation
- **BIP44 Derivation**: Standard Bitcoin key derivation path (m/44'/0'/0'/0/0)
- **Checkpoint System**: Fault-tolerant with progress saving every 5 minutes
- **Resume Capability**: Pause and resume execution without data loss
- **Optimized Cryptography**: Efficient ECDSA and PBKDF2 implementation

## 🛠️ Technologies

- **Go 1.19+**: High-performance systems programming language
- **Goroutines**: Lightweight concurrency model
- **Crypto Libraries**: Go's built-in cryptography packages
- **JSON**: Checkpoint serialization

## 📦 Installation

### Prerequisites
- Go 1.19 or higher
- Linux, macOS, or Windows

### Setup

1. Clone the repository:
```bash
git clone https://github.com/lucasandre16112000-png/08-bitcoin-puzzle-go.git
cd 08-bitcoin-puzzle-go
```

2. Build the executable:

**Linux/macOS:**
```bash
go build -o bitcoin_puzzle main.go
chmod +x bitcoin_puzzle
```

**Windows:**
```bash
go build -o bitcoin_puzzle.exe main.go
```

## 🚀 Running the Solver

**Linux/macOS:**
```bash
./bitcoin_puzzle
```

**Windows:**
```bash
bitcoin_puzzle.exe
```

## 📊 Performance Metrics

| Metric | Value |
|--------|-------|
| Speed | 60,000-120,000 combinations/sec |
| CPU Usage | 90-100% (normal) |
| Memory | Efficient streaming |
| Estimated Time | 6-12 months |
| Total Combinations | ~1.7 trillion |

## 🎯 Puzzle Blocks

The solver attempts combinations from 4 word blocks:

- **Block 1** (12 words): age, ancient, bench, chair, donkey, elder, man, mule, old, seat, senior, stool
- **Block 2** (12 words): cross, earth, galaxy, lunar, moon, orbit, planet, solar, space, sun, universe, world
- **Block 3** (15 words): beef, cruel, evil, food, fun, giggle, glad, happy, horror, joy, laugh, meat, monster, smile, steak
- **Block 4** (24 words): adult, boy, business, child, company, couple, daughter, deal, family, father, firm, girl, human, kid, man, mother, pair, parent, people, person, trade, venture, woman, young

## 🔄 How It Works

1. **Permutation Generation**: Generates all possible word combinations
2. **BIP39 Validation**: Validates checksum and word list
3. **Seed Generation**: Creates seed from mnemonic using PBKDF2
4. **Key Derivation**: Derives Bitcoin private key using BIP44 path
5. **Address Generation**: Generates Bitcoin address from key
6. **Target Matching**: Compares with target address
7. **Checkpoint**: Saves progress periodically

## 💾 Output Files

- `ENCONTRADO_*.txt`: Found mnemonic and private key
- `checkpoint_*.json`: Progress checkpoint for resume

## 🔒 Security

- Private keys stored locally only
- No network communication
- Secure random number generation
- Proper cryptographic libraries

## 📁 Project Structure

```
.
├── main.go           # Main solver implementation
├── go.mod            # Go module definition
├── go.sum            # Dependency checksums
└── README.md         # This file
```

## 📄 License

This project is open source and available under the MIT License.

## 👨‍💻 Author

Lucas André - [GitHub](https://github.com/lucasandre16112000-png)

## ⚠️ Disclaimer

This tool is for educational and research purposes only. Use at your own risk.
