package main

type Address struct {
	street      string
	city        string
	state       string
	postal_code string
	country     string
}

type user struct {
	id         int
	first_name string
	last_name  string
	email      string
	Address
}
