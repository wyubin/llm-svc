package cmd

import (
	"bytes"
	_ "embed"
	"log/slog"
	"os"

	"github.com/wyubin/llm-svc/utils/log"
	"github.com/wyubin/llm-svc/utils/viperkit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	KEY_EMBEDDING = "context"
)

//go:embed env_default
var byteEnv []byte

var (
	modelInference string
	ckDebug        bool

	uriDB      string
	uriOpenAI  string
	pathDBinfo string

	ctlCmd = &cobra.Command{
		Use:           "rag",
		Short:         "rag – prepare and chat with modelkb",
		Long:          ``,
		Version:       "0.1.0",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func Execute() error {
	defer log.LogExeTime("rag")()
	return ctlCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	viperkit.ReaderEnv(bytes.NewReader(byteEnv))
	viper.AutomaticEnv()
	ctlCmd.PersistentFlags().BoolVar(&ckDebug, "debug", false, "show debug log if needed")

	ctlCmd.PersistentFlags().StringVar(&uriDB, "uri-db", viper.GetString("QDRANT_URI"), "assign uri for vector store")
	ctlCmd.PersistentFlags().StringVar(&uriOpenAI, "uri-openai", viper.GetString("OPENAI_URI"), "assign uri for llm chat api")
	ctlCmd.PersistentFlags().StringVar(&pathDBinfo, "path-db-info", viper.GetString("PATH_DB_INFO"), "config path for db info")
	ctlCmd.PersistentFlags().StringVar(&modelInference, "model-inference", viper.GetString("MODEL_INFERENCE"), "assign model for inference")

	ctlCmd.AddCommand(createCmd)
	ctlCmd.AddCommand(importCmd)
	ctlCmd.AddCommand(chatCmd)
}

func initConfig() {
	// init logger
	var logLevel slog.Level = slog.LevelInfo
	if ckDebug {
		logLevel = slog.LevelDebug
	}
	log.InitLogger(logLevel, os.Stderr)
}
