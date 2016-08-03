package main

type provider struct {
	Name         string `json:"name"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}

func (d *provider) getName() string {
	return d.Name
}

func (d *provider) getClientID() string {
	return d.ClientID
}

func (d *provider) getClientSecret() string {
	return d.ClientSecret
}
