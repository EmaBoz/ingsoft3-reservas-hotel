# Inicialización de Solr para search-api

Este directorio contiene los scripts y configuraciones necesarios para inicializar Solr con el schema correcto para indexar habitaciones.

## Schema de Solr (schema.json)

El schema define los siguientes campos indexados:

- **id** (string): ID único de la habitación
- **number** (string): Número de habitación
- **type** (string): Tipo de habitación (single, double, suite, etc.)
- **status** (string): Estado (available, occupied, maintenance, reserved)
- **price** (double): Precio por noche
- **description** (text): Descripción de la habitación
- **capacity** (int): Capacidad de personas
- **floor** (int): Piso
- **has_wifi** (boolean): Tiene WiFi
- **has_ac** (boolean): Tiene aire acondicionado
- **has_tv** (boolean): Tiene TV
- **has_minibar** (boolean): Tiene minibar
- **created_at** (date): Fecha de creación
- **updated_at** (date): Fecha de última actualización

## Configuración manual del schema

Si el core ya fue creado por Docker pero necesitas actualizar el schema manualmente, puedes hacerlo con curl:

```bash
# Aplicar schema al core rooms-core
curl -X POST -H 'Content-type:application/json' \
  http://localhost:8983/solr/rooms-core/schema \
  --data-binary @schema.json
```

## Verificar el schema

```bash
# Ver todos los campos del schema
curl http://localhost:8983/solr/rooms-core/schema/fields

# Verificar health del core
curl http://localhost:8983/solr/rooms-core/admin/ping
```

## Consulta de ejemplo

```bash
# Buscar todas las habitaciones
curl "http://localhost:8983/solr/rooms-core/select?q=*:*&wt=json"

# Buscar habitaciones disponibles con WiFi
curl "http://localhost:8983/solr/rooms-core/select?q=*:*&fq=status:available%20AND%20has_wifi:true&wt=json"
```

## Re-indexar todas las habitaciones

Si necesitas re-indexar todas las habitaciones desde rooms-api, puedes hacerlo manualmente:

```bash
# 1. Eliminar todos los documentos del índice
curl http://localhost:8983/solr/rooms-core/update?commit=true \
  -H "Content-Type: application/json" \
  -d '{"delete":{"query":"*:*"}}'

# 2. Luego, crear/actualizar habitaciones en rooms-api
# Los eventos de RabbitMQ automáticamente indexarán en Solr
```
