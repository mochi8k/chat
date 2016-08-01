package main

type developer struct {
  name string `json:"name"`
  clientID string `json:"clientID"`
  clientSecret string `json:"clientSecret"`
}

func (d *developer) getName() string {
  return d.name
}

func (d *developer) getClientID() string {
  return d.clientID
}

func (d *developer) getClientSecret() string {
  return d.clientSecret
}
