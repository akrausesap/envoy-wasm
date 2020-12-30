# Important Inputs


WASM Golang SDK: https://github.com/tetratelabs/proxy-wasm-go-sdk & https://pkg.go.dev/github.com/tetratelabs/proxy-wasm-go-sdk@v0.12.0

WASM & Istio/Envoy: https://banzaicloud.com/blog/envoy-wasm-filter/

Setup Kind & Istio: https://www.danielstechblog.io/running-istio-on-kind-kubernetes-in-docker/

# Deploy httpbin sample: 
```
kubectl label namespace default istio-injection=enabled  
kubectl apply -f istio/httpbin_manifest.yaml
```

# Build & Deploy WASM


```
make build
kubectl create cm my-filter --from-file=my-filter.wasm
kubectl apply -f istio/filter.yaml
```