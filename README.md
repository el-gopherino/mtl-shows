# mtl-shows

Scrapes and serves information about live shows and events in Montreal's independent venues.

## Usage

```bash
go build -o mtl-shows .

./mtl-shows -serve  # scrape concurrently and run API server on :8080 with scheduler
./mtl-shows -conc   # scrape once (concurrent)
./mtl-shows -seq    # scrape once (sequential)
```

## API Endpoints

| Endpoint | Description |
|---|---|
| `GET /events` | All upcoming events |
| `GET /events/tonight` | Events happening today |
| `GET /events/tomorrow` | Events happening tomorrow |
| `GET /events/this-week` | Events within the next 7 days |
| `GET /events/this-weekend` | Events this Friday–Sunday |
| `GET /events/right-now` | Events starting soon or in progress |

## Venues

- Casa del Popolo
- La Sala Rossa
- La Sotterenea
- P'tit Ours
- La Toscadura
- Quai des Brumes
- Cafe Campus
- L'Hemisphere Gauche
- Le Verre Bouteille
- Turbo Haüs