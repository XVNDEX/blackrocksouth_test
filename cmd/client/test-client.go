package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/XVNDEX/blackrocksouth_test/internal/data"
	"github.com/XVNDEX/blackrocksouth_test/internal/service"
	"github.com/XVNDEX/blackrocksouth_test/pb"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
)

func main() {
	// 1 request: Creating token for given user
	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	authClient := pb.NewAuthServiceClient(conn)

	login := "user1"
	pwd := "qwerty1"
	req := &pb.SignInUserInput{
		Login: login,
		Pwd:   pwd,
	}

	log.Printf("1 request: SignInUser")
	log.Printf("Trying to get token for %v\n", login)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	responseSignIn, err := authClient.SignInUser(ctx, req)
	if err != nil {
		grpclog.Fatalf("authentication failed: %v", err)
		return
	}
	log.Println("Authentication succeeded, token created:", responseSignIn.GetToken())

	// Requests with credentials
	// Set up the credentials for the catalog connection
	creds, err := credentials.NewClientTLSFromFile(data.Path("ca_cert.pem"), "x.test.example.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(NewOauthAccess(service.FetchToken(responseSignIn.GetToken()))),
		grpc.WithTransportCredentials(creds),
	}

	connWithAuth, err := grpc.Dial("0.0.0.0:9090", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer connWithAuth.Close()

	catalogClient := pb.NewCatalogClient(connWithAuth)

	// 2 request: Get all products in catalog
	log.Printf("2 request: GetCatalog")
	requestGetProducts := &pb.GetCatalogRequest{}
	responseGetCatalog, err := catalogClient.GetProducts(ctx, requestGetProducts)
	if err != nil {
		grpclog.Fatalf("cannot get catalog: %v", err)
		return
	}
	log.Printf("catalog lenght: %v", len(responseGetCatalog.GetProducts()))
	log.Printf("catalog item 1: %v", responseGetCatalog.GetProducts()[0])
	log.Printf("catalog item 2: %v", responseGetCatalog.GetProducts()[1])
	log.Printf("catalog item 3: %v", responseGetCatalog.GetProducts()[2])

	// 3 request: Get certain product by id
	log.Printf("3 request: GetProduct")
	reqGetProduct := &pb.GetProductRequest{
		Id: responseGetCatalog.GetProducts()[1].GetId(),
	}

	responseGetProduct, err := catalogClient.GetProduct(ctx, reqGetProduct)
	if err != nil {
		grpclog.Fatalf("cannot get product by id: %v", err)
		return
	}
	p := responseGetProduct.GetProduct()
	log.Printf("product id: %v, price: %v, name: %v, qty: %v", p.GetId(), p.GetPrice(), p.GetName(), p.GetQty())
	log.Printf("price in BTC: %v", responseGetProduct.GetBtcPrice())
	log.Printf("price in ETH: %v", responseGetProduct.GetEthPrice())

	// 4 request: ConvertCurrency converts given price USD to ETH or BTC
	log.Printf("4 request: ConvertCurrency")
	reqConvertCurrency := &pb.ConvertCurrencyRequest{
		Price:    p.GetPrice(),
		Currency: "ETH",
	}

	responseConvertCurrency, err := catalogClient.ConvertCurrency(ctx, reqConvertCurrency)
	if err != nil {
		grpclog.Fatalf("cannot convert currency: %v", err)
		return
	}
	log.Printf("price in USD: %v", p.GetPrice())
	log.Printf("price in ETH: %v", responseConvertCurrency.GetPrice())
}

// oauthAccess supplies PerRPCCredentials from a given token.
type oauthAccess struct {
	token oauth2.Token
}

// NewOauthAccess constructs the PerRPCCredentials using a given token.
func NewOauthAccess(token *oauth2.Token) credentials.PerRPCCredentials {
	return oauthAccess{token: *token}
}

func (oa oauthAccess) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	ri, _ := credentials.RequestInfoFromContext(ctx)
	if err := credentials.CheckSecurityLevel(ri.AuthInfo, credentials.PrivacyAndIntegrity); err != nil {
		return nil, fmt.Errorf("unable to transfer oauthAccess PerRPCCredentials: %v", err)
	}
	return map[string]string{
		"authorization": oa.token.Type() + " " + oa.token.AccessToken,
	}, nil
}

func (oa oauthAccess) RequireTransportSecurity() bool {
	return true
}
