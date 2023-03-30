package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"

	"gopkg.in/mgo.v2"
)

var (
	connPool map[string]*mgo.Session
)

const rootPEM = `-----BEGIN CERTIFICATE-----
MIIECTCCAvGgAwIBAgICEAAwDQYJKoZIhvcNAQELBQAwgZUxCzAJBgNVBAYTAlVT
MRAwDgYDVQQHDAdTZWF0dGxlMRMwEQYDVQQIDApXYXNoaW5ndG9uMSIwIAYDVQQK
DBlBbWF6b24gV2ViIFNlcnZpY2VzLCBJbmMuMRMwEQYDVQQLDApBbWF6b24gUkRT
MSYwJAYDVQQDDB1BbWF6b24gUkRTIGV1LXNvdXRoLTEgUm9vdCBDQTAeFw0xOTEw
MzAyMDIxMzBaFw0yNDEwMzAyMDIxMzBaMIGQMQswCQYDVQQGEwJVUzETMBEGA1UE
CAwKV2FzaGluZ3RvbjEQMA4GA1UEBwwHU2VhdHRsZTEiMCAGA1UECgwZQW1hem9u
IFdlYiBTZXJ2aWNlcywgSW5jLjETMBEGA1UECwwKQW1hem9uIFJEUzEhMB8GA1UE
AwwYQW1hem9uIFJEUyBldS1zb3V0aC0xIENBMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAtEyjYcajx6xImJn8Vz1zjdmL4ANPgQXwF7+tF7xccmNAZETb
bzb3I9i5fZlmrRaVznX+9biXVaGxYzIUIR3huQ3Q283KsDYnVuGa3mk690vhvJbB
QIPgKa5mVwJppnuJm78KqaSpi0vxyCPe3h8h6LLFawVyWrYNZ4okli1/U582eef8
RzJp/Ear3KgHOLIiCdPDF0rjOdCG1MOlDLixVnPn9IYOciqO+VivXBg+jtfc5J+L
AaPm0/Yx4uELt1tkbWkm4BvTU/gBOODnYziITZM0l6Fgwvbwgq5duAtKW+h031lC
37rEvrclqcp4wrsUYcLAWX79ZyKIlRxcAdvEhQIDAQABo2YwZDAOBgNVHQ8BAf8E
BAMCAQYwEgYDVR0TAQH/BAgwBgEB/wIBADAdBgNVHQ4EFgQU7zPyc0azQxnBCe7D
b9KAadH1QSEwHwYDVR0jBBgwFoAUFBAFcgJe/BBuZiGeZ8STfpkgRYQwDQYJKoZI
hvcNAQELBQADggEBAFGaNiYxg7yC/xauXPlaqLCtwbm2dKyK9nIFbF/7be8mk7Q3
MOA0of1vGHPLVQLr6bJJpD9MAbUcm4cPAwWaxwcNpxOjYOFDaq10PCK4eRAxZWwF
NJRIRmGsl8NEsMNTMCy8X+Kyw5EzH4vWFl5Uf2bGKOeFg0zt43jWQVOX6C+aL3Cd
pRS5MhmYpxMG8irrNOxf4NVFE2zpJOCm3bn0STLhkDcV/ww4zMzObTJhiIb5wSWn
EXKKWhUXuRt7A2y1KJtXpTbSRHQxE++69Go1tWhXtRiULCJtf7wF2Ksm0RR/AdXT
1uR1vKyH5KBJPX3ppYkQDukoHTFR0CpB+G84NLo=
-----END CERTIFICATE-----`

type DbConf struct {
	Host   string
	Port   string
	User   string
	Pass   string
	CaPath string
}
type MongoItem struct {
	Id            string `bson:"_id"`
	Category_name string `bson:"category_name"`
	Create_time   int32  `bson:"create_time"`
	En_txt        string `bson:"en_txt"`
	En_type       int8   `bson:"en_type"`
	French_txt    string
	From          string `bson:"from"`
	German_txt    string
	Img_path      string
	Remark        string
	Word_limit    int32
	Zh_txt        string `bson:"zh_txt"`
	Zh_type       int8   `bson:"zh_type"`
}

type DbConnConf struct {
	Confs []*DbConf
}

func (conf *DbConnConf) Init() {
	connPool = make(map[string]*mgo.Session, 20)
}

func (conf *DbConnConf) GetConn(ctx context.Context) *mgo.Session {
	dbConf := conf.getConnConf(ctx)
	confIdent := conf.getConnIdent(dbConf)
	var conn *mgo.Session
	if _, ok := connPool[confIdent]; !ok {
		var err error
		var url = "mongodb://" + dbConf.User + ":" + dbConf.Pass + "@" + dbConf.Host + ":" + dbConf.Port
		fmt.Println(`mongo url:`, url)
		conn, err = mgo.Dial(url)
		if err != nil {
			fmt.Println("err:", err)
		}
		connPool[confIdent] = conn
	} else {
		conn = connPool[confIdent]
	}
	return conn.Clone()
}
func (conf *DbConnConf) GetSSLCon(ctx context.Context) *mgo.Session {
	dbConf := conf.getConnConf(ctx)
	rootCerts := x509.NewCertPool()
	ok := rootCerts.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		panic("failed to parse root certificat")
	}
	tlsConfig := &tls.Config{
		RootCAs:            rootCerts,
		InsecureSkipVerify: true,
	}
	var url = "mongodb://" + dbConf.User + ":" + dbConf.Pass + "@" + dbConf.Host + ":" + dbConf.Port
	fmt.Println(`mongo url:`, url)
	dialInfo, err := mgo.ParseURL(url)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println("dialInfo:", dialInfo)
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		if err != nil {
			log.Println(err)
		}
		return conn, err
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	return session
}

func (conf *DbConnConf) getConnIdent(dbconf *DbConf) string {
	return "xxx"
}

func (conf *DbConnConf) getConnConf(ctx context.Context) *DbConf {
	//TODO get session key by ctx,defalt return the first element
	return conf.Confs[0]
}
