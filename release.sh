#!/bin/bash

for pod in $(kubectl get pods -o jsonpath="{.items[*].metadata.name}" -l app=$1 -n swiply); do
  kubectl delete pod $pod -n swiply
  sleep 30
  kubectl get pods -n swiply -l app=$1
  sleep 30
done

echo "FINISHED!"
echo "Check logs:"
kubectl logs -l app=$1 -n swiply --tail=20