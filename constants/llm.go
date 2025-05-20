package constants

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
)

// The singleton LLM instance used to generate responses to user/agent queries.
var Llm *openai.LLM

var llmEmbedding = "text-embedding-3-small"
var llmModel = "gpt-4o-mini"

// Initialize constants for the tools package.
// In particular, initialize the OpenAI LLM model.
// Unlike Llama, LangChainGo obtains the OpenAI API key implicitly from the O/S environment.
func init() {
	var err error
	err = godotenv.Load("/Users/tmwong/code/arxiv-researcher-go/.env")
	if err != nil {
		log.Fatalln("failed while loading .env file: ", err)
	}
	Llm, err = openai.New([]openai.Option{
		openai.WithEmbeddingModel(llmEmbedding),
		openai.WithModel(llmModel),
	}...)
	if err != nil {
		log.Fatalln("failed while initializing LLM: ", err)
	}
}
