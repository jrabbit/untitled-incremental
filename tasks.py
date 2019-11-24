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
