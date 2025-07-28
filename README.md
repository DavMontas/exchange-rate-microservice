# Exchange Rate Offers Microservice

Este microservicio consulta y compara tasas de cambio entre tres proveedores externos en tiempo real, devolviendo la mejor oferta.

## 📦 Estructura del Proyecto
cmd/server/main.go # Entry point del servidor
- `cmd/server/main.go`: Punto de entrada del servidor HTTP
- `internal/application/`: Configuración y registro de dependencias
- `internal/exchangerate/`
  - `client/`: Clientes para cada proveedor (JSON, XML)
  - `service/`: Lógica de negocio (comparación y agregación de tasas)
  - `transport/`: API HTTP expuesta mediante Gin


1. **Clona el repositorio**  
    ```bash
    git clone https://github.com/DavMontas/exchange-rate-microservice.git
    ```
2. **Ejecuta el servidor**  
  ```bash
  go run cmd/server/main.go
  ```

3. **Haz una prueba** 
  ```bash
  curl -X POST http://localhost:8080/best-quote \
  -H "Content-Type: application/json" \
  -d '{"from":"USD","to":"EUR","amount":100}'
  ```


**Notes:** 
  Para ejecutar los tests: 

  ```bash
  go test -cover ./internal/exchangerate/client                                    
  go test -cover ./internal/exchangerate/service
  ```

***🐳 Usando Docker***
1. `Construye La imagen`
  ```bash
  docker build -t exchange-offers:latest .
  ```

2. `Ejecutar el container`
  ```bash
  docker run -p 8080:8080 exchange-offers:latest
  ```