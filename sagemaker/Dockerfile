FROM python:3.7-buster

WORKDIR /apps

# setup the serving
RUN apt-get -y update && apt-get install -y --no-install-recommends nginx git

# deps which algoflow depends on
COPY requirements-dep.txt .
RUN pip install --no-cache-dir -r requirements-dep.txt

COPY requirements.txt .
RUN --mount=type=ssh mkdir -p ~/.ssh/ && \
    ssh-keyscan github.com >> ~/.ssh/known_hosts && \
    pip install --no-cache-dir -r requirements.txt


# Keeps Python from generating .pyc files in the container
ENV PYTHONDONTWRITEBYTECODE=1
# Turns off buffering for easier container logging
ENV PYTHONUNBUFFERED=1
ENV PATH="/opt/program:${PATH}"

# Set up the program in the image
COPY container /opt/program

WORKDIR /opt/program

RUN chmod +x ./train && chmod +x ./serve


