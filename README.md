# Bus Booking System

_Aplicación e implementación del algoritmo de Suzuki-Kasami para asegurar exclusión mutua distribuida en un sistema de reserva de pasajes._

_Contexto:_ En buses con 48 asientos se asumen 3 tramos de precios:
- Tramo 1, asientos 01-16: $8000 CLP
- Tramo 2, asientos 17-32: $6000 CLP
- Tramo 3, asientos 33-48: $4000 CLP

---

- El archivo **mapa.txt** contiene una visualización de la distribución de asientos.
- **pasajeros.txt** contiene la lista de pasajeros a procesar, cada vez que se procese uno, se eliminara del archivo.
- **ganancias.txt** contiene la ganancia por tramos.
- **procesados.txt** contiene el log de la ejecución.

# Ejecución:

```
go build main.go
./main <n° process>
```

# Equipo:
- Rafael Aros Soto
- Franco Cabezas Poblete
- Paulina Vega Rivera

# Bugs conocidos:
- Replicación de compra
