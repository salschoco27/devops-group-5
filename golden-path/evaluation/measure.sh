#!/usr/bin/env bash
# Measurement script untuk evaluasi Golden Path.
# Mengukur waktu dari manifest diterapkan sampai pod aplikasi Ready.

set -e

MANIFEST_PATH="${1:-}"
NAMESPACE="${2:-golden-path-dev}"
APP_LABEL="${3:-app=sample-app}"

if [ -z "$MANIFEST_PATH" ]; then
  echo "Cara pakai:"
  echo "./measure.sh <path-manifest> [namespace] [label]"
  echo "Contoh:"
  echo "./measure.sh ../implementation/kubernetes/deployment.yaml golden-path-dev app=sample-app"
  exit 1
fi

echo "=== GOLDEN PATH DEPLOYMENT TIME MEASUREMENT ==="
echo "Timestamp : $(date)"
echo "Manifest  : $MANIFEST_PATH"
echo "Namespace : $NAMESPACE"
echo "Label     : $APP_LABEL"
echo ""

START_TIME=$(date +%s)

kubectl apply -f "$MANIFEST_PATH"

echo ""
echo "Menunggu pod Ready..."

kubectl wait \
  --for=condition=ready pod \
  -l "$APP_LABEL" \
  -n "$NAMESPACE" \
  --timeout=180s

END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo ""
echo "=== HASIL PENGUKURAN ==="
echo "Waktu deployment sampai pod Ready: ${DURATION} detik"
echo ""
echo "Status pod:"
kubectl get pods -n "$NAMESPACE"