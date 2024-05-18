FROM python:3.9
COPY . /app/
WORKDIR /app
RUN pip install --no-cache-dir requests urllib3==1.26.6
ENTRYPOINT ["python3", "toxic-repos-checker.py"]