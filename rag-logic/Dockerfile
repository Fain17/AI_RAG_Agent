FROM python:3.12-slim-bookworm

# Set working directory
WORKDIR /app

# Install system dependencies (for pdfminer and model ops)
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    libffi-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements
COPY requirements.txt .

# Install Python packages
RUN pip install --no-cache-dir -r requirements.txt

# Copy the app code
COPY app/ ./app/

# Expose port
EXPOSE 8000

# Run the embedder FastAPI app
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
