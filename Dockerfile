FROM python:3.12

# RUN apt-get -y update; apt-get debugpy

WORKDIR /app

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt

EXPOSE 5678

CMD [ "python", "backend/src/index.py"]
