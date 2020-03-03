## keycloak client 

main App: `keyclient.go`  
POC App: `POC-go/gokey.go`
```
# run Keycloak

$ docker run -d -p 8080:8080 --rm \
-e KEYCLOAK_USER=admin \
-e KEYCLOAK_PASSWORD=admin \
-v $(pwd):/tmp \
--name kc \
jboss/keycloak:8.0.1


# configure Keycloak

1. Create a realm named 'demo'
2. Create a client named 'demo-client'
3. Configure the 'demo-client' to be confidential (Settings >> 'Access Type' to confidential) 
and use 'http://localhost:5050/demo/callback' as a 'Valid Redirect URIs'
4. Create a user 'demo' with password 'demo'. Make sure to activate and 'impersonate' for this user.
5. Check access, log in: 'http://localhost:8080/auth/realms/demo/account/'
6. Get client secret (Clients >> demo-client >> Credentials >> Secret)


# run main App

$ go run keyclient.go

# in web browser connect to 'localhost:5050/home'
# it will redirect you to Keycloak
# provide demo/demo, it will redirect you to '/home'


# run POC App

$ go run gokey.go

# in web browser connect to 'localhost:5050/demo'
# it will redirect you to Keycloak
# provide demo/demo, it will give you token

$ TOKEN="..."

$ curl -i -XGET -H "Authorization: Bearer $TOKEN" localhost:5050/demo
```