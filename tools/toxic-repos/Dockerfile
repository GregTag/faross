FROM python:latest
COPY toxic-repos-checker.py /app/toxic-repos-checker.py
WORKDIR /app
RUN pip install --no-cache-dir requests
ENTRYPOINT ["python3", "toxic-repos-checker.py"]
