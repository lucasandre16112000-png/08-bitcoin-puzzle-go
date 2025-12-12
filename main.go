package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	bip39 "github.com/tyler-smith/go-bip39"
)

const (
	TARGET_ADDRESS      = "1At7z8J3t3JJiAqtBTyJuHdCMKx45HmyVp"
	CHECKPOINT_FILE     = "checkpoint_v3_updated.json"
	CHECKPOINT_INTERVAL = 250000
	BATCH_SIZE          = 5000
)

type Checkpoint struct {
	Processed uint64 `json:"processed"`
	LastCombo string `json:"last_combo"`
	Timestamp string `json:"timestamp"`
}

type Result struct {
	Mnemonic   string `json:"mnemonic"`
	Address    string `json:"address"`
	PrivateKey string `json:"private_key_wif"`
	PublicKey  string `json:"public_key"`
}

var (
	processed uint64
	found     bool
	result    *Result
	mu        sync.Mutex
)

func loadBlocks() (map[int][]string, error) {
	blocks := make(map[int][]string)
	
	for i := 1; i <= 4; i++ {
		filename := fmt.Sprintf("block%d.txt", i)
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("erro ao abrir %s: %v", filename, err)
		}
		defer file.Close()
		
		var words []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			word := strings.TrimSpace(scanner.Text())
			if word != "" {
				words = append(words, word)
			}
		}
		blocks[i] = words
	}
	
	return blocks, nil
}

func permutations(arr []string, r int) [][]string {
	var result [][]string
	var helper func([]string, []string)
	
	helper = func(current []string, remaining []string) {
		if len(current) == r {
			perm := make([]string, r)
			copy(perm, current)
			result = append(result, perm)
			return
		}
		
		for i := 0; i < len(remaining); i++ {
			newCurrent := append(current, remaining[i])
			newRemaining := append([]string{}, remaining[:i]...)
			newRemaining = append(newRemaining, remaining[i+1:]...)
			helper(newCurrent, newRemaining)
		}
	}
	
	helper([]string{}, arr)
	return result
}

func testMnemonic(mnemonic string) *Result {
	// Validar BIP39
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil
	}
	
	// Gerar seed
	seed := bip39.NewSeed(mnemonic, "")
	
	// Derivar chave BIP44: m/44'/0'/0'/0/0
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil
	}
	
	// m/44'
	purpose, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 44)
	if err != nil {
		return nil
	}
	
	// m/44'/0'
	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return nil
	}
	
	// m/44'/0'/0'
	account, err := coinType.Derive(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return nil
	}
	
	// m/44'/0'/0'/0
	change, err := account.Derive(0)
	if err != nil {
		return nil
	}
	
	// m/44'/0'/0'/0/0
	addressKey, err := change.Derive(0)
	if err != nil {
		return nil
	}
	
	// Obter endereço
	address, err := addressKey.Address(&chaincfg.MainNetParams)
	if err != nil {
		return nil
	}
	
	if address.EncodeAddress() == TARGET_ADDRESS {
		// Encontrado!
		privKey, _ := addressKey.ECPrivKey()
		wif, _ := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, true)
		pubKey, _ := addressKey.ECPubKey()
		
		return &Result{
			Mnemonic:   mnemonic,
			Address:    address.EncodeAddress(),
			PrivateKey: wif.String(),
			PublicKey:  hex.EncodeToString(pubKey.SerializeCompressed()),
		}
	}
	
	return nil
}

func worker(id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for mnemonic := range jobs {
		if found {
			return
		}
		
		res := testMnemonic(mnemonic)
		if res != nil {
			mu.Lock()
			if !found {
				found = true
				result = res
				saveResult(res)
			}
			mu.Unlock()
			return
		}
		
		atomic.AddUint64(&processed, 1)
	}
}

func saveResult(res *Result) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("ENCONTRADO_%s.txt", timestamp)
	
	content := fmt.Sprintf(`================================================================================
🎉 CARTEIRA ENCONTRADA! 🎉
================================================================================

Mnemonic:
%s

Endereço:
%s

Chave Privada (WIF):
%s

Chave Pública:
%s

================================================================================
`, res.Mnemonic, res.Address, res.PrivateKey, res.PublicKey)
	
	ioutil.WriteFile(filename, []byte(content), 0644)
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🎉 CARTEIRA ENCONTRADA! 🎉")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\nMnemonic: %s\n", res.Mnemonic)
	fmt.Printf("Chave Privada: %s\n\n", res.PrivateKey)
	fmt.Printf("Salvo em: %s\n", filename)
	fmt.Println(strings.Repeat("=", 80) + "\n")
}

func saveCheckpoint(processed uint64) {
	checkpoint := Checkpoint{
		Processed: processed,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	
	data, _ := json.MarshalIndent(checkpoint, "", "  ")
	ioutil.WriteFile(CHECKPOINT_FILE, data, 0644)
}

func main() {
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("BITCOIN PUZZLE - V3 UPDATED (GO)")
	fmt.Println(strings.Repeat("=", 80))
	
	fmt.Printf("\nEndereço alvo: %s\n\n", TARGET_ADDRESS)
	
	// Carregar blocos
	blocks, err := loadBlocks()
	if err != nil {
		fmt.Printf("ERRO: %v\n", err)
		return
	}
	
	fmt.Println("Palavras carregadas:")
	var totalCombos uint64 = 1
	for i := 1; i <= 4; i++ {
		n := len(blocks[i])
		perms := n * (n - 1) * (n - 2)
		totalCombos *= uint64(perms)
		fmt.Printf("  Bloco %d: %d palavras → P(%d,3) = %d\n", i, n, n, perms)
	}
	
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Total: %d combinações\n", totalCombos)
	fmt.Println(strings.Repeat("=", 80))
	
	numWorkers := runtime.NumCPU()
	fmt.Printf("\nConfigurações:\n")
	fmt.Printf("  Workers: %d\n", numWorkers)
	fmt.Printf("  Batch size: %d\n", BATCH_SIZE)
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("Iniciando busca...")
	fmt.Println(strings.Repeat("=", 80) + "\n")
	
	// Gerar permutações de cada bloco
	blockPerms := make([][]string, 4)
	for i := 1; i <= 4; i++ {
		blockPerms[i-1] = make([]string, 0)
		perms := permutations(blocks[i], 3)
		for _, perm := range perms {
			blockPerms[i-1] = append(blockPerms[i-1], strings.Join(perm, " "))
		}
	}
	
	// Criar workers
	jobs := make(chan string, BATCH_SIZE)
	var wg sync.WaitGroup
	
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}
	
	// Gerar combinações e enviar para workers
	startTime := time.Now()
	lastCheckpoint := time.Now()
	_ = lastCheckpoint // avoid unused
	
	go func() {
		for !found {
			time.Sleep(1 * time.Second)
			current := atomic.LoadUint64(&processed)
			elapsed := time.Since(startTime).Seconds()
			speed := float64(current) / elapsed
			remaining := float64(totalCombos-current) / speed / 3600
			
			fmt.Printf("\r  %.1f%% | %d/%d | %.0f/s | %.1fh restantes",
				float64(current)/float64(totalCombos)*100,
				current, totalCombos, speed, remaining)
			
			// Checkpoint
			if time.Since(lastCheckpoint) > 5*time.Minute {
				saveCheckpoint(current)
				lastCheckpoint = time.Now()
			}
		}
	}()
	
	for _, p1 := range blockPerms[0] {
		if found {
			break
		}
		for _, p2 := range blockPerms[1] {
			if found {
				break
			}
			for _, p3 := range blockPerms[2] {
				if found {
					break
				}
				for _, p4 := range blockPerms[3] {
					if found {
						break
					}
					
					mnemonic := fmt.Sprintf("%s %s %s %s", p1, p2, p3, p4)
					jobs <- mnemonic
				}
			}
		}
	}
	
	close(jobs)
	wg.Wait()
	
	if !found {
		fmt.Println("\n\n❌ Não encontrado")
		fmt.Printf("Processado: %d mnemonics\n", atomic.LoadUint64(&processed))
	}
}
