FROM python:3.12

# RUN apt-get -y update; apt-get debugpy

WORKDIR /app

COPY requirements.txt .

COPY db db
COPY backend backend
COPY frontend frontend
COPY Makefile Makefile

RUN pip install --no-cache-dir -r requirements.txt

VOLUME /app

EXPOSE 5678

CMD [ "python", "-m", "debugpy", "--listen", "0.0.0.0:5678", "--wait-for-client", "backend/src/index.py"]
