package db

const DBNAME = "hotel_reservation"

type Dropper interface {
	Drop() error
}
