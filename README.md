# go-hello-dependency-injection

## Usage

A simple program to show dependency injection.

### Invocation

```console
go-hello-dependency-injection
```

## Development

### Dependencies

#### Set environment variables

```console
export GOPATH="${HOME}/go"
export PATH="${PATH}:${GOPATH}/bin:/usr/local/go/bin"
export PROJECT_DIR="${GOPATH}/src/github.com/docktermj"
export REPOSITORY_DIR="${PROJECT_DIR}/go-hello-dependency-injection"
```

#### Download project

```console
mkdir -p ${PROJECT_DIR}
cd ${PROJECT_DIR}
git clone git@github.com:docktermj/go-hello-dependency-injection.git
```

#### Download dependencies

```console
cd ${REPOSITORY_DIR}
make dependencies
```

### Build

#### Local build

```console
cd ${REPOSITORY_DIR}
make
```

The results will be in the `${GOPATH}/bin` directory.

#### Docker build

Create `rpm` and `deb` installation packages.

```console
cd ${REPOSITORY_DIR}
make build
```

The results will be in the `${REPOSITORY_DIR}/target` directory.

### Test

```console
cd ${REPOSITORY_DIR}
make test-local
```

### Cleanup

```console
cd ${REPOSITORY_DIR}
make clean
```
