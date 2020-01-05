from invoke import task

@task
def build(c):
    c.run("go build -o main.wasm", env={"GOOS": "js", "GOARCH": "wasm"})

@task
def vet(c):
    c.run("go vet", env={"GOOS": "js", "GOARCH": "wasm"})

@task
def auto(c):
    c.run("ag -l | entr -rc inv build", pty=True)


@task
def tinygo(c, dev=False):
    "installed via debian package"
    # https://github.com/tinygo-org/tinygo
    c.run("/usr/local/tinygo/bin/tinygo build -o wasm.wasm -target=wasm main.go")