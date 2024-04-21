FROM python:3.9
COPY run.py /app/run.py
WORKDIR /app
RUN pip install --no-cache-dir requests urllib3==1.26.6
ENTRYPOINT ["python", "run.py"]