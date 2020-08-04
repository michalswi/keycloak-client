# simple keycloak oidc client 

The OIDC plugin needs three parameters to hook up with Keycloak (defined in `main.go`): 

- the client ID (clientID)
- the client secret (clientSecret)
- the discovery endpoint (keycloakURL)  

The discovery endpoint is needed to get information on where it can do authentication, token introspection, etc.

An example how to use Golang Keycloak Client (not OIDC) described in the [bottom](#-play-with-golang-keycloak-client). It's base on:  
https://github.com/Nerzal/gocloak

#### # run Keycloak

```
$ mkdir -p /tmp/kc
$ docker run -d --rm --name kc \
-p 8080:8080 \
-e KEYCLOAK_USER=admin \
-e KEYCLOAK_PASSWORD=admin \
-v /tmp/kc:/tmp \
jboss/keycloak:10.0.2
```

#### # configure Keycloak
```
- create a realm named 'demo'

- create a client named 'demo-client' with 'openid-connect' (it's by default)

- configure the 'demo-client' to be confidential (Settings >> 'Access Type' to 'confidential') 
and use 'http://localhost:5050/demo/callback' as a 'Valid Redirect URIs'

- create a user 'demo' with password 'demo'. Make sure to activate and 'impersonate' for this user (set new password)

- check access, log in to 'http://localhost:8080/auth/realms/demo/account/'

- log in to 'http://localhost:8080/auth/' and get client secret (Clients >> demo-client >> Credentials >> Secret)
```

#### # play with oidc
```
#1 
# run keycloak oidc client

$ SECRET=91da3bca-f4a0-4ead-baad-bc26e0b4298d
$ go run main.go $SECRET

# in a web browser connect to 'localhost:5050/home'
# it will redirect you to Keycloak
# provide 'demo/demo', it will redirect you to '/home'


#2 
# run POC app

$ SECRET=91da3bca-f4a0-4ead-baad-bc26e0b4298d
$ go run POC-go/gokey.go $SECRET

# in web browser connect to 'localhost:5050/demo'
# it will redirect you to Keycloak
# provide demo/demo, it will give you token

$ ACCESS_TOKEN=<>
$ curl -i -XGET -H "Authorization: Bearer $ACCESS_TOKEN" localhost:5050/demo
```

#### # play with Golang Keycloak Client

```
Add to configured Keycloak:

- new realm role 'demo-role'

- enable 'Service Accounts Enabled' for 'demo-client', once it's enabled new 'Service Account User' will be created (not visible in the Keycloak UI) with the name 'service-account-demo-client'

$ go run goclient/kclient.go
demo-client, 360e32f1-c0d4-4fb6-9179-2e70f5dfbb04
service-account-demo-client, 712cb93e-a091-4013-b3da-b9c84148d476
demo-role, b3735f63-1391-402d-a5f6-cb8b032dd84e
Realm role 'demo-role' added to user 'service-account-demo-client'

- once realm role is added, check in the Keycloak GUI (can't get that thru API for kc version <= 11.0.0), 
Clients >> select 'demo-client' >> Service Account Roles >> Realm Roles
```