from langchain.prompts import PromptTemplate
from langchain.chains import LLMChain
from langchain_ollama import OllamaLLM
from langchain.agents import initialize_agent, Tool, AgentType

# 1. Base LLM setup with Mistral
llm = OllamaLLM(
    model="mistral",
    base_url="http://192.168.5.192:11434",
    temperature=0.2
    )

# 2. Prompt for QA
qa_template = """You are a helpful assistant with access to code context.

Context:
{context}

Question:
{question}

Answer:"""

prompt = PromptTemplate.from_template(qa_template)
#qa_chain = LLMChain(llm=llm, prompt=prompt)
chain = prompt | llm


# 3. Optional tool
def dummy_tool(input: str) -> str:
    return f"[DummyTool] You asked about: {input}"

tools = [Tool(name="DummyTool", func=dummy_tool, description="Basic string echo tool")]

# 4. Agent setup
agent = initialize_agent(
    tools=tools,
    llm=llm,
    agent=AgentType.ZERO_SHOT_REACT_DESCRIPTION,
    verbose=True
)