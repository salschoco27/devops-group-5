#!/usr/bin/env bash
# Golden Path Deployment Measurement Script
# Measures the deployment time from manifest application until all application Pods reach the Ready state.

set -e

MANIFEST_PATH="${1:-}"
NAMESPACE="${2:-golden-path-dev}"
APP_LABEL="${3:-app=sample-app}"

if [ -z "$MANIFEST_PATH" ]; then
  echo "Usage:"
  echo "./measure.sh <manifest-path> [namespace] [label]"
  echo
  echo "Example:"
  echo "./measure.sh ../implementation/kubernetes/deployment.yaml golden-path-dev app=sample-app"
  exit 1
fi

echo
echo "══════════════════════════════════════════════════════════════════════"
echo "             ●─────────────────────★─────────────────────●           "
echo "                        G O L D E N   P A T H"
echo
echo "══════════════════════════════════════════════════════════════════════"
echo
echo "Timestamp : $(date)"
echo "Manifest  : $MANIFEST_PATH"
echo "Namespace : $NAMESPACE"
echo "Label     : $APP_LABEL"
echo

START_TIME=$(date +%s)

kubectl apply -f "$MANIFEST_PATH"

echo
echo "⏳ Waiting for Pods to become Ready..."
echo

kubectl wait \
  --for=condition=ready pod \
  -l "$APP_LABEL" \
  -n "$NAMESPACE" \
  --timeout=180s

END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo 
echo "══════════════════════════════════════════════════════════════════════"
echo "                      📝 Deployment Completed.                        "
echo "══════════════════════════════════════════════════════════════════════"
echo "Deployment Time : ${DURATION} second(s)"
echo "Status          : SUCCESS"
echo

echo "Current Pod Status"
echo "─────────────────────────────────────────────────────────────────────"
kubectl get pods -n "$NAMESPACE"
echo
echo
echo "══════════════════════════════════════════════════════════════════════"
echo "                 ✅ Measurement finished successfully.                "
echo "══════════════════════════════════════════════════════════════════════"