package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/wyubin/llm-svc/cmd/rag/ctl"
	"github.com/wyubin/llm-svc/utils/customflag"
	"github.com/wyubin/llm-svc/utils/log"
	"github.com/spf13/cobra"
)

var (
	pathInput  customflag.FlagPath = ""
	keyEmbed   string
	importProc ImportProc = ImportProc{}
	importCmd             = &cobra.Command{
		Use:   "import [flags] nameColl",
		Short: "import - embeded and import data into vector store",
		Long:  ``,
		Run:   importProc.run,
	}
)

type ImportProc struct {
	ctlDB  *ctl.DB
	ctlCfg *ctl.Conf
	ctlGpt *ctl.Gpt
}

// Processes a VCF file using the defined import protocols.
func (s *ImportProc) run(ccmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Logger.Warn("need to input one collection name")
		os.Exit(1)
	}
	// init ctl
	err := s.Init()
	if err != nil {
		log.Logger.Error(fmt.Sprintf("init failed: %s", err.Error()))
		os.Exit(1)
	}
	defer s.ctlDB.Close()
	// check nameColl exists in db, if has, exit with 1
	nameColl := args[0]
	if !s.ctlDB.HasColl(nameColl) {
		log.Logger.Warn(fmt.Sprintf("collection %s not found", nameColl))
		os.Exit(1)
	}
	// check info or use modelEmbedding
	coll, err := s.ctlCfg.GetCollInfo(nameColl)
	if err != nil {
		coll.Name = nameColl
		coll.Model = modelEmbedding
		err := s.ctlCfg.UpdateColl(coll)
		if err != nil {
			log.Logger.Error(fmt.Sprintf("update collection failed: %s", err.Error()))
			os.Exit(1)
		}
	}
	// file load
	fileInput, err := os.Open(pathInput.String())
	if err != nil {
		log.Logger.Error(fmt.Sprintf("open file failed: %s", err.Error()))
		os.Exit(1)
	}
	defer fileInput.Close()
	decoder := json.NewDecoder(fileInput)
	err = s.ctlDB.ChooseColl(nameColl)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("choose collection failed: %s", err.Error()))
		os.Exit(1)
	}
	log.Logger.Info("start import data...")
	itemNum, countData := 0, 0
	for {
		dataTmp := map[string]string{}
		err := decoder.Decode(&dataTmp)
		if err != nil {
			break
		}
		itemNum++
		// check keyEmbed exists in dataTmp, if not, warn and continue
		srcEmbed, found := dataTmp[keyEmbed]
		if !found {
			log.Logger.Warn(fmt.Sprintf("keyEmbed %s not found in data line: %d", keyEmbed, itemNum))
			continue
		}
		delete(dataTmp, keyEmbed)
		dataTmp[KEY_EMBEDDING] = srcEmbed
		// generate embedding
		embedding, err := s.ctlGpt.Embedding(srcEmbed, coll.Model)
		if err != nil {
			log.Logger.Error(fmt.Sprintf("generate embedding failed: %s", err.Error()))
			os.Exit(1)
		}
		// insert data
		err = s.ctlDB.InsertPoint(embedding, dataTmp)
		if err != nil {
			log.Logger.Error(fmt.Sprintf("insert data failed: %s", err.Error()))
			os.Exit(1)
		}
		countData++
	}
	log.Logger.Info(fmt.Sprintf("imported %d data into collection %s", countData, nameColl))
}

// init by flags
func (s *ImportProc) Init() error {
	uriDB := fmt.Sprintf("qdrant:%s", uriDB)
	db, err := ctl.NewDB(uriDB)
	if err != nil {
		return err
	}
	s.ctlDB = db
	cfg, err := ctl.NewConf(pathDBinfo)
	if err != nil {
		return err
	}
	s.ctlCfg = cfg
	// init ollama with modelEmbedding
	uriOpenAI := fmt.Sprintf("http://%s", uriOpenAI)
	api, err := ctl.NewGpt(uriOpenAI)
	if err != nil {
		return err
	}
	s.ctlGpt = api
	return nil
}

func init() {
	persistFlag := importCmd.PersistentFlags()
	persistFlag.Var(&pathInput, "file-input", "jsonl file to import data")
	importCmd.MarkPersistentFlagRequired("file-input")
	persistFlag.StringVar(&keyEmbed, "key-embed", "text", "key for generate embedding vectors")
}
