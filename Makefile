INFO = `git describe --tags --dirty --broken --abbrev=40`
PRE = CGO_ENABLED=0 GOOS=linux GOARCH=amd64
BINS = ./bin
API = ${BINS}/api
SEEDER = ${BINS}/seeder
BUILD = go build -o

all: ${API} ${SEEDER}

go.sum: go.mod
	go mod tidy
	go mod download

.PHONY: api
api: ${API}
${API}: go.sum
	${PRE} ${BUILD} ${API} ./cmd/api/main.go

.PHONY: seeder
seeder: ${SEEDER}
${SEEDER}: go.sum
	${PRE} ${BUILD} ${SEEDER} ./cmd/seeder/main.go

.PHONY: clean
clean:
	go clean
	rm ${API} ${SEEDER}

