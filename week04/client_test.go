package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

type Row struct {
	Id        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    int    `xml:"gender"`
}

var errBadOrderField = errors.New("ErrorBadOrderField")

var allowedOrderFields = map[string]struct{}{
	"":     {},
	"id":   {},
	"age":  {},
	"name": {},
}

func parseQuery(r *http.Request) (SearchRequest, error) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		return SearchRequest{}, fmt.Errorf("limit is not number: %w", err)
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		return SearchRequest{}, fmt.Errorf("offset is not number: %w", err)
	}
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")
	_, ok := allowedOrderFields[orderField]
	if !ok {
		return SearchRequest{}, errBadOrderField
	}
	orderBy, err := strconv.Atoi(r.URL.Query().Get("order_by"))
	if err != nil {
		return SearchRequest{}, fmt.Errorf("order_by is not number: %w", err)
	}
	if orderBy < -1 || orderBy > 1 {
		return SearchRequest{}, fmt.Errorf("order_by must be between -1 and 1")
	}

	return SearchRequest{
		Limit:      limit,
		Offset:     offset,
		Query:      query,
		OrderField: orderField,
		OrderBy:    orderBy,
	}, nil
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	errResp := SearchErrorResponse{}
	req, err := parseQuery(r)
	if err != nil {
		errResp.Error = err.Error()
		raw, _ := json.Marshal(errResp)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(raw)
		return
	}
	req.Limit = 10

}

func TestFindUsers(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer server.Close()

	tests := []struct {
		name       string
		client     SearchClient
		req        SearchRequest
		statusCode int
		want       *SearchResponse
		wantErr    bool
	}{
		{
			name: "access",
			client: SearchClient{
				AccessToken: "",
				URL:         server.URL,
			},
			req: SearchRequest{
				Limit: 10,
				Query: "user",
			},
			statusCode: http.StatusOK,
			want: &SearchResponse{
				nil,
				true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := tt.client
			got, err := client.FindUsers(tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got error %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
