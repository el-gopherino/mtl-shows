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
| `GET /` | All upcoming events |
| `GET /tonight` | Events happening today |
| `GET /tomorrow` | Events happening tomorrow |
| `GET /this-week` | Events within the next 7 days |
| `GET /this-weekend` | Events this Friday–Sunday |
| `GET /right-now` | Events starting soon or in progress |

## Venues

- Cafe Campus
- Casa del Popolo
- Club Soda
- L'Hemisphere Gauche
- MTelus
- Piranha Bar
- P'tit Ours
- Quai des Brumes
- La Sala Rossa
- La Sotterenea
- La Toscadura
- Turbo Haüs
- Le Verre Bouteille
