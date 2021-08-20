package link

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/schema"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/require"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/model"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/mock"
)

func TestList(t *testing.T) {
	numberOfLinks := 5
	links := randomLinks(numberOfLinks)
	linksHashes := ""
	for _, link := range links {
		linksHashes += link.Hash
	}

	testCases := []struct {
		name            string
		queryParameters interface{}
		buildStubs      func(store *mock.MockLinkInterface)
		checkResponse   func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:            "OK",
			queryParameters: nil,
			buildStubs: func(link *mock.MockLinkInterface) {
				link.EXPECT().
					FindBy(gomock.Any()).
					Times(1).
					Return(links, nil)
				link.EXPECT().CountByQuery(gomock.Any()).
					Times(1).
					Return(pointy.Uint64(uint64(numberOfLinks)), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var linkList linksList
				err = json.Unmarshal(data, &linkList)
				require.NoError(t, err)

				actual, err := json.Marshal(links)
				require.NoError(t, err)
				expected, err := json.Marshal(linkList.Links)
				require.NoError(t, err)

				require.JSONEq(t, string(actual), string(expected))
			},
		},
		{
			name:            "InvalidQueryParams",
			queryParameters: struct{ incorrect string }{"incorrect parameter"},
			buildStubs: func(link *mock.MockLinkInterface) {
				link.EXPECT().
					FindBy(gomock.Any()).
					Times(0)
				link.EXPECT().CountByQuery(gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:            "OKGetByHashes",
			queryParameters: struct{ ids string }{linksHashes},
			buildStubs: func(link *mock.MockLinkInterface) {
				link.EXPECT().
					GetByHashes(gomock.Any()).
					Times(1).
					Return(links, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				expected, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				actual, err := json.Marshal(links)
				require.NoError(t, err)

				require.JSONEq(t, string(actual), string(expected))
			},
		},
		{
			name:            "FindByInternalServerError",
			queryParameters: nil,
			buildStubs: func(link *mock.MockLinkInterface) {
				link.EXPECT().
					FindBy(gomock.Any()).
					Times(1).
					Return(nil, errors.New("error"))
				link.EXPECT().CountByQuery(gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:            "CountByQueryInternalServerError",
			queryParameters: nil,
			buildStubs: func(link *mock.MockLinkInterface) {
				link.EXPECT().
					FindBy(gomock.Any()).
					Times(1).
					Return(links, nil)
				link.EXPECT().CountByQuery(gomock.Any()).
					Times(1).
					Return(nil, errors.New("error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:            "GetByHashesInternalServerError",
			queryParameters: struct{ ids string }{linksHashes},
			buildStubs: func(link *mock.MockLinkInterface) {
				link.EXPECT().
					GetByHashes(gomock.Any()).
					Times(1).
					Return(nil, errors.New("error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			link := mock.NewMockLinkInterface(ctrl)

			tc.buildStubs(link)

			var getLinksURL = "/api/v1/links"

			if tc.queryParameters != nil {
				var encoder = schema.NewEncoder()
				form := url.Values{}

				err := encoder.Encode(tc.queryParameters, form)
				require.NoError(t, err)

				getLinksURL += "?" + form.Encode()
			}

			request := httptest.NewRequest(http.MethodGet, getLinksURL, nil)

			recorder := httptest.NewRecorder()
			lc := newLinkController(link, nil)

			lc.List(recorder, request)

			tc.checkResponse(recorder)
		})
	}
}

func randomLinks(n int) []*model.Link {
	faker := gofakeit.NewCrypto()

	links := make([]*model.Link, n)

	for i := 0; i < n; i++ {
		links[i] = &model.Link{
			ID:        faker.Uint64(),
			URL:       faker.URL(),
			Hash:      faker.LetterN(6),
			CreatedAt: faker.Date(),
		}
	}

	return links
}

func TestValidateQueryParams(t *testing.T) {
	testCases := []struct {
		name            string
		queryParameters *listParameters
		checkResponse   func(err error)
	}{
		{
			name: "OK",
			queryParameters: &listParameters{
				Page:  pointy.Uint64(1),
				Limit: pointy.Uint64(1),
				Sort:  pointy.String("url"),
				Order: pointy.String("asc"),
			},
			checkResponse: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "LimitToLow",
			queryParameters: &listParameters{
				Page:  pointy.Uint64(1),
				Limit: pointy.Uint64(0),
			},
			checkResponse: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "limit can not be less or equal 0")
			},
		},
		{
			name: "LimitToHigh",
			queryParameters: &listParameters{
				Page:  pointy.Uint64(1),
				Limit: pointy.Uint64(1000),
			},
			checkResponse: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "maximum limit is 100")
			},
		},
		{
			name: "PageToLow",
			queryParameters: &listParameters{
				Page:  pointy.Uint64(0),
				Limit: pointy.Uint64(1),
			},
			checkResponse: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "page can not be less or equal 0")
			},
		},
		{
			name: "SortIncorrect",
			queryParameters: &listParameters{
				Page:  pointy.Uint64(1),
				Limit: pointy.Uint64(1),
				Sort:  pointy.String("error"),
			},
			checkResponse: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "invalid sort value, available values: url, hash, created_at")
			},
		},
		{
			name: "OrderIncorrect",
			queryParameters: &listParameters{
				Page:  pointy.Uint64(1),
				Limit: pointy.Uint64(1),
				Order: pointy.String("error"),
			},
			checkResponse: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "invalid order value, available values: asc, desc")
			},
		},
		{
			name: "QueryToLong",
			queryParameters: &listParameters{
				Page:  pointy.Uint64(1),
				Limit: pointy.Uint64(1),
				Query: pointy.String(gofakeit.LetterN(1001)),
			},
			checkResponse: func(err error) {
				require.Error(t, err)
				require.EqualError(t, err, "query is to long")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err := validateQueryParams(tc.queryParameters)
			tc.checkResponse(err)
		})
	}
}
