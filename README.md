# simple keycloak oidc client 

```
# run Keycloak

$ mkdir -p /tmp/kc
$ docker run -d --rm --name kc \
-p 8080:8080 \
-e KEYCLOAK_USER=admin \
-e KEYCLOAK_PASSWORD=admin \
-v /tmp/kc:/tmp \
jboss/keycloak:8.0.1


# configure Keycloak

1. Create a realm named 'demo'
2. Create a client named 'demo-client'
3. Configure the 'demo-client' to be confidential (Settings >> 'Access Type' to 'confidential') 
and use 'http://localhost:5050/demo/callback' as a 'Valid Redirect URIs'
4. Create a user 'demo' with password 'demo'. Make sure to activate and 'impersonate' for this user (set new password).
5. Check access, log in to 'http://localhost:8080/auth/realms/demo/account/'
6. Log in to 'http://localhost:8080/auth/' and get client secret (Clients >> demo-client >> Credentials >> Secret)


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