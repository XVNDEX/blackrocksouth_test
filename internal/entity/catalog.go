package entity

import "github.com/XVNDEX/blackrocksouth_test/pb"

type Catalog struct {
	Products []*Product `json:"products"`
}

type Product struct {
	Id    int    `json:"id"`
	Price int    `json:"price"`
	Qty   int    `json:"qty"`
	Name  string `json:"name"`
}

func (p *Product) MapProductEntityToProto() *pb.Product {
	return &pb.Product{
		Id:    uint32(p.Id),
		Price: uint32(p.Price),
		Qty:   uint32(p.Qty),
		Name:  p.Name,
	}
}
