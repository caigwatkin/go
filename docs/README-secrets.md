# Secrets

## Token keys

Name the token:

```bash
export TOKEN_NAME="access_token"
```

Generate keys:

```bash
openssl req -nodes -x509 -newkey rsa:4096 -keyout ${TOKEN_NAME}_private_key.pem -out ${TOKEN_NAME}_public_key.pem
```

Optionally confirm keys:

```bash
openssl x509 -in ${TOKEN_NAME}_private_key.pem -text
openssl rsa -in ${TOKEN_NAME}_public_key.pem -check
```

Store in base64:

```bash
cat ${TOKEN_NAME}_private_key.pem | base64 > ${TOKEN_NAME}_private_key.txt
cat ${TOKEN_NAME}_public_key.pem | base64 > ${TOKEN_NAME}_public_key.txt
```
