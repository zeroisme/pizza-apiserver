SHELL = bash

OUTPUT := serving-cert-secret.yaml v1alpha1-apiservice.yaml v1beta1-apiservice.yaml

all: tls.key tls.crt $(OUTPUT)

tls.key tls.crt:
		#openssl req -new -x509 -subj "/CN=api.pizza-apiserver.svc"  -nodes -newkey rsa:4096 -keyout tls.key -out tls.crt -days 365
		openssl req -new -x509 -nodes -newkey rsa:4096 -keyout tls.key -out tls.crt -days 365 -config minimal.cnf 

$(OUTPUT): tls.key tls.crt
		sed 's,CERT,$(shell base64 -w0 tls.crt),;s,KEY,$(shell base64 -w0 tls.key),' $@.template > $@

clean:
		rm -f tls.key tls.crt $(OUTPUT)
