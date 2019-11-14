from invoke import task

@task
def build(c):
    c.run("go build -o main.wasm", env={"GOOS": "js", "GOARCH": "wasm"})