# Go IBAN service

## Setup instructions
- Clone this repository into a directory under $GOPATH/src
- Run `dep ensure` (https://github.com/golang/dep#installation)
- Run `make dev` to start the server

## Endpoints:
- `/` - States simple usage instructions
- `/validate` - Validates an IBAN if sent as POST field
- `/bban2iban` - Returns IBAN upon sending BBAN and country

## Options
**Host**

Override the default port to something else than `localhost`

**Port**

Override the default port to something else than `3000`

**Sanitize**

At the moment if you send hypens or any other non-alphanumeric characters (excluding space) the validation result will always be `invalid`. Unless this value is set to true, then all unwanted characters will be removed, and validation will be performed.

## Examples
**Validate Endpoint**
```
curl -d '{"IBAN": "NL44RABO0123456789"}' -H "Content-Type: application/json" -X POST http://localhost:3000/validate
```
**BBAN2IBAN Endpoint**
```
curl -d '{"BBAN":"RABO0123456789","COUNTRY":"NL"}' -H "Content-Type: application/json" -X POST http://localhost:3000/bban2iban
```