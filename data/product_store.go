package data

import (
	"cloud.google.com/go/bigtable"
	"context"
	"encoding/json"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

const (
	tableName = "products"

	columnFamilyProduct = "product"
	columnProduct       = "product"

	columnFamilyOffer = "offers"
	columnOfferPrefix = "offer_"

	columnFamilyReview = "reviews"
	columnReviewPrefix = "review_"
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

	err = adminClient.CreateColumnFamily(ctx, tableName, columnFamilyOffer)
	if err != nil {
		s, ok := status.FromError(err)
		if !ok || s.Code() != codes.AlreadyExists {
			return ProductStore{}, err
		}
	}

	err = adminClient.CreateColumnFamily(ctx, tableName, columnFamilyReview)
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

	return p.table.Apply(ctx, product.ID, mutation)
}

func (p ProductStore) SaveOffer(ctx context.Context, offer Offer) error {
	data, err := json.Marshal(offer)
	if err != nil {
		return err
	}

	mutation := bigtable.NewMutation()
	mutation.Set(columnFamilyOffer, columnOfferPrefix+offer.ID, bigtable.Now(), data)

	return p.table.Apply(ctx, offer.ProductID, mutation)
}

func (p ProductStore) SaveReview(ctx context.Context, review Review) error {
	data, err := json.Marshal(review)
	if err != nil {
		return err
	}

	mutation := bigtable.NewMutation()
	mutation.Set(columnFamilyReview, columnReviewPrefix+review.ID, bigtable.Now(), data)

	return p.table.Apply(ctx, review.ProductID, mutation)
}

func (p ProductStore) GetAllProducts(ctx context.Context, limit int64) ([]AggregatedProduct, error) {
	products := make([]AggregatedProduct, 0)

	err := p.table.ReadRows(ctx, bigtable.InfiniteRange(""), func(row bigtable.Row) bool {
		products = append(products, rowToAggregatedProduct(row))
		return true
	}, bigtable.RowFilter(bigtable.LatestNFilter(1)), bigtable.LimitRows(limit))
	if err != nil {
		return nil, err
	}

	return products, nil
}

// rowToAggregatedProduct maps a single bigtable row to an AggregatedProduct
// only works accurately if row is read with: bigtable.RowFilter(bigtable.LatestNFilter(1))
func rowToAggregatedProduct(row bigtable.Row) AggregatedProduct {
	var product AggregatedProduct

	// product columns
	for _, item := range row[columnFamilyProduct] {
		if item.Column == columnFamilyProduct+":"+columnProduct {
			product.Product = item.Value
		}
	}

	// offer columns
	for _, item := range row[columnFamilyOffer] {
		if strings.HasPrefix(item.Column, columnFamilyOffer+":"+columnOfferPrefix) {
			product.Offers = append(product.Offers, item.Value)
		}
	}

	// review columns
	for _, item := range row[columnFamilyReview] {
		if strings.HasPrefix(item.Column, columnFamilyReview+":"+columnReviewPrefix) {
			product.Reviews = append(product.Reviews, item.Value)
		}
	}

	return product
}
