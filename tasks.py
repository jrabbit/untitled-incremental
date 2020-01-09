from invoke import task
from fabric import Connection

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
def deploy(c):
    conn = Connection("trotsky")
    with conn.cd("~/src/untitled-incremental"):
        conn.run("git fetch")
        conn.run("git stash")
        conn.run("git pull")
        conn.run("git stash pop")
        conn.put("main.wasm")
        
@task
def tinygo(c, dev=False):
    "installed via debian package"
    # https://github.com/tinygo-org/tinygo
    # https://tinygo.org/getting-started/linux/
    c.run("/usr/local/tinygo/bin/tinygo build -o wasm.wasm -target=wasm main.go")