package main

type Venue struct {
	Name           string
	Group          string // group related venues (multiple venues for a single website) this is rare
	Link           string
	Website        string
	AllowedDomains []string
	Selector       string // CSS selector for events on page

	// for interactive map
	Latitude  float64
	Longitude float64

	Address      string
	Neighborhood string
}

// venue websites' CSS selectors for colly scraper -> event info + logo
const (
	CasaDelPopoloSelector    = "div.md\\:w-5\\/12.p-6"
	QuaiDesBrumesSelector    = "article.mec-event-article"
	CafeCampusSelector       = "div.noo-shevent-content"
	HemisphereGaucheSelector = "div.IFphb0"
	VerreBouteilleSelector   = "div.card-container"
	PiranhaBarSelector       = "article.eventlist-event"
	ClubSodaSelector         = "div.card.h-100"
)

const (
	TurboHausURL = "https://www.turbohaus.ca/calendrier?format=json"
)

var allVenues = map[string]Venue{
	"casa-del-popolo": {
		Name:    "Casa del Popolo",
		Address: "4873 Boul. St-Laurent",
		Group:   "casa-del-popolo-group",

		Link: "https://casadelpopolo.com/en/events/casa-del-popolo",
		AllowedDomains: []string{
			"casadelpopolo.com",
			"www.casadelpopolo.com",
		},
		Selector: CasaDelPopoloSelector,
		Website:  "https://casadelpopolo.com/en",

		Latitude:  45.521805,
		Longitude: -73.590431,
	},

	"la-sala-rossa": {
		Name:    "La Sala Rossa",
		Address: "4848 Boul. St-Laurent",
		Group:   "casa-del-popolo-group",

		Link: "https://casadelpopolo.com/en/events/la-sala-rossa",
		AllowedDomains: []string{
			"casadelpopolo.com",
			"www.casadelpopolo.com",
		},
		Selector: CasaDelPopoloSelector,
		Website:  "https://casadelpopolo.com/en",

		Latitude:  45.521771,
		Longitude: -73.590493,
	},

	"la-sotterenea": {
		Name:    "La Sotterenea",
		Address: "4848 Boul. St-Laurent",
		Group:   "casa-del-popolo-group",

		Link: "https://casadelpopolo.com/en/events/la-sotterenea",
		AllowedDomains: []string{
			"casadelpopolo.com",
			"www.casadelpopolo.com",
		},

		Selector: CasaDelPopoloSelector,
		Website:  "https://casadelpopolo.com/en",

		Latitude:  45.521771,
		Longitude: -73.590493,
	},

	"ptit-ours": {
		Name:    "Ptit Ours",
		Address: "5589 Avenue du Parc",
		Group:   "casa-del-popolo-group",

		Link: "https://casadelpopolo.com/en/events/ptit-ours",
		AllowedDomains: []string{
			"casadelpopolo.com",
			"www.casadelpopolo.com",
		},

		Selector: CasaDelPopoloSelector,
		Website:  "https://casadelpopolo.com/en",

		Latitude:  45.522644,
		Longitude: -73.602695,
	},

	"la-toscadura": {
		Name:    "La Toscadura",
		Address: "4388 St-Laurent",
		Group:   "casa-del-popolo-group",

		Link: "https://casadelpopolo.com/en/events/la-toscadura",
		AllowedDomains: []string{
			"casadelpopolo.com",
			"www.casadelpopolo.com",
		},

		Selector: CasaDelPopoloSelector,
		Website:  "https://casadelpopolo.com/en",

		Latitude:  45.519246,
		Longitude: -73.584909,
	},

	"quai-des-brumes": {
		Name:    "Quai des Brumes",
		Address: "4481 Rue Saint-Denis",

		Link: "https://quaidesbrumes.ca/calendrier/",
		AllowedDomains: []string{
			"quaidesbrumes.ca",
			"www.quaidesbrumes.ca",
		},

		Selector: QuaiDesBrumesSelector,
		Website:  "https://quaidesbrumes.ca",

		Latitude:  45.523917,
		Longitude: -73.582513,
	},

	"cafe-campus": {
		Name:    "Cafe Campus",
		Address: "57 Rue Prince-Arthur Est",

		Link: "https://spectacles.cafecampus.com/evenements",
		AllowedDomains: []string{
			"spectacles.cafecampus.com",
			"www.spectacles.cafecampus.com",
		},

		Selector: CafeCampusSelector,
		Website:  "https://www.cafecampus.com/",

		Latitude:  45.514541,
		Longitude: -73.572183,
	},

	"hemisphere-gauche": {
		Name:    "L'Hemisphere Gauche",
		Address: "221 Beaubien Est",

		Link: "https://www.hemispheregauche.com/?lang=en",
		AllowedDomains: []string{
			"hemispheregauche.com",
			"www.hemispheregauche.com",
		},

		Selector: HemisphereGaucheSelector,
		Website:  "https://www.hemispheregauche.com",

		Latitude:  45.532241,
		Longitude: -73.606866,
	},

	"verre-bouteille": {
		Name:    "Le Verre Bouteille",
		Address: "2112 Avenue du Mont-Royal Est",

		Link: "https://verrebouteille.com/shows.php",
		AllowedDomains: []string{
			"verrebouteille.com",
			"www.verrebouteille.com",
		},

		Selector: VerreBouteilleSelector,
		Website:  "https://verrebouteille.com",

		Latitude:  45.535373,
		Longitude: -73.572007,
	},

	"piranha-bar": {
		Name:    "Piranha Bar",
		Address: "680 Rue Sainte-Catherine Ouest",

		Link: "https://www.piranhabar.ca/events",
		AllowedDomains: []string{
			"piranhabar.ca",
			"www.piranhabar.ca",
		},

		Selector: PiranhaBarSelector,
		Website:  "https://www.piranhabar.ca",

		Latitude:  45.502818,
		Longitude: -73.569794,
	},

	"club-soda": {
		Name:    "Club Soda",
		Address: "1225 Boul. Saint-Laurent",

		Link: "https://clubsoda.ca/fr/evenements",
		AllowedDomains: []string{
			"clubsoda.ca",
			"www.clubsoda.ca",
		},

		Selector: ClubSodaSelector,
		Website:  "https://clubsoda.ca",

		Latitude:  45.509597,
		Longitude: -73.563217,
	},

	"mtelus": {
		Name:    "MTelus",
		Address: "59 Rue Sainte-Catherine Est",

		Link: "https://mtelus.com/en/events?display=list",
		AllowedDomains: []string{
			"mtelus.com",
			"www.mtelus.com",
		},

		// no Selector - parse JSON (JS-rendered)
		Website: "https://mtelus.com",

		Latitude:  45.510586,
		Longitude: -73.56321,
	},

	// JS-rendered
	"turbo-haus": {
		Name:    "Turbo Haüs",
		Address: "2040 Rue Saint-Denis",

		Link: "https://www.turbohaus.ca/cal",
		AllowedDomains: []string{
			"turbohaus.ca",
			"www.turbohaus.ca",
		},

		// no Selector - parse JSON
		Website: "https://www.turbohaus.ca",

		Latitude:  45.516304,
		Longitude: -73.566101,
	},
}
