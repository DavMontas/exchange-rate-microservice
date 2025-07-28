# Exchange Rate Offers Microservice

Este microservicio consulta y compara tasas de cambio entre tres proveedores externos en tiempo real, devolviendo la mejor oferta.

## 📦 Estructura del Proyecto
cmd/server/main.go # Entry point del servidor
internal/application/ # Configuración y registro de dependencias
internal/exchangerate/
client/ # Implementación de cada proveedor (JSON, XML)
service/ # Lógica de negocio (comparación de rates)
transport/ # HTTP API (Gin)


1. **Clona el repositorio**  

2. **Ejecuta el servidor**  
  ```bash
  go run cmd/server/main.go

3. **Haz una prueba** 
  ```bash
  curl -X POST http://localhost:8080/best-quote \
  -H "Content-Type: application/json" \
  -d '{"from":"USD","to":"EUR","amount":100}'


**Notes:** 
  Para ejecutar los tests: 