# mtl-shows

Scrapes and serves information about live shows and events in Montreal's independent venues.

## Requirements

- Go 1.26.4 or later
- Git

## Getting started

Clone the repo and fetch dependencies:

```sh
git clone https://github.com/el-gopherino/mtl-shows.git
cd mtl-shows
go mod download
```

## Running

The binary has three mutually exclusive modes, selected by flag:

```sh
# One-shot scrape, sequential (one venue at a time)
go run . -seq

# One-shot scrape, concurrent (all venues in parallel)
go run . -conc

# Start the API server on :8080, scrape on startup, then re-scrape every hour
go run . -serve
```

When running with `-serve`, open http://localhost:8080 in your browser. Available routes:

- `/` — all events
- `/right-now`
- `/tonight`
- `/tomorrow`
- `/this-week`
- `/this-weekend` (also `/this-weekend/friday`, `/saturday`, `/sunday`)

You can also build a binary instead of using `go run`:

```sh
go build -o mtl-shows
./mtl-shows -serve
```

Scraped event data is written to JSON files in the project directory (`output/`, `tonight/`, `tomorrow/`, etc.).

## Current venues scraped

- Bar Le Ritz
- Cafe Campus
- Casa del Popolo
- Club Soda
- Fairmount Theatre
- L'Hemisphere Gauche
- Le Ministère
- MTelus
- Piranha Bar
- P'tit Ours
- Quai des Brumes
- La Sala Rossa
- La Sotterenea
- La Toscadura
- Turbo Haüs
- Le Verre Bouteille
