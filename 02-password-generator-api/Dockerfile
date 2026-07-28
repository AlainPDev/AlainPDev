# Use a lightweight official Python image.
FROM python:3.12-slim

# Prevent Python from writing bytecode and buffer-free application logs.
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

# Store the application inside /app in the container.
WORKDIR /app

# Install dependencies before copying source code to improve build caching.
COPY requirements.txt .

RUN python -m pip install \
    --no-cache-dir \
    --upgrade pip \
    && python -m pip install \
    --no-cache-dir \
    -r requirements.txt

# Copy the application package into the container.
COPY app ./app

# Document the port used by the FastAPI server.
EXPOSE 8000

# Start the production server.
CMD ["fastapi","run","app/main.py","--host","0.0.0.0","--port","8000"]
