package constants

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
)

var Llm *openai.LLM // The LLM used to generate responses to user queries.

var llmEmbedding = "text-embedding-3-small"
var llmModel = "gpt-4o-mini"

func init() {
	// Initialize constants for the tools package. In particular, initialize the OpenAI LLM model.
	// Unlike Llama, LangChainGo obtains the OpenAI API key implicitly.
	var err error
	err = godotenv.Load("/Users/tmwong/code/arxiv-researcher-go/.env")
	if err != nil {
		log.Fatalln("Failed loading .env file:", err)
	}
	Llm, err = openai.New([]openai.Option{
		openai.WithEmbeddingModel(llmEmbedding),
		openai.WithModel(llmModel),
	}...)
	if err != nil {
		log.Fatalln("Failed initializing LLM:", err)
	}
}
