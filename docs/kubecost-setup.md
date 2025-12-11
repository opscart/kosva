# Kubecost Integration Guide

This guide covers integrating kosva with Kubecost for validating cost optimization recommendations.

## Prerequisites

- Kubernetes cluster (1.19+)
- kubectl configured
- Helm 3.x installed

## Option 1: Install Kubecost (If Not Already Installed)

### Free Tier Installation
```bash
# Add Kubecost Helm repository
helm repo add kubecost https://kubecost.github.io/cost-analyzer/
helm repo update

# Install Kubecost with free tier
helm install kubecost kubecost/cost-analyzer \
    --namespace kubecost \
    --create-namespace \
    --set kubecostToken="aGVsbUBrdWJlY29zdC5jb20=xm343yadf98"

# Wait for deployment
kubectl wait -n kubecost \
    --for=condition=ready pod \
    -l app=cost-analyzer \
    --timeout=300s
```

### Verify Installation
```bash
# Check pods
kubectl get pods -n kubecost

# Port-forward to access UI
kubectl port-forward -n kubecost svc/kubecost-cost-analyzer 9090:9090

# Open browser
open http://localhost:9090
```

**Note:** Kubecost needs 15-30 minutes after installation to collect initial metrics before generating recommendations.

---

## Option 2: Using Existing Kubecost Installation

If you already have Kubecost installed, skip to the integration section below.

---

## Integration Methods

### Method 1: Export from Kubecost UI (Recommended for Manual Validation)

**Steps:**

1. Access Kubecost UI: `http://localhost:9090`
2. Navigate to **Savings** → **Right-sizing** or **Cluster Right-sizing**
3. Click **Export** → **JSON**
4. Save as `kubecost-recommendations.json`
5. Run kosva:
```bash
kosva check --input kubecost-recommendations.json
```

**Pros:**
- Simple, no additional setup
- Works with any Kubecost version
- Good for one-time validations

**Cons:**
- Manual process
- Not suitable for automation

---

### Method 2: Kubecost API Export (Recommended for Automation)

**Steps:**
```bash
# Port-forward Kubecost service
kubectl port-forward -n kubecost svc/kubecost-cost-analyzer 9090:9090 &

# Export recommendations via API
curl http://localhost:9090/model/savings > kubecost-recs.json

# Validate with kosva
kosva check --input kubecost-recs.json
```

**API Endpoints:**

- **Savings recommendations:** `/model/savings`
- **Right-sizing:** `/model/resourceRecommendations`
- **Cluster optimization:** `/model/clusterOptimization`

**Pros:**
- Automatable
- Works in CI/CD pipelines
- No manual UI interaction

**Cons:**
- Requires API access
- Network connectivity needed

---

### Method 3: Direct API Integration (Recommended for Real-time Validation)

**Steps:**
```bash
# Option A: From within cluster (if kosva runs as pod)
kosva check --kubecost-url http://kubecost-cost-analyzer.kubecost.svc:9090

# Option B: From local machine (via port-forward)
kubectl port-forward -n kubecost svc/kubecost-cost-analyzer 9090:9090 &
kosva check --kubecost-url http://localhost:9090

# Option C: From remote machine (if Kubecost exposed via ingress)
kosva check --kubecost-url https://kubecost.company.com
```

**Pros:**
- Real-time validation
- No intermediate files
- Direct integration

**Cons:**
- Requires API access
- Network dependency

---

## CI/CD Pipeline Integration

### GitHub Actions
```yaml
name: Cost Validation

on:
  schedule:
    - cron: '0 9 * * 1'  # Weekly on Monday 9 AM
  workflow_dispatch:

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup kubectl
        uses: azure/setup-kubectl@v3
        with:
          version: 'v1.28.0'
      
      - name: Configure kubectl
        run: |
          echo "${{ secrets.KUBECONFIG }}" > kubeconfig
          export KUBECONFIG=./kubeconfig
      
      - name: Install kosva
        run: |
          wget https://github.com/opscart/kosva/releases/latest/kosva
          chmod +x kosva
      
      - name: Port-forward Kubecost
        run: |
          kubectl port-forward -n kubecost svc/kubecost-cost-analyzer 9090:9090 &
          sleep 5
      
      - name: Export recommendations
        run: |
          curl http://localhost:9090/model/savings > recommendations.json
      
      - name: Validate
        run: |
          ./kosva check --input recommendations.json
      
      - name: Upload report
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: validation-report
          path: recommendations.json
```

---

### Azure DevOps
```yaml
trigger: none

schedules:
- cron: "0 9 * * 1"
  displayName: Weekly cost validation
  branches:
    include:
    - main

pool:
  vmImage: 'ubuntu-latest'

steps:
- task: Kubernetes@1
  displayName: 'Setup kubectl'
  inputs:
    connectionType: 'Kubernetes Service Connection'
    kubernetesServiceEndpoint: '$(kubernetesServiceConnection)'

- bash: |
    # Port-forward Kubecost
    kubectl port-forward -n kubecost svc/kubecost-cost-analyzer 9090:9090 &
    sleep 5
    
    # Export recommendations
    curl http://localhost:9090/model/savings > recommendations.json
  displayName: 'Export Kubecost recommendations'

- bash: |
    # Download kosva
    wget https://github.com/opscart/kosva/releases/latest/kosva
    chmod +x kosva
    
    # Run validation
    ./kosva check --input recommendations.json
  displayName: 'Validate recommendations'

- task: PublishBuildArtifacts@1
  condition: always()
  inputs:
    pathToPublish: 'recommendations.json'
    artifactName: 'cost-validation'
```

---

### GitLab CI
```yaml
cost-validation:
  stage: validate
  image: bitnami/kubectl:latest
  script:
    - kubectl config use-context $KUBE_CONTEXT
    - kubectl port-forward -n kubecost svc/kubecost-cost-analyzer 9090:9090 &
    - sleep 5
    - curl http://localhost:9090/model/savings > recommendations.json
    - wget https://github.com/opscart/kosva/releases/latest/kosva
    - chmod +x kosva
    - ./kosva check --input recommendations.json
  artifacts:
    paths:
      - recommendations.json
    expire_in: 1 week
  only:
    - schedules
```

---

## Troubleshooting

### Issue: "Connection refused" when accessing Kubecost API

**Symptoms:**
```
ERROR: failed to connect to Kubecost API
```

**Solutions:**

1. **Verify Kubecost is running:**
```bash
kubectl get pods -n kubecost
```

2. **Check service:**
```bash
kubectl get svc -n kubecost
```

3. **Verify port-forward:**
```bash
# Kill existing port-forwards
killall kubectl

# Start fresh port-forward
kubectl port-forward -n kubecost svc/kubecost-cost-analyzer 9090:9090

# Test connectivity
curl http://localhost:9090/model/savings
```

4. **Check firewall/network policies:**
```bash
# Test from within cluster
kubectl run -it --rm debug --image=curlimages/curl --restart=Never -- \
  curl http://kubecost-cost-analyzer.kubecost.svc:9090/model/savings
```

---

### Issue: "No recommendations found"

**Symptoms:**
```
ERROR: recommendations list is empty
```

**Causes:**
- Kubecost hasn't collected enough data yet (needs 15-30 minutes after install)
- No cost optimization opportunities detected
- API endpoint incorrect

**Solutions:**

1. **Wait for data collection:**
```bash
# Check Kubecost logs
kubectl logs -n kubecost -l app=cost-analyzer --tail=50

# Look for: "Data collection complete"
```

2. **Verify data in UI:**
- Open Kubecost UI: `http://localhost:9090`
- Check if costs are showing
- Navigate to Savings tab

3. **Try different endpoint:**
```bash
# Try right-sizing endpoint
curl http://localhost:9090/model/resourceRecommendations
```

---

### Issue: "API returned 404"

**Symptoms:**
```
ERROR: kubecost API returned 404
```

**Causes:**
- Wrong Kubecost version (old API endpoints)
- Incorrect URL path
- Kubecost not fully initialized

**Solutions:**

1. **Check Kubecost version:**
```bash
kubectl get deployment -n kubecost kubecost-cost-analyzer -o yaml | grep image:
```

2. **Use correct endpoint for your version:**
- Kubecost 1.90+: `/model/savings`
- Older versions: `/model/resourceRecommendations`

3. **Verify API is ready:**
```bash
# Check readiness
kubectl get pods -n kubecost -o wide

# All pods should be Running and Ready (1/1)
```

---

### Issue: "Invalid JSON format"

**Symptoms:**
```
ERROR: failed to parse response: invalid JSON
```

**Causes:**
- Kubecost API returned HTML error page instead of JSON
- API authentication required
- Malformed response

**Solutions:**

1. **Check raw response:**
```bash
curl -v http://localhost:9090/model/savings | head -20
```

2. **Verify content type:**
```bash
curl -I http://localhost:9090/model/savings
# Should return: Content-Type: application/json
```

3. **Save and inspect:**
```bash
curl http://localhost:9090/model/savings > debug.json
cat debug.json | jq .  # Pretty-print JSON
```

---

## Advanced Configuration

### Running kosva as Kubernetes CronJob
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: kosva-validation
  namespace: kubecost
spec:
  schedule: "0 9 * * 1"  # Weekly on Monday 9 AM
  jobTemplate:
    spec:
      template:
        spec:
          serviceAccountName: kosva
          containers:
          - name: kosva
            image: your-registry/kosva:latest
            command:
            - /bin/sh
            - -c
            - |
              kosva check \
                --kubecost-url http://kubecost-cost-analyzer:9090 \
                --policy-dir /policies
          restartPolicy: OnFailure
          volumes:
          - name: policies
            configMap:
              name: kosva-policies
```

**Create ServiceAccount:**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kosva
  namespace: kubecost
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: kosva
  namespace: kubecost
rules:
- apiGroups: [""]
  resources: ["services"]
  verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: kosva
  namespace: kubecost
subjects:
- kind: ServiceAccount
  name: kosva
roleRef:
  kind: Role
  name: kosva
  apiGroup: rbac.authorization.k8s.io
```

---

## Performance Considerations

### API Response Times

- **Small clusters** (<50 nodes): < 2 seconds
- **Medium clusters** (50-200 nodes): 2-10 seconds
- **Large clusters** (200+ nodes): 10-30 seconds

### Caching

Kubecost caches recommendations for 1 hour by default. For more frequent validation:
```bash
# Force refresh by adding timestamp
curl "http://localhost:9090/model/savings?timestamp=$(date +%s)"
```

### Rate Limiting

When running in CI/CD, avoid making too many API calls:
```bash
# Cache the export
if [ ! -f recommendations.json ] || [ $(find recommendations.json -mmin +60) ]; then
  curl http://localhost:9090/model/savings > recommendations.json
fi

kosva check --input recommendations.json
```

---

## Next Steps

- [Policy Customization Guide](policy-guide.md)
- [Contributing](../CONTRIBUTING.md)
- [Example Policies](../policies/)