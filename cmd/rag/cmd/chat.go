package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wyubin/llm-svc/cmd/rag/ctl"
	"github.com/wyubin/llm-svc/utils/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	chatProc ChatProc = ChatProc{}
	chatCmd           = &cobra.Command{
		Use:   "chat [flags] nameColl question",
		Short: "chat - ask questions and get responses with vector store",
		Long:  ``,
		Run:   chatProc.run,
	}
)

type ChatProc struct {
	ctlDB  *ctl.DB
	ctlCfg *ctl.Conf
	ctlGpt *ctl.Gpt
}

// Processes a VCF file using the defined chat protocols.
func (s *ChatProc) run(ccmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.Logger.Warn("need to input one collection name, and question")
		os.Exit(1)
	}
	nameColl, question := args[0], args[1]
	// init ctl
	err := s.Init()
	if err != nil {
		log.Logger.Error(fmt.Sprintf("init failed: %s", err.Error()))
		os.Exit(1)
	}
	defer s.ctlDB.Close()
	// check nameColl exists in info, if not, exit with 1
	coll, err := s.ctlCfg.GetCollInfo(nameColl)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("get collection info failed: %s", err.Error()))
		os.Exit(1)
	}
	// make embedding of question
	embedding, err := s.ctlGpt.Embedding(question, coll.Model)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("make embedding failed: %s", err.Error()))
		os.Exit(1)
	}
	// retrieve vector store
	log.Logger.Info("retrieve vector store...")
	s.ctlDB.ChooseColl(coll.Name)
	contexts, err := s.ctlDB.VectorSearch(embedding, 3)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("retrieve vector store failed: %s", err.Error()))
		os.Exit(1)
	}
	// chat with rag
	s.ctlGpt.SetDirTmpl(filepath.Join(viper.GetString("PATH_TMPL"), "chat"))
	resp, err := s.ctlGpt.ChatWithContexts(modelInference, question, contexts)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("chat failed: %s", err.Error()))
		os.Exit(1)
	}
	fmt.Printf("Response from %s: %s\n", modelInference, resp)
}

// init by flags
func (s *ChatProc) Init() error {
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
}
