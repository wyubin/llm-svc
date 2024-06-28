package cmd

import (
	"fmt"
	"os"

	"github.com/wyubin/llm-svc/cmd/rag/ctl"
	"github.com/wyubin/llm-svc/utils/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	modelEmbedding string
	createProc     CreateProc = CreateProc{}
	createCmd                 = &cobra.Command{
		Use:   "create [flags] nameColl [...nameColl]",
		Short: "create - collection for rag",
		Long:  ``,
		Run:   createProc.run,
	}
)

type CreateProc struct {
	ctlDB  *ctl.DB
	ctlCfg *ctl.Conf
	ctlGpt *ctl.Gpt
}

// Processes a VCF file using the defined create protocols.
func (s *CreateProc) run(ccmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Logger.Warn("need to input at least one collection name")
		os.Exit(1)
	}
	// init ctl
	err := s.Init()
	if err != nil {
		log.Logger.Error(fmt.Sprintf("init failed: %s", err.Error()))
		os.Exit(1)
	}
	defer s.ctlDB.Close()
	for _, nameColl := range args {
		// check nameColl exists in db, if has, exit with 1
		if s.ctlDB.HasColl(nameColl) {
			log.Logger.Warn(fmt.Sprintf("collection %s already exists in db", nameColl))
			os.Exit(1)
		}
		// check info of collections if not exist, use from args and update collections
		coll, err := s.ctlCfg.GetCollInfo(nameColl)
		if err != nil {
			coll.Name = nameColl
			coll.Model = modelEmbedding
		}
		// check modelEmbedding exists in uri-openai, if not, exit with 1
		hasModel, err := s.ctlGpt.HasModel(coll.Model)
		if err != nil {
			log.Logger.Error(err.Error())
			os.Exit(1)
		}
		if !hasModel {
			log.Logger.Warn(fmt.Sprintf("modelEmbedding %s not found in uri-openai", coll.Model))
			os.Exit(1)
		}
		// get vector-length of modelEmbedding and update info if models[info] not exists
		embedding, err := s.ctlGpt.Embedding("test", coll.Model)
		if err != nil {
			log.Logger.Error(fmt.Errorf("get embedding failed: %w", err).Error())
			os.Exit(1)
		}
		modelInfo := ctl.Model{
			Name:      coll.Model,
			VectorLen: len(embedding),
		}
		s.ctlCfg.UpdateModel(modelInfo)
		// create collection in db
		err = s.ctlDB.CreateColl(coll.Name, modelInfo.VectorLen)
		if err != nil {
			log.Logger.Error(fmt.Errorf("create collection failed: %w", err).Error())
			os.Exit(1)
		}
		// update info with coll info and model info
		err = s.ctlCfg.UpdateColl(coll)
		if err != nil {
			log.Logger.Error(fmt.Errorf("update info failed: %w", err).Error())
			os.Exit(1)
		}
		log.Logger.Info(fmt.Sprintf("collection %s created", coll.Name))
	}
}

// init by flags
func (s *CreateProc) Init() error {
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
	persistFlag := createCmd.PersistentFlags()
	persistFlag.StringVar(&modelEmbedding, "model-embedding", viper.GetString("MODEL_EMBEDDING"), "assign model embedding")
}
