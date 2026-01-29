package model

type Sector struct {
	IdSector       int    `json:"idSector"`
	SectorName     string `json:"sectorName"`
	SectorEditDate string `json:"sectorEditDate"`
	SectorEditUser int    `json:"sectorEditUser"`
	SecVMPL        string `json:"secVMPL"`
	PerformerId    int    `json:"performerId"`
	Dtact          string `json:"dtact"`
	TicketSize     string `json:"ticketSize"`
}
