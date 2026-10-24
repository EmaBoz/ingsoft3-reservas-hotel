#!/bin/sh
# Corre los tests de reservations-api, mide la cobertura y falla si queda por
# debajo del umbral. Es el ENTRYPOINT de la etapa "test" del Dockerfile.
set -e -o pipefail

# Umbral sobre cobertura de sentencias (Go no mide ramas). Justificado en decisiones.md.
UMBRAL=60
OUT=${OUT_DIR:-/out}
mkdir -p "$OUT"

# Se mide todo el módulo y después se saca lo que no entra en la cuenta:
# cmd y config (arranque), domain (datos), repositories y events (Mongo y
# RabbitMQ, eso lo verifica la integración). Un paquete nuevo entra solo.
CGO_ENABLED=0 go test -v -count=1 ./... -coverpkg=./... -coverprofile=/tmp/cobertura-completa.out | tee "$OUT/tests.txt"
grep -v -E "reservations-api/(cmd|config|domain|repositories|events)/" /tmp/cobertura-completa.out > "$OUT/coverage.out"

go tool cover -func="$OUT/coverage.out" > "$OUT/coverage.txt"
go tool cover -html="$OUT/coverage.out" -o "$OUT/coverage.html"

TOTAL=$(tail -n 1 "$OUT/coverage.txt" | awk '{print $NF}' | tr -d '%')
echo "Cobertura de sentencias: $TOTAL% (umbral: $UMBRAL%)"

if awk "BEGIN { exit !($TOTAL < $UMBRAL) }"; then
  echo "ERROR: la cobertura de sentencias ($TOTAL%) no llega al umbral ($UMBRAL%)"
  exit 1
fi
