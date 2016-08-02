package main

type developer struct {
	Name         string `json:"name"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

func (d *developer) getName() string {
	return d.Name
}

func (d *developer) getClientID() string {
	return d.ClientID
}

func (d *developer) getClientSecret() string {
	return d.ClientSecret
}
