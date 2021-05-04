package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/renatoaguimaraes/job-scheduler/internal/worker/proto"
	"github.com/renatoaguimaraes/job-scheduler/pkg/worker/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var config = conf.Config{ServerAddress: "localhost:8080", LogFolder: os.TempDir()}

func TestStartAuthnAuthzAdminUser(t *testing.T) {
	// creates server
	serv := createTestServer(t, clientca, servercert, serverkey)
	defer serv.Stop()
	// load client credentials
	clientcred, err := loadClientCredentials(serverca, admincert, adminkey)
	require.NoError(t, err)
	// connects to the server
	conn, err := grpc.Dial(config.ServerAddress, grpc.WithTransportCredentials(clientcred))
	require.NoError(t, err)
	// creates the client
	client := proto.NewWorkerServiceClient(conn)
	// calls admin function
	res, err := client.Start(context.Background(), &proto.StartRequest{Name: "ls"})
	require.NoError(t, err)
	assert.NotEmpty(t, res.JobID)
}

func TestStartAuthnAuthzUnauthorizedUser(t *testing.T) {
	// creates server
	serv := createTestServer(t, clientca, servercert, serverkey)
	defer serv.Stop()
	// load client credentials
	clientcred, err := loadClientCredentials(serverca, usercert, userkey)
	require.NoError(t, err)
	// connects to the server
	conn, err := grpc.Dial(config.ServerAddress, grpc.WithTransportCredentials(clientcred))
	require.NoError(t, err)
	// creates the client
	client := proto.NewWorkerServiceClient(conn)
	// calls admin function
	res, err := client.Start(context.Background(), &proto.StartRequest{Name: "ls"})
	assert.NotNil(t, err)
	stat, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, "unauthorized, user does not have privileges enough", stat.Message())
	assert.Nil(t, res)
}

// Utilities

func TestUntrustedUser(t *testing.T) {
	// creates server
	serv := createTestServer(t, clientca, servercert, serverkey)
	defer serv.Stop()
	// load client credentials
	clientcred, err := loadClientCredentials(serverca, unauthcert, unauthkey)
	require.NoError(t, err)
	// connects to the server
	conn, err := grpc.Dial(config.ServerAddress, grpc.WithTransportCredentials(clientcred))
	require.NoError(t, err)
	// creates the client
	client := proto.NewWorkerServiceClient(conn)
	// calls admin function
	res, err := client.Start(context.Background(), &proto.StartRequest{Name: "ls"})
	assert.NotNil(t, err)
	stat, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unavailable, stat.Code())
	assert.Equal(t, "connection closed", stat.Message())
	assert.Nil(t, res)
}

func createTestServer(t *testing.T, ca, cert, key []byte) *grpc.Server {
	// load server credentials
	servercred, err := loadServerCredentials(clientca, servercert, serverkey)
	require.NoError(t, err)
	// creates server
	serv, lis, err := createServer(config, servercred)
	require.NoError(t, err)
	// starts the server
	go func() {
		serv.Serve(lis)
	}()
	time.Sleep(time.Second)
	return serv
}

func loadServerCredentials(ca, cert, key []byte) (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(ca) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}
	serverCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}
	return credentials.NewTLS(config), nil
}

func loadClientCredentials(ca, cert, key []byte) (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(serverca) {
		return nil, errors.New("failed to add CA's certificate")
	}
	clientCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}
	return credentials.NewTLS(tlsConfig), nil
}

// server ca, certficate and key

var serverca = []byte(`-----BEGIN CERTIFICATE-----
MIIFoTCCA4mgAwIBAgIUfTom3hNDy2c2uBbnpdeYx/tiZ2AwDQYJKoZIhvcNAQEL
BQAwYDELMAkGA1UEBhMCQlIxEjAQBgNVBAsMCVNlcnZlciBDQTESMBAGA1UEAwwJ
bG9jYWxob3N0MSkwJwYJKoZIhvcNAQkBFhpyZW5hdG9hZ3VpbWFyYWVzQGdtYWls
LmNvbTAeFw0yMTA0MjcxNjMwMzJaFw0yMjA0MjcxNjMwMzJaMGAxCzAJBgNVBAYT
AkJSMRIwEAYDVQQLDAlTZXJ2ZXIgQ0ExEjAQBgNVBAMMCWxvY2FsaG9zdDEpMCcG
CSqGSIb3DQEJARYacmVuYXRvYWd1aW1hcmFlc0BnbWFpbC5jb20wggIiMA0GCSqG
SIb3DQEBAQUAA4ICDwAwggIKAoICAQDIOFu+OuF9s63XFlTAAKiYuUq+kzz3N4oJ
2tp6aIU55Tgw7eDyH512caGY0EOlmtc+R9caIJ4Oj54vEI+4dsM2ixbOmxpRJVSq
aJaI7W5d8Fc12DQAiF9mOkcIlXT1+9RxklXobdF9smo3rCFia+5gxKe3r1A6QNpS
+ZRZSfdnryEUBXVYFCGfLCuYEo79vxfaprilNXaF2enN5s/6XZLCx/1DEfXJFEhu
0IS0nML5LbzatYFLvG5hhA/cR8AfXiygtrYEhVzKd1kSpuxfUxr0PunpojLWOALe
m7hM2t206J8XAYk0AWK8gTlNLWDhy0WTmHZUC762pNNclrrf+ovHxvMJL3ZHJfSB
OFUpV6ZVimoNaO2APQhuwNnVt/b4KLBBZU4If1f5FNe/9U+WwLY1hITW4n5HYlch
RQrMmpJaOBPVuZA9T87SNRg7LUDXaMPIik1h0ki+lB/6YlROS1zLxDiN3tnbqSA4
jVMD9hN1JyOg5x6HgafG0Gv421PLw59PoqJVOlYAlCGRKBfavmt2BW2DoxLGBCq6
AwpANP2FFk0ZlblioTQeI8yTuuLlqRzwz2nUucASicr9ZOtrHTiS/3FFFwlcpM3L
4eASMSWh5O7puJY4l2xUQyWRHST77PmwC0/w02VPD14+CCE+ZfRlrKq55Ivnm/Yo
8368aSQ/YwIDAQABo1MwUTAdBgNVHQ4EFgQUlG+73tZFx84fCXNmV7yDO8H5RRYw
HwYDVR0jBBgwFoAUlG+73tZFx84fCXNmV7yDO8H5RRYwDwYDVR0TAQH/BAUwAwEB
/zANBgkqhkiG9w0BAQsFAAOCAgEAXDXcA12sENHbLrGErCquHtQ2tg0DuIgrE8Rv
0V8wDgLkurfkKqmHTl/e167d8XJq4lo9qlYOqOzpzMUuHgZ/3GVPBihinADVSSif
hrV74XwgpN1AJfG58SoT3UvoA9XhrWGxggPSzMlLhqYXWNphzj4nnM4c0wgmVudx
L1fg4ur6zXse2GpfdTNDJHJOsrGLLaMpmAQefnA5pj5hghGf+4qXANG3moEXZr1a
WXx+lo9JZucwZzzq+4qEy78LJJLSkrrMf+JOsl+6mfw/BMNKz0dkIc4B3dO1JEey
N8UpQjZCgplyja0cWDKLfEEpBaJJ/2AsAC/8jKMOe1DW3eIXWmXRR+dKlfbFBO4D
UBUHJmq0a2UjD1of4wrnXqyWWZ7A8tqQsnkIZmrwghMDu/Vifk0vGXfXdNVwL/YL
VWlMljcDJAWe4By03dkU14VhDZTLP/p3c/7uyq/SSpJyAxd+nnMiD6vXYe+P2STY
jRdDEHxmBNF+cAkihlNLHgSo0GuNUVJ+9YfQFhXlnQ1MgAvmcQOkItO557fezvj/
pATYfD3Bl/s0GzGMZ4IbyvNQL3JUJo4YZcr3CwKTmy69yKOQHiT94uTDEbrfGmth
1wz1eYHF41N2vdKaEj71piLamQ8MxPW+SmIFjvHEO0Lu5zAs9dhfn1Gi90KO+Dj6
+RfawT8=
-----END CERTIFICATE-----`)

var servercert = []byte(`-----BEGIN CERTIFICATE-----
MIIFYzCCA0ugAwIBAgIUZNjBXa1G9wiYJXlFwa6IoV5Zq+0wDQYJKoZIhvcNAQEL
BQAwYDELMAkGA1UEBhMCQlIxEjAQBgNVBAsMCVNlcnZlciBDQTESMBAGA1UEAwwJ
bG9jYWxob3N0MSkwJwYJKoZIhvcNAQkBFhpyZW5hdG9hZ3VpbWFyYWVzQGdtYWls
LmNvbTAeFw0yMTA0MjcxNjMwMzJaFw0yMTA2MjYxNjMwMzJaMF0xCzAJBgNVBAYT
AkJSMQ8wDQYDVQQLDAZTZXJ2ZXIxEjAQBgNVBAMMCWxvY2FsaG9zdDEpMCcGCSqG
SIb3DQEJARYacmVuYXRvYWd1aW1hcmFlc0BnbWFpbC5jb20wggIiMA0GCSqGSIb3
DQEBAQUAA4ICDwAwggIKAoICAQDGY5G8w6ot8qL0EeNeG1rFu8o0OejOSHE8jUov
ePPEHrRgqdfrv5xhHKZ24LYB0iQjT8PcY2cBYkdR0NNRbxtOBDggba4haUl4+Cw/
7mw2FYDjTk/fcunUQcQQB6sxFuYb7OIjfbw89tv2np6TLegPrQOVwRSKXvNFTD1v
hdqWssQHbyrZpnRzZBcTdO530u2ZhkRLA1uVrtScIRqFspOdn4WsxVgtqy3+BOPz
Y+cfYejJJZd76l96wZmgUdQcwg+Q1oiMDpE1p3hvO/JPGG6xYxu69ZJVvUHOjIzg
pEWWPh5ci8ZkE3DSC2r+XV9fMaiHy8AveOuqK9hevAoH56IR1TgPzaE4JIp3igPz
3bQoPtuDhdZpqcgFs4gXcPhg0BgLrw1twEsXiHjiSZ/dbSb9Rb/9ZFF+a0pe6zSd
3UCicvPPK6zu9y80oOvCAWa7bAQ6gW2/LklK3ontSRynL1TRvfCIemaPlTXJJ0hX
fj/Y16Kka+VxDaEo7yqTW3Oz6FPICjmqOM1leTpMXf50EWRKxMJLv6npfO05C/6+
6cxArEuFrdNDlMDCLPoq5Se6Mmohc92Huh2RRo6qMjLXI9LIdpah0PNVLh+xxA3f
UxPbkkDFYqvBNGdJ1D/odBYrZrR0AM3Pj1JVOABeObIJS83LVxnYBOiRjGXNqB+4
fvC6lQIDAQABoxgwFjAUBgNVHREEDTALgglsb2NhbGhvc3QwDQYJKoZIhvcNAQEL
BQADggIBAE6HT4YAzsT5uV7yp4BAauHvhvNFtPoP27hAa4B6ivZ4/0+L/FWRZNnY
jOndy+8KwgB6yqEUgrP208j5gVNXwegi/r33Btmvg2p3sFtvtVsD4jdoFk1BYfKl
T5P10aKfz0J44LJCLshNm5sSa+ief1M/pAr97FXbZoTdfO9T4gtLt2cDJSDiZKp8
xK+7huwPfimpOqSALncWgB4dtSOzSYUbf/cZaomsCrwvFB0i2sWKx2TKJ9ZU65F8
JG/jhv0jlKj6ZXNQi9S/MmSe3ZawazIOBmgv+HzmxVTRuzTIqnmiSlo5lfOYK6hQ
Fj76jB1k3ft6zQBaYAehh41MTPMWZtajxt1QqF+UEqOTdt6Fjwxr8fBJeAbwqARn
hMtULrmzF3GffkYdox1K3jLSr+ULcEJBLjIV02pjfchxCZgDyIUWyz462zhDnSvK
TCvOYHQjB7U1GQacAzVo5lsiMcTLWBX5IcJLmqa+TLEdK1HBh0sAYEJA4bVvFguc
E6KJsZyXquPBhP713+M1eClLCSlMDG3MDGD68SzHaVB0Nm46d58ItokJmDVxVuOR
aUi35dutkia4AF9xq/m+zTKUg8ejaOeKyPJA1+BUsaLZ6tQXzrp8YiGWlB2maiYB
mbDsZBwcghHEQ3Sf2v7SHKZG56aCHgIvMOhc5tvJn1H2gx+ysSWl
-----END CERTIFICATE-----`)

var serverkey = []byte(`-----BEGIN PRIVATE KEY-----
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
AoBGNvN7PhYoAgL0P0ejPaIm8mO42A==
-----END PRIVATE KEY-----`)

// client ca, certficate and key

var clientca = []byte(`-----BEGIN CERTIFICATE-----
MIIFoTCCA4mgAwIBAgIUb/csWgspk+6MTr6rngqHuAYyyhowDQYJKoZIhvcNAQEL
BQAwYDELMAkGA1UEBhMCQlIxEjAQBgNVBAsMCUNsaWVudCBDQTESMBAGA1UEAwwJ
bG9jYWxob3N0MSkwJwYJKoZIhvcNAQkBFhpyZW5hdG9hZ3VpbWFyYWVzQGdtYWls
LmNvbTAeFw0yMTA0MjcxNjMwMzJaFw0yMjA0MjcxNjMwMzJaMGAxCzAJBgNVBAYT
AkJSMRIwEAYDVQQLDAlDbGllbnQgQ0ExEjAQBgNVBAMMCWxvY2FsaG9zdDEpMCcG
CSqGSIb3DQEJARYacmVuYXRvYWd1aW1hcmFlc0BnbWFpbC5jb20wggIiMA0GCSqG
SIb3DQEBAQUAA4ICDwAwggIKAoICAQDvrEfB5XO55JUgaEjwo/dhOUMbLEy8ACpv
j0ZRWkgeAylW+0rAH1u5TdEejlkj5mQ2eR10DEQ4wK5wTgyT4fgDp7UGtcovMpOZ
8J1vX52HQnBkmDei2NUwhekJsgyHJFSWAzxLnZLyoIfil9xkbqmqZN54bifB0gZK
4C3PyE215Sq7szVUKmIp7tZ8+gico5QQgrJkdlHSHgGNlQlBz8oiUetI0w8xn1lO
LPQGeCz2Y7YQs/oH1xVDQeqi1Qgmep9X2It73tFY4+kwdLlwTq4hsuxlxpdOKBfp
cArX1zrZgFvfBtS54VyUl9kBexGFVMsDFWdlBY/93DpXT44NiY4/K/JT/OwKrX9d
AxgKpKUytLX55TsauI1nhvyozAi6+JFUPivCnTdgbV3nMuDjJLOK1pNbAcQTfIuL
bG6y6P+5x92kJ96bPgFWTkdfRL71ZObm+En/i4G/7wF6FZGIXa9cX+8WlyYx2cuR
d973gKRK65uDqmw8vHYYnifNFy3ltexxPFLeAA3H9DZrvMoPqQv/XWOeb3POB7ie
OkFmEfLyQ2NEHoeQAUdDB/vDb+xfF1CiMUuraSTdwSoPxgHPJsQvegSdN5avNPyj
mpgCCkw2AGEEpbubfOuEjB96yk0u3IdWSmv2F0KQslHGpAtlpQLUQTg8DIK6+HNM
vD4ieXGu4QIDAQABo1MwUTAdBgNVHQ4EFgQUj732vTo88plf0xL4O1rQyFUFXeow
HwYDVR0jBBgwFoAUj732vTo88plf0xL4O1rQyFUFXeowDwYDVR0TAQH/BAUwAwEB
/zANBgkqhkiG9w0BAQsFAAOCAgEAXRG3H4sriv6FlNJVJvEPssd94sRcAQ8tj+LD
DKeoIIfWL499hxtr0Ko1V1HJXhT2oB/E6Wh3LMxhM+fTRSmbi5d82ej8erDVv0+a
dLN6zuhdk2TmIsiBgJrIur7qyvpsf8qfxoU9U6YTYOI6Xnbx+9tk812Um+Kfph88
buFHxlewzzeLK+vIEbX5eMD+wUq7f5oB9iwT2KLWoO6qotslioHPsy/Wmmq0fnOp
QU5ZtdKjHB7ZUj0Esb1/fxixYicK/t9RyzbGcZZ8JJ+mbSxwk0Yi6+zGchFSKiCx
LIyU8fO9jr1sIQ7rZ5hmxx9ZH2acyjhtwU14bGKgf8fqAWCEUKNWlAfxaZLtK28h
layisir2pF76EHRqGrLrEtuJPezodeQiu+91y7GaiXHegar4h8vNjJtz8sBrBfu9
8SWgmi9Jcdz83dYtM7pyZ24VZrdRCx2Y2vHU0YOSNUwK96kMTM4ixt94pNzSouSo
MyZaMTJYliq/2cPG2AMVly3lmfkcl+tUQH5BdbOxpCfBQ6rZglYwq4biqTc/kGIM
8nMK5JLMc+uRhZU5Wec2fzgiovvr1nKSZHYkCvzih3qkn/5DKV0+YeD5Tk5xhpqp
bIwaKXK0q9UYWb4LYAGSBpmv9+DLJEjzOrdQQm5mBYLPPZTbA8x1mkck8wU9KZCf
0lbaB2A=
-----END CERTIFICATE-----`)

var admincert = []byte(`-----BEGIN CERTIFICATE-----
MIIFZjCCA06gAwIBAgIUVDESpCW7CAsDHuKMEWkO26dRwCMwDQYJKoZIhvcNAQEL
BQAwYDELMAkGA1UEBhMCQlIxEjAQBgNVBAsMCUNsaWVudCBDQTESMBAGA1UEAwwJ
bG9jYWxob3N0MSkwJwYJKoZIhvcNAQkBFhpyZW5hdG9hZ3VpbWFyYWVzQGdtYWls
LmNvbTAeFw0yMTA0MjcxNjMwMzNaFw0yMTA2MjYxNjMwMzNaMF0xCzAJBgNVBAYT
AkJSMQ8wDQYDVQQLDAZDbGllbnQxEjAQBgNVBAMMCWxvY2FsaG9zdDEpMCcGCSqG
SIb3DQEJARYacmVuYXRvYWd1aW1hcmFlc0BnbWFpbC5jb20wggIiMA0GCSqGSIb3
DQEBAQUAA4ICDwAwggIKAoICAQD7XbygX4jAFrDPjr97il6ZmKpHzDf6PkPFl6zq
r9Y91bbPryP+7o1k909H9DcOGs3kQc3Eraqkw5ppvQcjubfqOgtPvoxK6eVgLHzQ
SxlIcFEyNZuWYh77cikOod5KpJJYdwcO6jf6iq0+GEdX55X7U5y+FrHC//9yZQGw
6niH274ieZBPdg/sjJUA82vjsAhN3KOSW500U60eqNaDsCKMsKRpckBYnTvQkHQy
iKxvSD2VVaJDmfntFbmWWvgEbeKg8tv1yJ5efmcC8bLEFpnmGWElu183RVHTKTKx
1DIMiBzTv8DClZ4DYc0BFY4lPq8FOrkWpNI0nc1OOIGTbfF8vMzyF2MKPKBR/oHN
ZE5i4uZquZfGk7EJ0aujkxHZ5Iz5QW0twylWQMUs/l+R6GZqMkNwERXq0NGA3hvu
BYHixsGjHM081umWOHSvCn3En1nphWgR9tkPq3n92rXVFIEPiaCjatQJDD5VupFd
+OteiH/y1o8AzRKRN2VX7DElyxGy8AWq/OTQXeAMpgUTT3GRyRHPeqon3WVnNasq
qiWqXh82UPiQpsKqvyOUIIz7F9msgHKxmil48V4Q2mZ9fMZEP7LcUIO2iw6puY0q
JQk7jVGVYLNESrWUEENo5hXliWoT3Y+hmZwI6oB+d/iBXvwvl+wSjhUuf+q/CB5P
eRwgYQIDAQABoxswGTAXBgcqhkjOVggBBAwMCmFkbWluLHVzZXIwDQYJKoZIhvcN
AQELBQADggIBAAdyXYyrnMF2ZYd/rQXqvhZCi2/p+4O6Qn94onYPYQTMXC3R5P/C
nyd8BxR4lQtnYcLFvHQvwtgsL4dk+3gWIAbCN2ktg7sJg22vNN6qS70+xlPg6OWu
iMb6B/OuOAAUoo8kR47NZ/5a1nG0QNf9WIwCvwMyx4DsVDjFaT9at17i9obqHPtZ
BhPrkvfeNknRGC8IW/Z0A8zxYIz/MaiKuPBnJmWgQE+bkenIAYEfQgttJo7N9fHH
3KWNel2BMU6kaYsxgQiJUd0Geo5fL5VztC1oG00wieQreqvN4ZuH8hoJ8/aWjED6
z5RqnwzKdHE5qI/9AVUJbXOKFsWSWJcGx/VWeG2A5+p5ZRFH6N99gUWkct8dr2IX
fYvlIc326BQ1MHukVz+N8X1ztr97F7NlNtBAXzViITcn/zrKxaYEGjwOIkgMYgjx
jK7Pe+Ni9H3QpxMx8WALMu48mU1GKJ61GEOWt2uRanptiby2+9lO+npqkQ/M4ipR
IYP8O4Oi0Lons8Fr6KLKDbDOJRrMSR7e3FG3c02oZ2Nzk6833fpo+iIrdN16gOZB
0VRYGVTfb9GxNFS08pewNw3jKKuH4WLYUzA5EXHDbtf6kCiQo8O15pWf+VI7xh5i
ORSvoUyYq8SYbgcHkxN7XJ/StYE0g4uhn39+/fp8xjB3nsTLFlkQE86s
-----END CERTIFICATE-----`)

var adminkey = []byte(`-----BEGIN PRIVATE KEY-----
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
TYNTy30VTPdXVLBo+gTF/EEfIl5/RpHz
-----END PRIVATE KEY-----`)

var usercert = []byte(`-----BEGIN CERTIFICATE-----
MIIFYDCCA0igAwIBAgIUVDESpCW7CAsDHuKMEWkO26dRwCQwDQYJKoZIhvcNAQEL
BQAwYDELMAkGA1UEBhMCQlIxEjAQBgNVBAsMCUNsaWVudCBDQTESMBAGA1UEAwwJ
bG9jYWxob3N0MSkwJwYJKoZIhvcNAQkBFhpyZW5hdG9hZ3VpbWFyYWVzQGdtYWls
LmNvbTAeFw0yMTA1MDQwMjM4NDRaFw0yMTA3MDMwMjM4NDRaMF0xCzAJBgNVBAYT
AkJSMQ8wDQYDVQQLDAZDbGllbnQxEjAQBgNVBAMMCWxvY2FsaG9zdDEpMCcGCSqG
SIb3DQEJARYacmVuYXRvYWd1aW1hcmFlc0BnbWFpbC5jb20wggIiMA0GCSqGSIb3
DQEBAQUAA4ICDwAwggIKAoICAQDdn8QmYzMZOj31gZR8UO4d91T325Ax9pqOuNuj
6Ep5osZ6bk9gEN5dOsh6P6PDSI4HV+IwpbAj9AiID9mRKLyRDb+Z0FqEpLTZNCDN
ua1zqaz7Gm5IyWGt1t/zL/p6RlqAkKm+lX+waY19LbL73F88cB6748YYdP4pt3vm
olduQt88b57OYhEmGfODp8ub/uK5lsS0k5Dv2+I/6InyWEuzpsHpzkF6VB2ilAbp
iinko3EibPq4FQojanDGc+cICdXtCD/z3hBAYqaMyzdPtZt/EEcRcCrhkFQCvEZ0
yvjlvxt1oBx8TksmTQwiusPG5Mk2rVXDdKXT9bH7og4iECzv5LjbSZBiJsmQ46+p
k4rlUoCH30yElpKresM+DPot8GaRnRfljebRoWm5WxnK3Bfxk9pA6rojURd5Mp7o
kJnwzoXSkfwdz7gSL206Avp3+nwefaAHgNPbDDIcB4hIgelowuHTYJAb5e1OO0zd
r5PX6vf/ERhFgVS1/iTbCUNAzXJXL1gK7czpOXN/tXdJbUK5Ngiw3Wt8bUUTjuqf
pf7MQCI6nIkfuTgafDUKd3HEEIHgRELaQShZyRX3xLYRHtS8sGsMALVFNR7zdqR3
h9zgj55oSWP/Xzc++Cy7ahifhV0IaSgQCW1WyFE/2FCcgKqLYOniMrWE6O0W8y1K
lDwEfwIDAQABoxUwEzARBgcqhkjOVggBBAYMBHVzZXIwDQYJKoZIhvcNAQELBQAD
ggIBALAc9zter/JG6dsAfJH3kfIQQrocOf9MMKaAsohB3jHq7qZb0sBRUSlbrypr
BG9zilvcloeQE0JSKPhoL/348OauNl3ABI3ekpmjZiqitKG7sZcTd0L2M94M/RFX
pTHeBmJLvij4EsP/xhgowP4mkvXjM59j3aQUMcWSsLq1eNGkadJRKfie/gTxCU6W
HvEpGMS6hrHWKRPDqow8DROhvsDeNPXPxHzlD3dl+xpOdNQlKqF5nXR8LvCmUEJ0
jLyHi7LwbCC9a5aIbnJkR9L9paT2/MSl9lNbMkSWtzb6k43e85CdGhRb0zcJpgIO
t5UDOty+Q6rMZmqMeSl5yhFfFmdBoS7WV6QG/pdoSZO8O7Rs5Mllck90FBb9VNG+
Om3bco7LaitCIyXN4nsT3E+U3h+Z8nSaZCkZDU1xru+VC1DlzBgE9bOpz1RjWMYU
bA6GjAJB5kJZDzttDsVyeYP56CTQmIn895ToLx8VjEha7IjOM2WudaL3cX8PGpAA
8Opeis5VgvHKER7k7onxazcLomDfBTrOnsqVnxVbGxWStAGwA8qn6yKAeIbfrLYk
fwypaFcWktjH93WX5LraH4Id8Vhi7erALyyykaPZuO8fDJ7if9YYqsaNpm6GAzR4
XO2jvvVbhzKDA5OuOVpvnzUXmcO54gCgcrpHca9317thR/mh
-----END CERTIFICATE-----`)

var userkey = []byte(`-----BEGIN PRIVATE KEY-----
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
vauhSrsJEsHZbLbBbNOnkRljL1IiiZ0=
-----END PRIVATE KEY-----`)

var unauthcert = []byte(`-----BEGIN CERTIFICATE-----
MIIFczCCA1ugAwIBAgIUFhOYWF6fyuq+HABQlUZIoINIWaQwDQYJKoZIhvcNAQEL
BQAwZjELMAkGA1UEBhMCQlIxGDAWBgNVBAsMD1VuYXV0aG9yaXplZCBDQTESMBAG
A1UEAwwJbG9jYWxob3N0MSkwJwYJKoZIhvcNAQkBFhpyZW5hdG9hZ3VpbWFyYWVz
QGdtYWlsLmNvbTAeFw0yMTA1MDQwMzEzNDBaFw0yMTA3MDMwMzEzNDBaMGoxCzAJ
BgNVBAYTAkJSMRwwGgYDVQQLDBNVbmF1dGhvcml6ZWQgQ2xpZW50MRIwEAYDVQQD
DAlsb2NhbGhvc3QxKTAnBgkqhkiG9w0BCQEWGnJlbmF0b2FndWltYXJhZXNAZ21h
aWwuY29tMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAq+4h5jmnadge
dKCxW/ImVszI9ojvAXWxWUJwIp0ubEFk90L/cLijbF6b8c69hvFdhbJXVG4sXqw2
royzfpIzLGGTEyl3fKoYmGDuojeHDWhmpGcM2sN8Pm9bIEJ53YhBeUiKFNNkL5Fp
POdy8NTIhmg2yi4tEfu8L5xOn8M9FeApLVP29INY+v1T6VIXPL+khGBqt2YQPSwG
e0y7YkUcxYKsrWmwnz8cnwOSA9PzeuNL5AcihWzX1JgXzvJ+ClWHtHLCYSkmkNa/
QKtxC371SEANYDTixqy+KmhIVWPIjXaN4oLN+Y29vYzh0GhA8feJ4Tpy7Gmk8YNN
g/Y1i4Ds295AJ/rLuZ3CrR1/WgRL9Tmfi3TkeFcB8dDfuLhMPDJdL+2KLKe753ok
LFzhq1GEx+ARvnISHwLW9Av0va1EDAFmXSu6fqMgyT7w0Rc9ElyvXP8F0HSUzB5a
OvV3V2QKjkqC3C3bBQzl9axu+eYjSFaJODSbwBxL9eZdrycZwti+VkSqubj7xg4r
3QYcKqleRvWDTqj+l9LJEP0dqQJLOJIQdo++m+/HYQd5xfSL6+XpCOJN58Q71DYH
RZZ3gzDQxCr2WILyYaXk1XRA7u0dg8vk3fsic/AELdxx73mbNjyiQlQX1JFJHhJ8
pK5dMmJuhivGbv76Y2Voz3AIcEL5pGcCAwEAAaMVMBMwEQYHKoZIzlYIAQQGDAR1
c2VyMA0GCSqGSIb3DQEBCwUAA4ICAQCdoDoJI1ml9kYYG8L019Ps7FC8DAcENS0s
eMG1ck9+YPp5QRmLkV+242dQgJLkK+bEASK2RE00uFpaEVfCu3GpFuIo65nL+UPI
BbgOejYhUOwZ0zNMMYoNW1R05kEzzJ3wVFw9FWH2kAP3hsYwKuthuhA6ryGCEUPJ
NTeoDDSCFNlQ/KTgX1r888MOVfig70XW8ODxy/8iqD8IesNQdDX59mCT8gj+ikYS
oqG3ZbQcuMNEv7OzD1CapleIQ7N4BotZq5nSNDdLpahUbKWwVX8+fh0nP5qLbWVu
o+WBrqV1Pen+iq4bL+Ml28A043S/8/diju8fo2xYC+x9ayepi1F2E2QM9bPMz7iT
oeGa6nVv2A0n1Bz9baVHJp555058qQGI5FWICryUBPDXDpeYQ+3B8pF6SAbZ2VD0
aR95ndjxtzZkQDW204TrHGcPKxbxZJK8kZfGG9mk9KuKLriNKltaLCJiBKbpH7de
2RLIttRrH5Olu16MIu+U0lzbSjSHwb/qweZ0KQ3POyH6OTfmnCL2P0agQJdUeHN3
eVMcxGrYjabFsGAASk7Hwp93gfb39nMMIv/MvZ69A8W2LF4RrlahjKbyC4GVe8jC
qPJ3P3NlUyf2RRV1wQAIn72Na+XaQlGT3hRBJVXkUKsFvFFuOJy8bny99t9XHTZK
VcFOvd7irw==
-----END CERTIFICATE-----`)

var unauthkey = []byte(`-----BEGIN PRIVATE KEY-----
MIIJQgIBADANBgkqhkiG9w0BAQEFAASCCSwwggkoAgEAAoICAQCr7iHmOadp2B50
oLFb8iZWzMj2iO8BdbFZQnAinS5sQWT3Qv9wuKNsXpvxzr2G8V2FsldUbixerDau
jLN+kjMsYZMTKXd8qhiYYO6iN4cNaGakZwzaw3w+b1sgQnndiEF5SIoU02QvkWk8
53Lw1MiGaDbKLi0R+7wvnE6fwz0V4CktU/b0g1j6/VPpUhc8v6SEYGq3ZhA9LAZ7
TLtiRRzFgqytabCfPxyfA5ID0/N640vkByKFbNfUmBfO8n4KVYe0csJhKSaQ1r9A
q3ELfvVIQA1gNOLGrL4qaEhVY8iNdo3igs35jb29jOHQaEDx94nhOnLsaaTxg02D
9jWLgOzb3kAn+su5ncKtHX9aBEv1OZ+LdOR4VwHx0N+4uEw8Ml0v7Yosp7vneiQs
XOGrUYTH4BG+chIfAtb0C/S9rUQMAWZdK7p+oyDJPvDRFz0SXK9c/wXQdJTMHlo6
9XdXZAqOSoLcLdsFDOX1rG755iNIVok4NJvAHEv15l2vJxnC2L5WRKq5uPvGDivd
BhwqqV5G9YNOqP6X0skQ/R2pAks4khB2j76b78dhB3nF9Ivr5ekI4k3nxDvUNgdF
lneDMNDEKvZYgvJhpeTVdEDu7R2Dy+Td+yJz8AQt3HHveZs2PKJCVBfUkUkeEnyk
rl0yYm6GK8Zu/vpjZWjPcAhwQvmkZwIDAQABAoICAFtu//406RCEC+ZQUyInzDXb
IIDj399x3MgwafwuhTKzMGPC8J/cwaRvSBW3sdli4S6p4oGXOed7RHVdhFOqoqZW
meXV5qKYvw3CdwYz65G41OVXnGF3FssJY3frgm3K+a7rYeujZCNb/JCUMe9b6ex4
3bJ2DigDLVhQkwupxAGvPZbKkYoFlDuen5J0wsDe1jpEIFy6XaZ6lFPcQccIJD3x
ua1biL+Xy7vRJUT94f2XhZOWK8izUrDP1dGL6nXKeKoBYdhUDxt16S4WPr8zXHR5
A+QHHBc2ZMeev+EReOWqh9hPYgT2WVLF6o1v9CH6WAw8jkOmBz8dagrk1CMFhHCm
Mdcjym9XGyb3dSEeD8RjFin/D8jtEyZVq+SXoKkWV71G3HHkB2OS9vcjc7YpmZ65
bLgViZFqjppmyW1gOuovolTJ0OtSUJZEy4F4WOZAuY1XAdcLLe4YR4FzMW+DSPYc
ofeEn3P/0W7jLiOry3sPFn40Hfiw6pNJ8dkHDBLTSMnT5Po/crPVuJecp82do97w
zJpsbGGZMniRPSA9v6RyU+VuufTwJq8TZ6g12+/zmFjrZeEuWjX0/PTW+6PGx5vF
yI+HkPb1PBuJbIHLZtwpChJjnigam6mvBZ0b8dI10bfbRyGoiCZPpOW4ht76IJP9
GuEgMTPHvo+WH4D3CGrJAoIBAQDXFayAxtWv0/d5BdJM3Lt9weYjsidz1rnANqaI
g6NCy8z2ezFxNT7URVQ6malEzTx5b3xpaO9bR+WHY9Cx+wsO/mhVyA247FyVPwwl
i6FZ3D+ieKltMGRmnTIBOZxFPC3IxtqBjahvptXbDF8l5GhO+dLNpCK6b7AFWM+E
MnqKIfPJ2XlBuZFpQrOH8eisl/AbL/soyA2MSxD9Y8NP13IVhkgIY1Oi9B07fFkF
YeoA/xSQtl0WHC6Y2CNxjyRONbKWExwKMXT8J/12+qHPi3Fz48njOAWrJttlRArN
yiv1odDW7N4U5A3iZt8YdlA1PubbxPLOkq2S+TOAwd+FGizbAoIBAQDMouVhqm8d
oq/V2Rft9VCk94tKrPjz07iCwhJq+eZpiuGvsBlCKn60YWkfsme2HpKGfmcQ7T06
OA+s7cvhK8vfz85gIhpG+1sXKv+y7R3Dl4IpfJa7yfbdk34tYxrq6qkBmRJsSD63
sXfaKZP1CpQt8vjSav5ag2rX3bHgfd2yeIhCE+DlKgugNs/50T9XA3YSc/lCIET/
x7IKMBlXdik8FWzcwW8F2hvcVbeFMdp4DKgG08oj848kzioGpGhqhG1KQjTeWQ92
e5sb5WrTxsxaS1AIMXIPkehQ4wucaC/sftkfGUiXocd+ml7xB1SiwWpA63wX021G
RgpiwzYxu3ZlAoIBAQC02RuH2CgaxI7j6rouOLdJgX0B7K1xoE3lSc5zIMRoyLYZ
VfJ3rv6aO4UcFNIX8L31mYExnLMNvIFJeusii+R7gfy2jBTFtbAPKckL1MEhrqiu
8uf60GLLSUefchJswLH6jQFzR2omH1DX8yoZ0VjHdxYCJQ2yV3DYvhkWnd8dQHkc
8sgbNG4I5LXSC9zJKqQQKCL14mhJ22B4vqad/piFcrgknWfYr4vY1bCAbxj0J4bi
OwRAMAgaKjp7JJGxVUNJHt5Hz9f+oZ1kqk7eFPCbMiAx9owFohF9r12qNWELEzln
ThlZ1Xu7LyZNCkZczvwMNCi4+uoJl5HvpJRN9XlbAoIBADkI5dnUwOeHom599Xdz
Oijgfcgwcaqzxedb4/pA8IFWHhTzhfa1FV99DumwtctCtiAaNuu206vYWDYgiQSX
Sll099Y/aweBox/P8jiScgtDvRmHChQI9G9JXo+T4bq78KLrYQEhGanlIryBfiV7
71TnNYagH4hmvG6x8ZPaQOIvSfrww2vbziW3YTwFoBvGcTAvdreBevm3VN4WDdoc
qt1+MztMBn+hDMbadS4AeR2gmWxdtydSCQF/HKOBS0D06+kYjteyBudFCDQ8OwwP
ioFYIFRIQ7wnNJLm1SOgvkqyCB3s1Bi/FHUq3W9cVbCPK+gwOgQB/6DogJiKRB6U
WykCggEAaNdr8q+J0MGaaKzL2YHizo0T3vRSaIIa2nxafNTg3ZqCIcIf1BwX26+v
TIEwiJCC/l9v3+yJonZG1hXLV4EIqxFWvSc6fwXTvqNSYmSK1/8h3tBuH+FeFCgT
5iT+KPVgk6QvUPpNW6OmjWrUGeJYrV3NbVcZjI7lffR6p1NMVjyG1WxrRSDcjm5+
XevVvUOzRvdwjD6Qf0YHXEZ08XJDrrK6H0FsCxbkPojKYf8y5IGF2k1Ajk3++QVR
BwLfqar8p93ODkgVmd4TrdQh0X1OItPz07EY7HsEYL1eZ/ZQd2MVHXdCumUntCGe
lDB2NY6XhiC7TZgpR2y07TPkioq1kA==
-----END PRIVATE KEY-----`)
