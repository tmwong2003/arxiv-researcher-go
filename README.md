`arxiv-researcher-go`
=====================

This RAG-based chatbot searches for research papers related to user-supplied topics.
The agent takes a short topic keyword phrase as input and searches for relevant papers in its own private Pinecone knowledge database.
If it finds no papers in its database, it expands its search to arXiv.
After completing its search, the agent will display a list of any relevant papers it found
and download the papers to the local file system.

# Execution environment setup

To run the chatbot, you need have accounts with OpenAI and Pinecone. When you are ready, follow these instructions:

1. Download and install the [Go runtime and development environment]((https://go.dev/doc/install)).
2. Get an API key to access [OpenAI LLMs](https://help.openai.com/en/articles/4936850-where-do-i-find-my-openai-api-key).
3. Get an API key to read from and write to a [Pinecone vector store](https://docs.pinecone.io/guides/projects/manage-api-keys).
4. Create a new index named `arxiv-researcher-playground` and configure it to use `text-embedding-3-small` with the default dimension of 1536.
5. Copy the `.env.example` file in the local repository root to a new `.env` file, and fill in these environment variables:
   1. `OPENAI_API_KEY`
   2. `PINECONE_API_KEY`
   3. `PINECONE_HOST_NAME`

# Private knowledge database population

To populate the "private" knowledge database to index papers for a specific topic, run the following command:
```
$ go run cmd/indexer/main.go <topic keyword>
```
in the root of the cloned repository, where \<topic keyword\> is a short phrase describing the topic of interest.
The indexer will search arXiv for papers relevant to the topic, and save information about the papers in the Pinecone vector store.

# Paper search

To search for research papers related to a given topic, run the following command:
```
$ go run cmd/agent/main.go <topic keyword>
```
in the root of the cloned repository, where \<topic keyword\> is a short phrase describing the topic of interest.
The agent takes a short topic keyword phrase as input and searches for relevant papers in its knowledge database.
If it finds no papers in its database, it expands its search to arXiv.
After completing its search, the agent will display a list of any relevant papers it found
and download the papers to the local file system.

# Acknowledgements

I built this chatbot after completing the Udemy course
"[Build an AI Agent (OpenAI, LlamaIndex, Pinecone & Streamlit)](https://www.udemy.com/course/build-an-ai-agent-openai-llamaindex-pinecone-streamlit/?couponCode=CP130525US)"
hosted by [David Armendáriz](https://www.udemy.com/user/david-4271/).
David's course guides you through the design and implementation of an arXiv search chatbot in Python that uses OpenAI and Pinecone for its LLM and vector store backends,
and [LlamaIndex](https://pypi.org/project/llama-index-core/) as the agent framework. For my project, I decided ti reimplement the chatbot in Go to teach myself the language,
but instead used [LangChainGo](https://tmc.github.io/langchaingo/docs/) as the agent framework
as the Go port of LlamaIndex is somewhat incomplete.
