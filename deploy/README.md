# Deploying gpuctl


## Building Docker images freestanding

Due to how Docker context works, you need to do this from the top level directory.

```console
alona@Ashtabula:~/dev/gpuctl$ podman build -f ./deploy/control.Dockerfile .
alona@Ashtabula:~/dev/gpuctl$ podman build -f ./deploy/frontend.Dockerfile .
```

## npm ERR! EMFILE: too many open files

You may see an error like:

```
[1/3] STEP 4/6: RUN npm install
npm ERR! code EMFILE
npm ERR! syscall open
npm ERR! path /root/.npm/_cacache/index-v5/7d/8e/9676576fe239de89dec5e769bcbad0def29c3e4fd33b2caebf5c78716e7e
npm ERR! errno -24
npm ERR! EMFILE: too many open files, open '/root/.npm/_cacache/index-v5/7d/8e/9676576fe239de89dec5e769bcbad0def29c3e4fd33b2caebf5c78716e7e'

npm ERR! A complete log of this run can be found in: /root/.npm/_logs/2024-02-08T11_54_42_543Z-debug-0.log
Error: building at STEP "RUN npm install": while running runtime: exit status 232
```

This can be solved with the `--ulimit` flag:

```console
alona@Ashtabula:~/dev/gpuctl$ podman build --ulimit=4096:4096 -f ./deploy/frontend.Dockerfile .
```