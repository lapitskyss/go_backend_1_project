package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/services/linksrv"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/store/mock"
)

var faker = gofakeit.NewCrypto()

func TestCreate(t *testing.T) {
	randomLink := getRandomLink()

	testCases := []struct {
		name          string
		body          *createLinkRequest
		buildStubs    func(store *mock.MockLinkStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: &createLinkRequest{
				URL: randomLink.URL,
			},
			buildStubs: func(link *mock.MockLinkStore) {
				link.EXPECT().
					GetByURL(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil, nil)
				link.EXPECT().
					GetNextId(gomock.Any()).
					Times(1).
					Return(randomLink.ID, nil)
				link.EXPECT().
					Add(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				var response = &singleLinkResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), response)
				require.NoError(t, err)
				require.Equal(t, response.URL, randomLink.URL)
			},
		},
		{
			name: "IncorrectURL",
			body: &createLinkRequest{
				URL: "incorrect url",
			},
			buildStubs: func(link *mock.MockLinkStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "LinkAlreadyExist",
			body: &createLinkRequest{
				URL: randomLink.URL,
			},
			buildStubs: func(link *mock.MockLinkStore) {
				link.EXPECT().
					GetByURL(gomock.Any(), gomock.Any()).
					Times(1).
					Return(randomLink, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var response = &singleLinkResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), response)
				require.NoError(t, err)
				require.Equal(t, response.URL, randomLink.URL)
				require.Equal(t, response.Hash, randomLink.Hash)
				require.Equal(t, response.CreatedAt, randomLink.CreatedAt)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			var err error

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLinkStore := mock.NewMockLinkStore(ctrl)
			linkHandler := linkHandlerForTest(mockLinkStore)
			tc.buildStubs(mockLinkStore)

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewReader(data))
			recorder := httptest.NewRecorder()

			linkHandler.Create(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestList(t *testing.T) {
	testCases := []struct {
		name          string
		query         string
		buildStubs    func(store *mock.MockLinkStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			query: "",
			buildStubs: func(link *mock.MockLinkStore) {
				linkChannel := func(_ interface{}, _ interface{}) (<-chan linksrv.Link, <-chan error) {
					linkChannel := make(chan linksrv.Link, 1)
					errorChannel := make(chan error, 1)
					defer close(linkChannel)
					defer close(errorChannel)

					return linkChannel, errorChannel
				}

				link.EXPECT().
					GetByHashes(gomock.Any(), gomock.Any()).
					Times(1).
					DoAndReturn(linkChannel)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "IncorrectQueryParams",
			query:      "?test=example",
			buildStubs: func(link *mock.MockLinkStore) {},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockLinkStore := mock.NewMockLinkStore(ctrl)
			linkHandler := linkHandlerForTest(mockLinkStore)
			tc.buildStubs(mockLinkStore)

			request := httptest.NewRequest(http.MethodPost, "/api/v1/links"+tc.query, nil)
			recorder := httptest.NewRecorder()

			linkHandler.List(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func getRandomLink() *linksrv.Link {
	return &linksrv.Link{
		ID:        faker.Uint64(),
		URL:       faker.URL(),
		Hash:      faker.LetterN(6),
		CreatedAt: faker.Date(),
	}
}

func linkHandlerForTest(mockLinkStore *mock.MockLinkStore) *LinkHandler {
	logger := zap.NewNop()
	linkService := linksrv.NewLinkService(logger, mockLinkStore)

	return NewLinkHandler(logger, linkService)
}
