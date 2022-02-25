# Running this exmaple project on WIndows 10 WSL and VS Code
## Installation
### WSL 2 on WIndows 10
1. [Install WSL as a feature of your Windows 10.](https://docs.microsoft.com/en-us/windows/wsl/install)
    * [You will need to make sure WSL is enabled and on version 2 for your windows edition.](https://docs.microsoft.com/en-us/windows/wsl/install#upgrade-version-from-wsl-1-to-wsl-2)
1. [You will need to make sure you have VSCode installed.](https://code.visualstudio.com/download)
1. [Install the extension to enbable VSCode to work with WSL.](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-wsl)
1. Clone this repo.

### Troublshooting internet access on WSL
**["Temporary failure resolving {url}"](https://stackoverflow.com/questions/55649015/could-not-resolve-host-github-com-only-in-windows-bash)**  
Instruction copied here for postarity. On the instance;  
`sudo bash -c 'echo "[network]\ngenerateResolvConf = false" > /etc/wsl.conf && rm /etc/resolv.conf && echo "nameserver 1.1.1.1" > /etc/resolv.conf'`  

### Installing Docker on WSL, protoc engine and GO
Run these commands;
```
./dev/docker-setup-on-wsl.sh
./dev/install-go.sh
./dev/install-protoc.sh
./dev/build-service.sh
```

### Testing
Tests are done with Testify.  They are located in the `/test` folder. You can run them manually by using the same command the docker file uses; `go test -timeout 2s -v ./test/...` (or simply `go test ./...`)

Also in the test folder is the test client.  This quick go program proves the client created by the protoc can talk to the server created by the protoc. While not a robust test suite, it demonstrates that the client can be tested from the server code if you need to. You can compile it using a provided script (`dev\test-client-compile.sh`).

### Building
The project can be built youself (`go build -o ./bin/ipauthorize ./cmd/ipauthorize`), though it is intended to be run as a container.

You can create the docker image with;
```
docker build . -t dockersample
```

### Troublshooting "Docker has no internet access"
In WSL2, I had to use the host network because the Docker engine has no access to the internet.
```
docker build --network host . -t dockersample
docker run --network host -p 9080:9080 -p 9079:9079 $(docker image ls dockersample -q)
```

### Running
```
docker run -p 9080:9080 -p 9079:0079 $(docker image ls dockersample -q)
```
You can now use the sample call script; `dev/samplecall.sh`  
--or--  
You can use the compiled test client from earlier; `bin/testclient`


