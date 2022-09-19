# Cryptex
Cryptex is a CLI application to manage and in a fast and smart way your secrets. It is native app for Mac OS, Linux and Windows.

## Compiling your protocol buffers

Now that you have a .proto, the next thing you need to do is generate the classes:

1. If you haven't installed the compiler, [download](https://developers.google.com/protocol-buffers/docs/downloads) the package and follow the instructions in the README.

2. Run the following command to install the Go protocol buffers plugin:
```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
3. Now run the compiler, specifying the source directory (where your application's source code lives – the current directory is used if you don't provide a value), the destination directory (where you want the generated code to go; often the same as $SRC_DIR), and the path to your .proto. In this case, you would invoke:
```shell
protoc -I=. --go_out=. protos/box.proto
```


