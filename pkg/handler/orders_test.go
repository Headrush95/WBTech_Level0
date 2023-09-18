package handler

import (
	"WBTech_Level0/models"
	"WBTech_Level0/pkg/repository"
	"WBTech_Level0/pkg/service"
	mock_service "WBTech_Level0/pkg/service/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"testing"
)

var (
	testOrderJSON = `{
	  "order_uid": "b563feb7b2b84b6test",
	  "track_number": "WBILMTESTTRACK",
	  "entry": "WBIL",
	  "delivery": {
		"name": "Test Testov",
		"phone": "+9720000000",
		"zip": "2639809",
		"city": "Kiryat Mozkin",
		"address": "Ploshad Mira 15",
		"region": "Kraiot",
		"email": "test@gmail.com"
	  },
	  "payment": {
		"transaction": "b563feb7b2b84b6test",
		"request_id": "",
		"currency": "USD",
		"provider": "wbpay",
		"amount": 1817,
		"payment_dt": 1637907727,
		"bank": "alpha",
		"delivery_cost": 1500,
		"goods_total": 317,
		"custom_fee": 0
	  },
	  "items": [
		{
		  "chrt_id": 9934930,
		  "track_number": "WBILMTESTTRACK",
		  "price": 453,
		  "rid": "ab4219087a764ae0btest",
		  "name": "Mascaras",
		  "sale": 30,
		  "size": "0",
		  "total_price": 317,
		  "nm_id": 2389212,
		  "brand": "Vivienne Sabo",
		  "status": 202
		}
	  ],
	  "locale": "en",
	  "internal_signature": "",
	  "customer_id": "test",
	  "delivery_service": "meest",
	  "shardkey": "9",
	  "sm_id": 99,
	  "date_created": "2021-11-26T06:22:19Z",
	  "oof_shard": "1"
	}`
	testInvalidDataOrderJSON = `{
	  "order_uid": "b563feb7b2b84b6tost12",
	  "track_number": "WBILMTESTTRACK",
	  "entry": "WBIL",
	  "delivery": {
		"name": "Test Testov",
		"phone": "+9720000000",
		"zip": "2639809",
		"city": "Kiryat Mozkin",
		"address": "Ploshad Mira 15",
		"region": "Kraiot",
		"email": "test@gmail.com"
	  },
	  "payment": {
		"transaction": "b563feb7b2b84b6test",
		"request_id": "",
		"currency": "USD",
		"provider": "wbpay",
		"amount": 1817,
		"payment_dt": 1637907727,
		"bank": "alpha",
		"delivery_cost": 1500,
		"goods_total": 317,
		"custom_fee": 0
	  },
	  "items": [
		{
		  "chrt_id": 9934930,
		  "track_number": "WBILMTESTTRACK",
		  "price": 453,
		  "rid": "ab4219087a764ae0btest",
		  "name": "Mascaras",
		  "sale": 30,
		  "size": "0",
		  "total_price": 317,
		  "nm_id": 2389212,
		  "brand": "Vivienne Sabo",
		  "status": 202
		}
	  ],
	  "locale": "en",
	  "internal_signature": "",
	  "customer_id": "test",
	  "delivery_service": "meest",
	  "shardkey": "9",
	  "sm_id": 99,
	  "date_created": "2021-11-26T06:22:19Z",
	  "oof_shard": "1"
	}`
)

func TestHandler_GetOrderCache(t *testing.T) {
	type mockBehavior func(r *mock_service.MockOrderCache, uid string)
	var order models.Order
	err := json.Unmarshal([]byte(testOrderJSON), &order)
	if err != nil {
		t.Fatalf("error occurred while unmarshalling test order: %s", err.Error())
	}
	// из-за кастомного типа даты paymentDate приходится заново маршалить в JSON
	orderJSON, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("error occurred while marshalling test order: %s", err.Error())
	}

	tests := []struct {
		name                 string
		inputUid             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "OK",
			inputUid: "b563feb7b2b84b6test",
			mockBehavior: func(r *mock_service.MockOrderCache, uid string) {
				r.EXPECT().GetOrder(uid).Return(order, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: string(orderJSON),
		},
		{
			name:     "incorrect order uid",
			inputUid: "b563feb7b2b84b6",
			mockBehavior: func(r *mock_service.MockOrderCache, uid string) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"incorrect order uid"}`,
		},
		{
			name:     "no value in cache",
			inputUid: "b563feb7b2b84b6test",
			mockBehavior: func(r *mock_service.MockOrderCache, uid string) {
				r.EXPECT().GetOrder(uid).Return(models.Order{}, repository.CacheHasNoValue)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"error occurred while trying get order with uid b563feb7b2b84b6test: there is no such order uid in the cache"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockOrderCache(c)
			test.mockBehavior(repo, test.inputUid)

			services := &service.Service{OrderCache: repo}
			handler := Handler{services: services}

			r := gin.New()
			r.GET("/orders/:id", handler.GetOrderById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/orders/"+test.inputUid, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_CreateOrder(t *testing.T) {
	type mockBehavior func(r *mock_service.MockOrdersDB, c *mock_service.MockOrderCache, order models.Order)
	//var order models.Order
	//err := json.Unmarshal([]byte(testOrderJSON), &order)
	//if err != nil {
	//	t.Fatalf("error occurred while unmarshalling test order: %s", err.Error())
	//}

	tests := []struct {
		name                 string
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: testOrderJSON,
			mockBehavior: func(r *mock_service.MockOrdersDB, c *mock_service.MockOrderCache, order models.Order) {
				r.EXPECT().CreateOrder(order).Return(nil)
				c.EXPECT().PutOrder(order).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"order_uid":"b563feb7b2b84b6test"}`,
		},
		{
			name:      "Order already exists",
			inputBody: testOrderJSON,
			mockBehavior: func(r *mock_service.MockOrdersDB, c *mock_service.MockOrderCache, order models.Order) {
				r.EXPECT().CreateOrder(order).Return(errors.New("unable to inserting data into \"orders\" table: pq: duplicate key value violates unique constraint \"orders_pkey\""))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"error occurred while creating order DB entry: unable to inserting data into \"orders\" table: pq: duplicate key value violates unique constraint \"orders_pkey\""}`,
		},
		{
			name:      "Invalid JSON data",
			inputBody: testInvalidDataOrderJSON,
			mockBehavior: func(r *mock_service.MockOrdersDB, c *mock_service.MockOrderCache, order models.Order) {
				r.EXPECT().CreateOrder(order).Return(validator.ValidationErrors{})
			},
			expectedStatusCode: 400,
			/*
				отсутсвует конкретика по ошибке в виду необходиомсти проверки результата, возвращаемого CreateOrder,
				на содержание ошибок валидатора и возвращение соответствующего статус кода
			*/
			expectedResponseBody: `{"error":"error occurred while creating order DB entry: "}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockOrdersDB(c)
			cache := mock_service.NewMockOrderCache(c)

			var order models.Order
			err := json.Unmarshal([]byte(test.inputBody), &order)
			if err != nil {
				t.Fatalf("error occurred while unmarshalling test order: %s", err.Error())
			}

			test.mockBehavior(repo, cache, order)

			services := &service.Service{OrdersDB: repo, OrderCache: cache}
			handler := Handler{services: services}

			r := gin.New()
			r.POST("/orders/create", handler.CreateOrder)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/orders/create", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
