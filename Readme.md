# Important Things

WASM & Istio/Envoy: https://banzaicloud.com/blog/envoy-wasm-filter/

Setup Kind & Istio: https://www.danielstechblog.io/running-istio-on-kind-kubernetes-in-docker/

Deploy httpbin: 
```
kubectl label namespace default istio-injection=enabled  
kubectl apply -f isti/manifest.yaml
```

# Deploy WASM

```
kubectl create cm my-filter --from-file=my-filter.wasm
kubect apply -f filter.yaml
```