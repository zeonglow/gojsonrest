package jsonrest

import "log"

type IDtype int64

type Storable interface {
	New() Storable
	Show(id IDtype) error
	Create() error
	Delete(id IDtype) error
}

// Contact Root level object
type Contact struct {
	ID int64 `jsonapi:"primary,contact"`
	// people have multiple names too, and none, keeping it simple
	Name       string             `jsonapi:"attr,name"`
	Physicals  *[]AddressPhysical `jsonapi:"relation,address_physical"`
	Emails     *[]AddressEmail    `jsonapi:"relation,address_email"`
	Telephones *[]Telephone       `jsonapi:"relation,telephone"`
}

func (_ *Contact) New() Storable {
	return &Contact{}
}

func (data *Contact) Show(id IDtype) error {
	log.Printf("Show for '%T' called with %d", data, id)
	return nil
}

func (data *Contact) Create() error {
	log.Printf("Create for '%T' called", data)
	log.Printf("Trying to insert %v", data)
	log.Printf("Trying to Physicals %v", data.Physicals)

	return nil
}

func (data *Contact) Delete(id IDtype) error {
	log.Printf("Delete for '%T' called with %d", data, id)
	return nil
}

// AddressPhysical : many to one relationship with Contact
type AddressPhysical struct {
	ID int64 `jsonapi:"primary,address_physical"`
	// Summer house, work, home etc
	AddressName string `jsonapi:"attr,address_name"`

	Line1    string `jsonapi:"attr,line1"`
	Line2    string `jsonapi:"attr,line2"`
	Line3    string `jsonapi:"attr,line3"`
	City     string `jsonapi:"attr,city"`
	Postcode string `jsonapi:"attr,postcode"`
	Country  string `jsonapi:"attr,country"`
}

func (_ *AddressPhysical) New() Storable {
	return &AddressPhysical{}
}

func (data *AddressPhysical) Show(id IDtype) error {
	log.Printf("Show for '%T' called with %d", data, id)
	return nil
}

func (data *AddressPhysical) Create() error {
	log.Printf("Create for '%T' called", data)
	return nil
}

func (data *AddressPhysical) Delete(id IDtype) error {
	log.Printf("Delete for '%T' called with %d", data, id)
	return nil
}

// AddressEmail :  many to one relationship with Contact
type AddressEmail struct {
	ID          int64  `jsonapi:"primary,address_email"`
	AddressName string `jsonapi:"attr,address_name"`
	Email       string `jsonapi:"attr,email"`
}

func (_ *AddressEmail) New() Storable {
	return &AddressEmail{}
}

func (data *AddressEmail) Show(id IDtype) error {
	log.Printf("Show for '%T' called with %d", data, id)
	return nil
}

func (data *AddressEmail) Create() error {
	log.Printf("Create for '%T' called", data)
	return nil
}

func (data *AddressEmail) Delete(id IDtype) error {
	log.Printf("Delete for '%T' called with %d", data, id)
	return nil
}

// Telephone :  many to one relationship with Contact
type Telephone struct {
	ID int64 `jsonapi:"primary,telephone"`
	// type of telephone, work, home, mobile etc
	TelephoneName string `jsonapi:"attr,telephone_name"`
	TelephoneNum  string `jsonapi:"attr,telephone_num"`
}

func (_ *Telephone) New() Storable {
	return &Telephone{}
}

func (data *Telephone) Show(id IDtype) error {
	log.Printf("Show for '%T' called with %d", data, id)
	return nil
}

func (data *Telephone) Create() error {
	log.Printf("Create for '%T' called", data)
	log.Printf("Trying to insert %v", data)
	return nil
}

func (data *Telephone) Delete(id IDtype) error {
	log.Printf("Delete for '%T' called with %d", data, id)
	return nil
}
