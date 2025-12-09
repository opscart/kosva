# Using kosva with Kubecost

## Prerequisites

You need Kubecost installed in your cluster. If you don't have it:
```bash
# Install Kubecost (free tier)
helm repo add kubecost https://kubecost.github.io/cost-analyzer/
helm install kubecost kubecost/cost-analyzer \
    --namespace kubecost --create-namespace \
    --set kubecostToken="aGVsbUBrdWJlY29zdC5jb20=xm343yadf98"
```

## Using kosva with Kubecost

### Option 1: Direct API Access (Recommended)

If kosva runs inside your cluster:
```bash
kosva check --kubecost-url http://kubecost-cost-analyzer.kubecost.svc:9090
```

If you have port-forwarded Kubecost to localhost:
```bash
kubectl port-forward -n kubecost svc/kubecost-cost-analyzer 9090:9090

# In another terminal:
kosva check --kubecost-url http://localhost:9090
```

### Option 2: Export and Validate (Offline)
```bash
# 1. Export from Kubecost UI
# Navigate to: Kubecost UI → Savings → Export as JSON
# Save as: kubecost-recommendations.json

# 2. Validate
kosva check --input kubecost-recommendations.json
```

### Option 3: Kubecost API Export
```bash
# Get your Kubecost endpoint
KUBECOST_URL="http://localhost:9090"

# Fetch recommendations
curl "${KUBECOST_URL}/model/savings" > recommendations.json

# Validate
kosva check --input recommendations.json
```

## Example Workflow
```bash
# 1. Install Kubecost (if not already)
helm install kubecost kubecost/cost-analyzer -n kubecost --create-namespace

# 2. Wait for it to collect data (15-30 minutes)
kubectl wait -n kubecost --for=condition=ready pod -l app=cost-analyzer

# 3. Port-forward
kubectl port-forward -n kubecost svc/kubecost-cost-analyzer 9090:9090 &

# 4. Run kosva validation
kosva check --kubecost-url http://localhost:9090

# 5. Review safe recommendations and implement
```

## Troubleshooting

### "Connection refused"
- Ensure Kubecost is running: `kubectl get pods -n kubecost`
- Verify port-forward: `netstat -an | grep 9090`

### "No recommendations found"
- Kubecost needs 15-30 minutes to collect data after installation
- Check Kubecost UI to see if recommendations are available

### "API returned 404"
- Ensure you're using the correct endpoint: `/model/savings`
- Check Kubecost version (need v1.90+)