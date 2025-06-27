from fastapi import FastAPI
from app.routes import routes

app = FastAPI(title="Business Logic API", version="1.0")
app.include_router(routes.router)
