package link

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/model"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/mock"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/repository"
)

func TestCreate(t *testing.T) {
	randomLink := randomLink()

	testCases := []struct {
		name          string
		body          *createLinkParams
		buildStubs    func(store *mock.MockLinkInterface)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: &createLinkParams{
				URL: randomLink.URL,
			},
			buildStubs: func(link *mock.MockLinkInterface) {
				link.EXPECT().
					GetByURL(gomock.Any()).
					Times(1).
					Return(nil, nil)
				link.EXPECT().
					GetNextId().
					Times(1).
					Return(pointy.Uint64(randomLink.ID), nil)
				link.EXPECT().Add(gomock.Any()).
					Times(1).
					Return(randomLink, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				expectedResult, err := json.Marshal(randomLink)
				require.NoError(t, err)
				require.JSONEq(t, string(expectedResult), recorder.Body.String())
			},
		},
		{
			name: "EmptyBody",
			body: nil,
			buildStubs: func(link *mock.MockLinkInterface) {

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.JSONEq(t, "{\"status\":false, \"error\":\"body is empty\"}", recorder.Body.String())
			},
		},
		{
			name: "IncorrectURL",
			body: &createLinkParams{
				URL: "",
			},
			buildStubs: func(link *mock.MockLinkInterface) {

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.JSONEq(t, "{\"status\":false, \"error\":\"incorrect URL\"}", recorder.Body.String())
			},
		},
		{
			name: "ToLongURL",
			body: &createLinkParams{
				URL: "https://" + gofakeit.DigitN(10001) + ".com",
			},
			buildStubs: func(link *mock.MockLinkInterface) {

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				require.JSONEq(t, "{\"status\":false, \"error\":\"URL is to long\"}", recorder.Body.String())
			},
		},
		{
			name: "CannotGetLinkByUrl",
			body: &createLinkParams{
				URL: randomLink.URL,
			},
			buildStubs: func(link *mock.MockLinkInterface) {
				link.EXPECT().
					GetByURL(gomock.Any()).
					Times(1).
					Return(nil, errors.New("error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "LinkAlreadyExist",
			body: &createLinkParams{
				URL: randomLink.URL,
			},
			buildStubs: func(link *mock.MockLinkInterface) {
				link.EXPECT().
					GetByURL(gomock.Any()).
					Times(1).
					Return(randomLink, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "FailGetNextId",
			body: &createLinkParams{
				URL: randomLink.URL,
			},
			buildStubs: func(link *mock.MockLinkInterface) {
				link.EXPECT().
					GetByURL(gomock.Any()).
					Times(1).
					Return(nil, nil)
				link.EXPECT().
					GetNextId().
					Times(1).
					Return(nil, errors.New("error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "FailSaveUrl",
			body: &createLinkParams{
				URL: randomLink.URL,
			},
			buildStubs: func(link *mock.MockLinkInterface) {
				link.EXPECT().
					GetByURL(gomock.Any()).
					Times(1).
					Return(nil, nil)
				link.EXPECT().
					GetNextId().
					Times(1).
					Return(pointy.Uint64(randomLink.ID), nil)
				link.EXPECT().Add(gomock.Any()).
					Times(1).
					Return(randomLink, errors.New("error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			var err error

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			link := mock.NewMockLinkInterface(ctrl)

			tc.buildStubs(link)

			data := []byte("")
			if tc.body != nil {
				data, err = json.Marshal(tc.body)
				require.NoError(t, err)
			}

			request := httptest.NewRequest(http.MethodPost, "/api/v1/links", bytes.NewReader(data))

			recorder := httptest.NewRecorder()
			lc := newLinkController(link, nil)

			lc.Create(recorder, request)

			tc.checkResponse(recorder)
		})
	}
}

func randomLink() *model.Link {
	faker := gofakeit.NewCrypto()

	return &model.Link{
		ID:        faker.Uint64(),
		URL:       faker.URL(),
		Hash:      faker.LetterN(6),
		CreatedAt: faker.Date(),
	}
}

func newLinkController(linkInterface repository.LinkInterface, redirectLogInterface repository.RedirectLogInterface) *linkController {
	ctx := context.Background()
	logger := zap.NewNop()
	sugar := logger.Sugar()

	var store mock.Store
	store.SetLink(linkInterface)
	store.SetRedirectLog(redirectLogInterface)

	return New(ctx, sugar, &store)
}
