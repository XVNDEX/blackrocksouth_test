package service

import (
	"context"
	"github.com/XVNDEX/blackrocksouth_test/internal/entity"
	"github.com/XVNDEX/blackrocksouth_test/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	BTCToUSD float64 = 16817.67
	ETHToUSD float64 = 1218.58
)

// CatalogServer is the server API for Catalog service.
type CatalogServer struct {
	pb.UnimplementedCatalogServer
	Catalog *entity.Catalog
}

// NewCatalogServer returns new CatalogServer
func NewCatalogServer(catalog *entity.Catalog) *CatalogServer {
	return &CatalogServer{
		Catalog: catalog,
	}
}

func (c *CatalogServer) GetProducts(ctx context.Context, in *pb.GetCatalogRequest) (*pb.GetCatalogResponse, error) {
	resp := new(pb.GetCatalogResponse)
	var products []*pb.Product
	for _, product := range c.Catalog.Products {
		products = append(products, product.MapProductEntityToProto())
	}
	resp.Products = products
	return resp, nil
}

func (c *CatalogServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	for _, product := range c.Catalog.Products {
		if uint32(product.Id) == req.GetId() {
			return &pb.GetProductResponse{
				Product:  product.MapProductEntityToProto(),
				BtcPrice: float64(product.Price) / BTCToUSD,
				EthPrice: float64(product.Price) / ETHToUSD,
			}, nil
		}
	}

	return nil, status.Errorf(codes.InvalidArgument, "incorrect id of product")
}

func (c *CatalogServer) ConvertCurrency(ctx context.Context, req *pb.ConvertCurrencyRequest) (*pb.ConvertCurrencyResponse, error) {
	switch curr := req.Currency; curr {
	case "BTC":
		return &pb.ConvertCurrencyResponse{
			Price: float64(req.Price) * BTCToUSD,
		}, nil
	case "ETH":
		return &pb.ConvertCurrencyResponse{
			Price: float64(req.Price) * ETHToUSD,
		}, nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unacceptable currency")
	}
}
