# STAGE BUILD IMAGE
FROM golang:1.23-alpine AS build_api
WORKDIR /ocr-api
COPY /ocr-api/go.mod /ocr-api/go.sum ./
RUN go mod download
COPY /ocr-api ./
RUN go build -o main /ocr-api/cmd/api/main.go

# DEVELOPMENT
FROM alpine:3.20.1 AS development_api
WORKDIR /ocr-api
COPY --from=build_api /ocr-api/main /ocr-api/main
COPY --from=build_api /ocr-api/.air.toml /ocr-api/.air.toml
RUN apk add --no-cache curl poppler-utils \
    && curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b /usr/local/bin
COPY --from=build_api /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"
EXPOSE ${PORT_OCR_DEV}
CMD ["air", "-c", ".air.toml"]

# PRODUCTION
FROM alpine:3.20.1 AS production_api
WORKDIR /ocr-api
COPY --from=build_api /ocr-api/main /ocr-api/main
RUN apk add --no-cache poppler-utils
EXPOSE ${POST_OCR}
CMD ["./main"]



FROM python:3.10-slim AS build_manage
WORKDIR /manage
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1 
RUN pip install --upgrade pip 
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
       make build-essential libssl-dev zlib1g-dev \
       libbz2-dev libreadline-dev libsqlite3-dev wget curl llvm \
       libncurses5-dev libncursesw5-dev xz-utils tk-dev libffi-dev \
       liblzma-dev python3-openssl git python3-dev \
       default-libmysqlclient-dev pkg-config \
       libmariadb-dev-compat libmariadb-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
COPY /manage/requirements.txt /manage/
RUN pip install --no-cache-dir -r requirements.txt

# PRODUCTION
FROM python:3.10-slim AS manage
WORKDIR /manage
COPY --from=build_manage /usr/local/lib/python3.10/site-packages/ /usr/local/lib/python3.10/site-packages/
COPY --from=build_manage /usr/local/bin/ /usr/local/bin/
COPY --from=build_manage /usr/lib/x86_64-linux-gnu/libmariadb.so.3 /usr/lib/x86_64-linux-gnu/libmariadb.so.3
COPY /manage/ .
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1 
EXPOSE ${PORT_MANAGE}
