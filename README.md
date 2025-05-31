`arxiv-researcher-go`
=====================

This RAG-based agent searches for research papers on a given topic of interest.
The agent takes a short phrase describing the topic
and searches for relevant papers in its own private knowledge database.
If it finds no papers in its private database,
it expands its search to arXiv.
After completing its search,
the agent will display a list of any relevant papers it found
and download the papers to the local file system.

# Execution environment setup

To run the chatbot,
you need to have accounts with OpenAI and Pinecone and corresponding API keys.
When you are ready, follow these instructions:

1. Download and install the [Go runtime and development environment]((https://go.dev/doc/install)).
1. Get an API key to access [OpenAI LLMs](https://help.openai.com/en/articles/4936850-where-do-i-find-my-openai-api-key).
1. Get an API key to read from and write to a [Pinecone vector store](https://docs.pinecone.io/guides/projects/manage-api-keys).
1. Create a new index named `arxiv-researcher-playground` and configure it to use `text-embedding-3-small` with the default dimension of 1536.
1. Copy the `.env.example` file in the local repository root to a new `.env` file, and fill in these environment variables:
   1. `OPENAI_API_KEY`
   1. `PINECONE_API_KEY`
   1. `PINECONE_HOST_NAME`

# Private knowledge database population

To demonstate the utility of a RAG-based agent,
we first need to populate a private knowledge database with papers from arXiv on a specifc topic of interest.
To populate the knowledge database, run
```
$ go run cmd/indexer/main.go <topic phrase>
```
in the root of the cloned repository,
where `<topic phrase>` is a query phrase describing the topic.
The indexer takes the query phrase,
searches arXiv for relevant papers,
and saves metadata for the papers (including abstracts) in its database.

# Paper search

To search for papers in the knowledge base on some general topic of interest, run
```
$ go run cmd/agent/main.go <topic phrase>
```
in the root of the cloned repository,
where `<topic phrase\>` is a query phrase describing the topic.
The agent takes the query phrase,
and searches its knowledge database for relevant papers.
If the agent finds no relevant papers,
it expands its search to arXiv.
After completing its search,
the agent will display a list of any relevant papers it found
and download the papers to the local file system.

# Acknowledgements

I built this chatbot after completing the Udemy course
"[Build an AI Agent (OpenAI, LlamaIndex, Pinecone & Streamlit)](https://www.udemy.com/course/build-an-ai-agent-openai-llamaindex-pinecone-streamlit/?couponCode=CP130525US)"
hosted by [David Armendáriz](https://www.udemy.com/user/david-4271/).
David's course guides you through the design and implementation of an arXiv search chatbot in Python that uses OpenAI and Pinecone for its LLM and vector store backends,
and [LlamaIndex](https://pypi.org/project/llama-index-core/) as the agent framework. For my project, I decided ti reimplement the chatbot in Go to teach myself the language,
but instead used [LangChainGo](https://tmc.github.io/langchaingo/docs/) as the agent framework
as the Go port of LlamaIndex is somewhat incomplete.
