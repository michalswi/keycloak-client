package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Nerzal/gocloak/v6"
)

// https://godoc.org/gopkg.in/nerzal/gocloak.v6

func main() {
	gclient := gocloak.NewClient("http://localhost:8080/")
	ctx := context.Background()

	token, err := gclient.LoginAdmin(ctx, "admin", "admin", "master")
	if err != nil {
		log.Println(err)
	}

	// clients
	// curl -s http://localhost:8080/auth/admin/realms/demo/clients -H "Authorization: Bearer $TOKEN" | jq '.[] | select(.clientId=="demo-client").id' --raw-output
	clientid := "demo-client"
	view := true
	clients := gocloak.GetClientsParams{
		ClientID:     &clientid,
		ViewableOnly: &view,
	}
	c, err := gclient.GetClients(ctx, token.AccessToken, "demo", clients)
	if err != nil {
		log.Println(err)
	}
	// fmt.Printf("%+v\n", *c[0])
	fmt.Printf("%v, %v\n", *c[0].ClientID, *c[0].ID)

	// service account user
	// curl -s http://localhost:8080/auth/admin/realms/demo/clients/$CLIENTID/service-account-user -H "Authorization: Bearer $TOKEN" | jq
	sa, err := gclient.GetClientServiceAccount(ctx, token.AccessToken, "demo", *c[0].ID)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%v, %v\n", *sa.Username, *sa.ID)

	// roles
	// curl -s http://localhost:8080/auth/admin/realms/demo/roles -H "Authorization: Bearer $TOKEN" | jq '.[] | select(.name=="demo-role").id' --raw-output
	roleName := "demo-role"
	ro, err := gclient.GetRealmRole(ctx, token.AccessToken, "demo", roleName)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%v, %v\n", *ro.Name, *ro.ID)

	// add realm role to service account user (not to 'demo-client' directly)
	sarolename := *ro.Name
	saroleid := *ro.ID
	roles := gocloak.Role{
		Name: &sarolename,
		ID:   &saroleid,
	}

	err = gclient.AddRealmRoleToUser(ctx, token.AccessToken, "demo", *sa.ID, []gocloak.Role{roles})
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("Realm role '%s' added to user '%s'\n", sarolename, *sa.Username)

	// once realm role is added you have to check in the Keycloak GUI (can't get that thru API for kc version <= 10.0.2)
	// Clients >> select 'demo-client' >> Service Account Roles
}
