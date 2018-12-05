progs are tiny utilities written in Go.

but why?

 * readablity (no more `if [ $# -eq 42]; then echo $1; fi;` 
 * static binary (contains everything, no need to bring any userspace stuff)
 * painless crosscompilation (just set GOOS, GOARCH and run `go build`)
 * extensibility (there are tons of libraries one can use)
 * ???
 * the gopher!!

| prog | description |
| ---- | ----------- |
| `bundle` | builds Go application, puts it into Docker image, and pushes it to the registry |
| `chrome-tmp` | starts chrome in Incognito mode with ephemeral profile |
| `dockmach` | run `docker` via docker-machine without configuring envrironment for each terminal |
| `frontprobe` | quickly probe domain fronting availability |
| `getsignal` | fetch latest Signal APK |
| `goplay` | run a local .go file in Go Playground |
| `helmsettag` | finds a Helm release by name and sets its image tag to the desired one |
| `kubecred` | simplified certificate signing process on Kubernetes cluster |
| `kubectx` | makes kubectl context switching faster |
| `kubewrap` | wraps CLI tools to make them connect to Kubernetes pods directly: ```kubewrap curl nginx-6lapfe/healthz``` |
| `mkremote` | adds default git remote based on current Go project location |
| `nodepod` | looks up full Kubernetes pod name by prefix and node |
| `ppod` | fetches CPU profile and trace of Go app by podname it is running in |
| `redir` | spins up an http server that redirects to provided url |
| `sembump` | makes git tag semver bumping a no-brainer |
| `whatonion` | calculates onion address from private key file |
