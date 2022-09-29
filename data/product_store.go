package data

import (
	"cloud.google.com/go/bigtable"
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	tableName           = "products"
	columnFamilyProduct = "product"
	columnProduct       = "product"
	columnOffer         = "offer"
	columnReview        = "review"
)

type ProductStore struct {
	table *bigtable.Table
}

func NewProductStore(ctx context.Context, project string, instance string) (ProductStore, error) {
	client, err := bigtable.NewClient(ctx, project, instance)
	if err != nil {
		return ProductStore{}, err
	}

	adminClient, err := bigtable.NewAdminClient(ctx, project, instance)
	if err != nil {
		return ProductStore{}, err
	}
	_, err = adminClient.TableInfo(ctx, tableName)
	if err != nil {
		s, ok := status.FromError(err)
		if !ok && s.Code() != codes.NotFound {
			return ProductStore{}, err
		}

		err = adminClient.CreateTable(ctx, tableName)
		if err != nil {
			return ProductStore{}, err
		}
	}

	err = adminClient.CreateColumnFamily(ctx, tableName, columnFamilyProduct)
	if err != nil {
		s, ok := status.FromError(err)
		if !ok || s.Code() != codes.AlreadyExists {
			return ProductStore{}, err
		}
	}

	return ProductStore{table: client.Open(tableName)}, nil
}

func (p ProductStore) SaveProduct(ctx context.Context, product Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	mutation := bigtable.NewMutation()
	mutation.Set(columnFamilyProduct, columnProduct, bigtable.Now(), data)

	return p.table.Apply(ctx, product.Id, mutation)
}

func (p ProductStore) SaveOffer(ctx context.Context, offer Offer) error {
	data, err := json.Marshal(offer)
	if err != nil {
		return err
	}

	mutation := bigtable.NewMutation()
	mutation.Set(columnFamilyProduct, columnOffer, bigtable.Now(), data)

	return p.table.Apply(ctx, offer.ProductID, mutation)
}

func (p ProductStore) SaveReview(ctx context.Context, review Review) error {
	data, err := json.Marshal(review)
	if err != nil {
		return err
	}

	mutation := bigtable.NewMutation()
	mutation.Set(columnFamilyProduct, columnReview, bigtable.Now(), data)

	return p.table.Apply(ctx, review.ProductID, mutation)
}

func (p ProductStore) GetAllProducts(ctx context.Context, limit int64) ([]AggregatedProduct, error) {
	products := make([]AggregatedProduct, 0)

	err := p.table.ReadRows(ctx, bigtable.InfiniteRange(""), func(row bigtable.Row) bool {
		var product AggregatedProduct
		for _, item := range row[columnFamilyProduct] {
			var err error
			switch item.Column {
			case columnFamilyProduct + ":" + columnProduct:
				product.Product = item.Value
				break
			case columnFamilyProduct + ":" + columnOffer:
				product.Offer = item.Value
				break
			case columnFamilyProduct + ":" + columnReview:
				product.Review = item.Value
				break
			}
			if err != nil {
				return false
			}
		}
		products = append(products, product)
		return true
	}, bigtable.RowFilter(bigtable.LatestNFilter(1)), bigtable.LimitRows(limit))
	if err != nil {
		return nil, err
	}

	return products, nil
}
