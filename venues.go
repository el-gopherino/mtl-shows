package main

type Venue struct {
	Name           string
	Group          string // group related venues (multiple venues for a single website)
	Links          []string
	AllowedDomains []string
	Selector       string // CSS selector for events on page
	Address        string
	Neighborhood   string
	Website        string
}

// venue websites' CSS selectors and for event info and event logo
const (
	CasaDelPopoloSelector    = "div.md\\:w-5\\/12.p-6"
	QuaiDesBrumesSelector    = "article.mec-event-article"
	CafeCampusSelector       = "div.noo-shevent-content"
	HemisphereGaucheSelector = "div.IFphb0"
	BarLeRitzPDBSelector     = "article.eventlist-event eventlist-event--upcoming eventlist-event--hasimg eventlist-hasimg is-loaded"
)

var allVenues = map[string]Venue{

	// --------------------------------- CASA group ---------------------------------
	"casa-del-popolo": {
		Name:    "Casa del Popolo",
		Address: "4873 St-Laurent",
		Group:   "casa",
		Links: []string{
			"https://casadelpopolo.com/en/events/casa-del-popolo",
		},
		AllowedDomains: []string{
			"casadelpopolo.com",
			"www.casadelpopolo.com",
		},
		Selector: CasaDelPopoloSelector,
	},

	"la-sala-rossa": {
		Name:    "La Sala Rossa",
		Address: "4848 St-Laurent",
		Group:   "casa",
		Links: []string{
			"https://casadelpopolo.com/en/events/la-sala-rossa",
		},
		AllowedDomains: []string{
			"casadelpopolo.com",
			"www.casadelpopolo.com",
		},
		Selector: CasaDelPopoloSelector,
	},

	"la-sotterenea": {
		Name:    "La Sotterenea",
		Address: "4848 St-Laurent",
		Group:   "casa",
		Links: []string{
			"https://casadelpopolo.com/en/events/la-sotterenea",
		},
		AllowedDomains: []string{
			"casadelpopolo.com",
			"www.casadelpopolo.com",
		},
		Selector: CasaDelPopoloSelector,
	},

	"ptit-ours": {
		Name:    "Ptit Ours",
		Address: "5589 Avenue du Parc",
		Group:   "casa",
		Links: []string{
			"https://casadelpopolo.com/en/events/ptit-ours",
		},
		AllowedDomains: []string{
			"casadelpopolo.com",
			"www.casadelpopolo.com",
		},
		Selector: CasaDelPopoloSelector,
	},

	"la-toscadura": {
		Name:    "La Toscadura",
		Address: "4388 St-Laurent",
		Group:   "casa",
		Links: []string{
			"https://casadelpopolo.com/en/events/la-toscadura",
		},
		AllowedDomains: []string{
			"casadelpopolo.com",
			"www.casadelpopolo.com",
		},
		Selector: CasaDelPopoloSelector,
	},
	// ------------------------------------- Fin Casa ---------------------------------------------------------

	"quai-des-brumes": {
		Name:    "Quai des Brumes",
		Address: "4481 Rue Saint-Denis",
		Links: []string{
			"https://quaidesbrumes.ca/calendrier/",
		},
		AllowedDomains: []string{
			"quaidesbrumes.ca",
			"www.quaidesbrumes.ca",
		},
		Selector: QuaiDesBrumesSelector,
		Website:  "https://quaidesbrumes.ca",
	},

	"cafe-campus": {
		Name:    "Cafe Campus",
		Address: "57 Rue Prince-Arthur Est",
		Links: []string{
			"https://spectacles.cafecampus.com/evenements",
		},
		AllowedDomains: []string{
			"spectacles.cafecampus.com",
			"www.spectacles.cafecampus.com",
		},
		Selector: CafeCampusSelector,
		Website:  "https://www.cafecampus.com/",
	},

	"hemisphere-gauche": {
		Name:    "L'Hemisphere Gauche",
		Address: "221 Beaubien Est",
		Links: []string{
			"https://www.hemispheregauche.com/?lang=en",
		},
		AllowedDomains: []string{
			"hemispheregauche.com",
			"www.hemispheregauche.com",
		},
		Selector: HemisphereGaucheSelector,
		Website:  "https://www.hemispheregauche.com",
	},

	//"bar-le-ritz-pdb": {
	//	Name:    "Bar Le Ritz PDB",
	//	Address: "179 Rue Jean-Talon Ouest",
	//	Links: []string{
	//		"https://www.barleritzpdb.com/",
	//	},
	//	AllowedDomains: []string{
	//		"barleritzpdb.com",
	//		"www.barleritzpdb.com",
	//	},
	//	Selector: BarLeRitzPDBSelector,
	//	Website:  "https://www.barleritzpdb.com/",
	//},
}
