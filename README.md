# coraserver

## Requirements
- Azure Account
## Getting `clientID, clientSecret, tenant`
1. Sign into https://portal.azure.com

2. Create a new application and set the `callbackURI` to match what you are
   gonna use in production

3. Get the `clientID` from the overview page.

4. Generate a client secret: In the application's overview page, navigate to
   "Certificates & secrets" and click on "New client secret". Enter a
   description, choose an expiration option, and click "Add". Make a note of
   the generated *client secret value*. Please note that this will only appear
   once, if you refresh the page or login later this will be hidden.

5. If your organization allows you to access the Azure Active Directory then
   `tenant` will be your tenant ID but, if you are using a personal account you
   should set `tenant` to be `common` even though you will have a tenant ID.

## ```config.json```
The format of `config.json` is as follows
```json
{
  "clientID": "YOUR_CLIENT_ID",
  "clientSecret": "YOUR_SECRET_KEY",
  "redirectURL": "http://localhost:8080/oauth/callback",
  "scopes": [
    "openid"
  ],
  "tenant": "common"
}
```
## Building and Running
```bash
git clone https://github.com/deebakkarthi/coraserver
go mod tidy
go build
./coraserver
```
