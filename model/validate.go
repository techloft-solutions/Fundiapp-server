package model

import "github.com/asaskevich/govalidator"

func (b Booking) Validate() error {
	_, err := govalidator.ValidateStruct(b)
	if err != nil {
		return err
	}
	return nil
}

func (r Request) Validate() error {
	_, err := govalidator.ValidateStruct(r)
	if err != nil {
		return err
	}
	return nil
}

func (p Profile) Validate() error {
	_, err := govalidator.ValidateStruct(p)
	if err != nil {
		return err
	}
	return nil
}
