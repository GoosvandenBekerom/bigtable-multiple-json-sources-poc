package endpoints

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/GoosvandenBekerom/intaker-bigtable-poc/data"
	"github.com/google/uuid"
)

type ProductAPI struct {
	store data.ProductStore
}

func NewProductAPI(store data.ProductStore) ProductAPI {
	return ProductAPI{store: store}
}

func (api ProductAPI) GenerateProducts(response http.ResponseWriter, request *http.Request) {
	amount := 10
	var err error
	if q := request.URL.Query().Get("amount"); q != "" {
		amount, err = strconv.Atoi(q)
		if err != nil {
			fail(response, err)
			return
		}
	}

	var ids []string

	for i := 0; i < amount; i++ {
		id := uuid.New().String()

		log.Printf("generating product with id %s\n", id)

		err = api.store.SaveProduct(request.Context(), data.Product{
			ID:          id,
			Title:       "title of " + id,
			Description: "description of " + id,
		})
		if err != nil {
			fail(response, err)
			return
		}

		for offerCounter := 0; offerCounter < rand.Intn(5); offerCounter++ {
			err = api.store.SaveOffer(request.Context(), data.Offer{
				ProductID:    id,
				ID:           uuid.New().String(),
				PriceInCents: rand.Intn(1000_00),
			})
			if err != nil {
				fail(response, err)
				return
			}
		}

		for reviewCounter := 0; reviewCounter < rand.Intn(5); reviewCounter++ {
			err = api.store.SaveReview(request.Context(), data.Review{
				ProductID: id,
				ID:        uuid.New().String(),
				Rating:    rand.Intn(5),
				Message:   "message of review " + strconv.Itoa(reviewCounter) + " of product " + id,
			})
			if err != nil {
				fail(response, err)
				return
			}
		}
	}

	response.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(response)
	err = encoder.Encode(ids)
	if err != nil {
		fail(response, err)
		return
	}
}

func (api ProductAPI) GetProducts(response http.ResponseWriter, request *http.Request) {
	var limit int
	var err error
	if q := request.URL.Query().Get("limit"); q != "" {
		limit, err = strconv.Atoi(q)
	}

	products, err := api.store.GetAllProducts(request.Context(), int64(limit))
	if err != nil {
		fail(response, err)
	}

	response.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(response)
	err = encoder.Encode(products)
	if err != nil {
		fail(response, err)
	}
}

func fail(response http.ResponseWriter, err error) {
	response.WriteHeader(http.StatusInternalServerError)
	response.Write([]byte(err.Error()))
}
