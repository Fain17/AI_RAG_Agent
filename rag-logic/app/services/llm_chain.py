from langchain.agents import AgentType, Tool, initialize_agent
from langchain.prompts import PromptTemplate
from langchain_ollama import OllamaLLM

# 1. Base LLM setup with Mistral
llm = OllamaLLM(
    model="mistral", base_url="http://192.168.5.192:11434", temperature=0.2
)

# 2. Prompt for QA
qa_template = """You are a highly intelligent assistant with access to \
relevant context extracted from documents and database schemas.

Your job is to answer questions accurately using the provided information. \
When answering, follow these rules:

1. Use ONLY the information provided in the context below
2. If the context doesn't contain enough information to answer the question, \
say "I don't have enough information to answer this question"
3. Be specific and cite relevant parts of the context when possible
4. If asked about database operations, provide SQL examples when appropriate
5. Keep your answers concise but complete

Context:
{context}

Question: {question}

Answer:"""


prompt = PromptTemplate.from_template(qa_template)
# qa_chain = LLMChain(llm=llm, prompt=prompt)
chain = prompt | llm


# 3. Optional tool
def dummy_tool(input: str) -> str:
    return f"[DummyTool] You asked about: {input}"


tools = [
    Tool(
        name="DummyTool", func=dummy_tool, description="Basic string echo tool"
    )
]

# 4. Agent setup
agent = initialize_agent(
    tools=tools,
    llm=llm,
    agent=AgentType.ZERO_SHOT_REACT_DESCRIPTION,
    verbose=True,
)
