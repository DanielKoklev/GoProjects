package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const difficulty = 1

type Block struct {
	Index int
	Timestamp string
	Data int
	Hash string
	PreviousHash string
	Difficulty int
	Nonce string // number used once
}

var blockChain []Block

type Message struct {
	Data int
}

var mtx = &sync.Mutex{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := Block{}
		genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), "", difficulty, ""}
		spew.Dump(genesisBlock)
		mtx.Lock()
		blockChain = append(blockChain, genesisBlock)
		mtx.Unlock()
		}()
	log.Fatal(run())
}

func run() error {
	mux := makeMuxRouter()
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server is running on port: ", httpPort)
	s := &http.Server{
		Addr: ":" + httpPort,
		Handler: mux,
		ReadTimeout: 10*time.Second,
		WriteTimeout: 10*time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockChain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockChain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(blockChain, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var m Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mtx.Lock()
	newBlock := generateBlock(blockChain[len(blockChain) - 1], m.Data)
	mtx.Unlock()

	if isBlockValid(newBlock, blockChain[len(blockChain) - 1]) {
		blockChain = append(blockChain, newBlock)
		spew.Dump(blockChain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func isBlockValid(newBlock Block, oldBlock Block) bool {
	if oldBlock.Index + 1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PreviousHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.Data) + block.PreviousHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}


func generateBlock(oldBlock Block, Data int) Block{
	var newBlock Block
	t := time.Now()
	
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.Data = Data
	newBlock.PreviousHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

	for i := 0; ; i++ {
		hex := fmt.Sprintf("/%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(calculateHash(newBlock), "Doing more work!")
			time.Sleep(time.Second)
			continue
		} else {
			fmt.Println(calculateHash(newBlock), "The work is done!")
			newBlock.Hash = calculateHash(newBlock)
			break
		}
	}
	return newBlock
}