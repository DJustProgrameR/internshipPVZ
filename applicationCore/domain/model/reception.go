package model

type ReceptionStatus string

const (
	ReceptionInProgress ReceptionStatus = "in_progress"
	ReceptionClosed     ReceptionStatus = "close"
)

type Reception struct {
	ID       string
	PVZID    string
	DateTime string
	Status   ReceptionStatus
}
