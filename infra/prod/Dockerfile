FROM python:3.12

WORKDIR /app

COPY requirements.txt .

COPY db db
COPY backend backend
COPY frontend frontend
COPY Makefile Makefile

RUN pip install --no-cache-dir -r requirements.txt

VOLUME /app

EXPOSE 5678

CMD [ "python", "backend/src/index.py"]
