# Learn Go

## Running the application

### Todo api

```
go run ./cmd/todo-api/main.go
```

## CGO

CGO is required for the sqlite connection.

To enable CGO set the environment variable

```
CGO_ENABLED=1
```

### Windows

Install the `gcc` toolchain using MSYS2 ([installer link](https://github.com/msys2/msys2-installer/releases/download/2024-01-13/msys2-x86_64-20240113.exe))

After installing open a MSYS2 terminal (Run MSYS2 now checkbox in installer) in the opened terminal run:

```
pacman -S --needed base-devel mingw-w64-ucrt-x86_64-toolchain

```

When asked to make a selection use all (default, pres `enter`)

When prompted to confirm the installation pres `y`

After installing add the MSYS2 bin path to the system or user path. (default install location = `C:\msys64\ucrt64\bin`)
