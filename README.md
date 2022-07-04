# FileServer

A simple file-server in golang, supports password & fake password.

## Install

```bash
go install github.com/flyfy1/fileserver@latest
```

## Usage

```bash
fileserver -h
# Usage of fileserver:
#   -d string
#         the directory of static file to host. Default to current dir (default ".")
#   -fake string
#         if provided, serve a fake webpage instead
#   -p string
#         port to serve on (default "8100")
#   -secret string
#         secret to browse file. If provided, would check this secret cookie

fileserver -secret true-secret -fake false-secret
# 2022/07/04 08:38:35 Serving . on HTTP port: 8100
```
