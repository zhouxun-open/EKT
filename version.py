import os

main = "cmd/enode/main.go"

def get_version():
    for line in open(main):
        line = line.strip()
        if not line.startswith("version"):
            continue

        items = map(lambda x:x.strip(), line.split("="))
        if len(items) != 2 or items[0] != "version":
            continue

        return items[1][1:-1]

print get_version()
