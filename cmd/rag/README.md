# intro
利用 ollama 跟 qdrant 來進行 RAG 的 chat

# compile
```shell
go build -o bin/rag cmd/rag/main.go
```

# use
## create
```shell
./bin/rag create test
```

## import
```shell
path_input=/root/user_pvc/mygit/annotator-vcfgo/tmp/dbpedia_sample.jsonl
./bin/rag import --file-input $path_input test
```

## chat
```shell
./bin/rag chat test "When did the Monarch Company exist?"
```
