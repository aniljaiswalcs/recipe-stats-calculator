package stats

// DeliveriesByPostcodeAndTime is used to count the number of deliveries to a postcode within given time.
type DeliveriesByPostcodeAndTime struct {
	Postcode string
	From     uint
	To       uint
}
