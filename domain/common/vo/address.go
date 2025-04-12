package vo

// Address is a value object representing a physical address
type Address struct {
	ReceiverName string `json:"receiver_name"`
	Phone        string `json:"phone"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	District     string `json:"district"`
	City         string `json:"city"`
	Province     string `json:"province"`
	PostalCode   string `json:"postal_code"`
	IsDefault    bool   `json:"is_default"`
}

// Equals checks if two addresses are equal
func (a Address) Equals(other Address) bool {
	return a.ReceiverName == other.ReceiverName &&
		a.Phone == other.Phone &&
		a.AddressLine1 == other.AddressLine1 &&
		a.AddressLine2 == other.AddressLine2 &&
		a.District == other.District &&
		a.City == other.City &&
		a.Province == other.Province &&
		a.PostalCode == other.PostalCode
}
