package main

// venue websites' CSS selectors for events
const (
	CasaDelPopoloSelector = "div.md\\:w-5\\/12.p-6"
	QuaiDesBrumesSelector = "article.mec-event-article"
	CafeCampusSelector    = "div.noo-shevent-content"
)

var allVenues = map[string]Venue{

	/*
		DEBUT CASA & co
	*/

	"casa-del-popolo": {
		Name:  "Casa del Popolo",
		Group: "casa",
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
		Name:  "La Sala Rossa",
		Group: "casa",
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
		Name:  "La Sotterenea",
		Group: "casa",
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
		Name:  "Ptit Ours",
		Group: "casa",
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
		Name:  "La Toscadura",
		Group: "casa",
		Links: []string{
			"https://casadelpopolo.com/en/events/la-toscadura",
		},
		AllowedDomains: []string{
			"casadelpopolo.com",
			"www.casadelpopolo.com",
		},
		Selector: CasaDelPopoloSelector,
	},

	/*
		FIN CASA & co
	*/

	// ----------------------------------------------------------------------------------------------------

	"quai-des-brumes": {
		Name:    "Quai des Brumes",
		Address: "4481 Rue Saint-Denis, Montréal",
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
		Address: "57 Rue Prince-Arthur Est, Montréal",
		Links: []string{
			"https://spectacles.cafecampus.com/evenements",
		},
		AllowedDomains: []string{
			"spectacles.cafecampus.com",
			"www.spectacles.cafecampus.com",
			//"cafecampus.com",
			//"www.cafecampus.com",
		},
		Selector: CafeCampusSelector,
		Website:  "https://www.cafecampus.com/",
	},
}
