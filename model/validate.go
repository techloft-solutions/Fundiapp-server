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

func (l Location) Validate() error {
	_, err := govalidator.ValidateStruct(l)
	if err != nil {
		return err
	}
	return nil
}

func (c Category) Validate() error {
	_, err := govalidator.ValidateStruct(c)
	if err != nil {
		return err
	}
	return nil
}

func (r Review) Validate() error {
	_, err := govalidator.ValidateStruct(r)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) Validate() error {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		return err
	}
	return nil
}

func (b Bid) Validate() error {
	_, err := govalidator.ValidateStruct(b)
	if err != nil {
		return err
	}
	return nil
}

func (u User) Validate() error {
	_, err := govalidator.ValidateStruct(u)
	if err != nil {
		return err
	}
	return nil
}
