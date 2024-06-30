FROM cosmtrek/air as run
WORKDIR /app
ENV air_wd=/app

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY .air.toml .

COPY . .
