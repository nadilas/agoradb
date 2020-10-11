* [Client generation](#client-generation)


#### Client generation 

##### Compiler: protoc or flatc

###### FlatBuffers
Take a look at Google's guide: [Building the compiler](https://google.github.io/flatbuffers/flatbuffers_guide_building.html)

###### Protocol Buffers

Here are some OS-specific options for installing the binary. These instructions also install basic .proto files like wrappers.proto, any.proto and descriptor.proto. (Those files aren’t needed by proto-lens itself, but they may be useful for other language bindings/plugins.)

####### Mac OS X
If you have Homebrew (which you can get from [https://brew.sh](https://brew.sh)), just run:

```
# depends on protobuf (protoc)
brew install protoc-gen-go
```
If you see any error messages, run brew doctor, follow any recommended fixes, and try again. If it still fails, try instead:

brew upgrade protobuf
Alternately, run the following commands:

```bash
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-osx-x86_64.zip
sudo unzip -o protoc-3.7.1-osx-x86_64.zip -d /usr/local bin/protoc
sudo unzip -o protoc-3.7.1-osx-x86_64.zip -d /usr/local 'include/*'
rm -f protoc-3.7.1-osx-x86_64.zip
```

####### Linux
Run the following commands:

```bash
curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip
sudo unzip -o protoc-3.7.1-linux-x86_64.zip -d /usr/local bin/protoc
sudo unzip -o protoc-3.7.1-linux-x86_64.zip -d /usr/local 'include/*'
rm -f protoc-3.7.1-linux-x86_64.zip
```
Alternately, manually download and install protoc from [here](https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip).
