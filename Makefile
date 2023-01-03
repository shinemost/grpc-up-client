# 2.生成SAN的服务端证书 生成服务端私钥（serve.key）–>生成服务端证书请求（server.csr）–>CA对服务端请求文件签名，生成服务端证书（server.pem）
# 2.1 生成服务端证书私钥
# 2.2 根据私钥server.key生成证书请求文件server.csr
# 2.3 请求CA对证书请求文件签名，生成最终证书文件
genclientcerts:
	rm -rf certs/client.*
	openssl genrsa -out certs/clinet.key 2048 
	openssl req -new -nodes -key certs/clinet.key -out certs/clinet.csr -subj "/C=cn/OU=myorg/O=mytest/CN=localhost" -config ./openssl.cnf -extensions v3_req 
	openssl x509 -req -days 365 -in certs/clinet.csr -out certs/clinet.pem -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -extfile ./openssl.cnf -extensions v3_req 

mock:
	rm -rf mocks/*
	mockgen -package mocks -destination mocks/prodinfo_mock.go github.com/shinemost/grpc-up-client/pbs ProductInfoClient

.PHONY: genclientcerts mock