package db

const DBNAME = "hotel_reservation"

type Dropper interface {
	Drop() error
}

type Store struct {
	User  UserStore
	Org   OrgStore
	Event EventStore
}
