FROM asia-southeast1-docker.pkg.dev/tk-dev-micro/base-image/golang:1.22-bullseye as builder

ARG GITHUB_TOKEN

WORKDIR /src/

RUN mkdir ~/.ssh
RUN touch ~/.ssh/known_hosts
RUN ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

COPY . .

RUN go test ./...
RUN CGO_ENABLED=0 go build -buildvcs=false -ldflags "-s -w" -o {{SERVICE_NAME}} ./cmd/{{PROVIDER}}

# create final docker image
FROM asia-southeast1-docker.pkg.dev/tk-dev-micro/base-image/distroless-go

WORKDIR /app

COPY --from=builder /src/{{SERVICE_NAME}} .

ENTRYPOINT ["/app/{{SERVICE_NAME}}"]