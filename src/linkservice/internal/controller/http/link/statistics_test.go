package link

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/golang/mock/gomock"
	"github.com/lapitskyss/go_backend_1_project/src/linkservice/internal/repository/mock"
	"github.com/openlyinc/pointy"
	"github.com/stretchr/testify/require"
)

func TestStatistics(t *testing.T) {
	randomLink := randomLink()
	numberOfRedirects := gofakeit.Uint64()

	testCases := []struct {
		name          string
		hash          string
		buildStubs    func(link *mock.MockLinkInterface, redirectLog *mock.MockRedirectLogInterface)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			hash: randomLink.Hash,
			buildStubs: func(link *mock.MockLinkInterface, redirectLog *mock.MockRedirectLogInterface) {
				link.EXPECT().
					GetByHash(gomock.Any()).
					Times(1).
					Return(randomLink, nil)
				redirectLog.EXPECT().
					CountRedirects(randomLink.ID).
					Times(1).
					Return(pointy.Uint64(numberOfRedirects), nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var statisticsInfo statisticsInfo
				err = json.Unmarshal(data, &statisticsInfo)
				require.NoError(t, err)
				require.Equal(t, numberOfRedirects, statisticsInfo.Redirects)
			},
		},
		{
			name: "HashNotFound",
			hash: "",
			buildStubs: func(link *mock.MockLinkInterface, redirectLog *mock.MockRedirectLogInterface) {
				link.EXPECT().
					GetByHash(gomock.Any()).
					Times(0)
				redirectLog.EXPECT().
					CountRedirects(randomLink.ID).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "GetByHashError",
			hash: randomLink.Hash,
			buildStubs: func(link *mock.MockLinkInterface, redirectLog *mock.MockRedirectLogInterface) {
				link.EXPECT().
					GetByHash(gomock.Any()).
					Times(1).
					Return(nil, errors.New("error"))
				redirectLog.EXPECT().
					CountRedirects(randomLink.ID).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "LinkNotFound",
			hash: randomLink.Hash,
			buildStubs: func(link *mock.MockLinkInterface, redirectLog *mock.MockRedirectLogInterface) {
				link.EXPECT().
					GetByHash(gomock.Any()).
					Times(1).
					Return(nil, nil)
				redirectLog.EXPECT().
					CountRedirects(randomLink.ID).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "CountRedirectsError",
			hash: randomLink.Hash,
			buildStubs: func(link *mock.MockLinkInterface, redirectLog *mock.MockRedirectLogInterface) {
				link.EXPECT().
					GetByHash(gomock.Any()).
					Times(1).
					Return(randomLink, nil)
				redirectLog.EXPECT().
					CountRedirects(randomLink.ID).
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
			redirectLog := mock.NewMockRedirectLogInterface(ctrl)

			tc.buildStubs(link, redirectLog)

			request := httptest.NewRequest(http.MethodGet, "/links/"+tc.hash+"/statistics", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("hash", tc.hash)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			recorder := httptest.NewRecorder()
			lc := newLinkController(link, redirectLog)

			lc.Statistics(recorder, request)

			tc.checkResponse(recorder)
		})
	}
}
