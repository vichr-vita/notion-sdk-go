package notion

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileUploadsListBuildsPaginationQuery(t *testing.T) {
	client := newFileUploadsTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/file_uploads?page_size=50&start_cursor=cursor_a", req.URL.String())
		require.Empty(t, req.Header.Get("Content-Type"))

		return fileUploadListResponse()
	})

	resp, err := client.FileUploads.List(context.Background(), &Pagination{
		StartCursor: "cursor_a",
		PageSize:    50,
	})
	require.NoError(t, err)
	require.Equal(t, "list", resp.Object)
	require.Equal(t, "file_upload", resp.Type)
	require.Len(t, resp.Results, 1)
	require.Equal(t, "upload_a", resp.Results[0].ID)
	require.Equal(t, "pending", resp.Results[0].Status)
	require.Equal(t, "document.pdf", *resp.Results[0].Filename)
	require.Equal(t, int64(12), *resp.Results[0].ContentLength)
	require.Equal(t, 3, resp.Results[0].NumberOfParts.Total)
	require.Equal(t, 1, resp.Results[0].NumberOfParts.Sent)
	require.JSONEq(t, fileUploadResponse(), string(resp.Results[0].Raw))
}

func TestFileUploadsListWithStatusBuildsQuery(t *testing.T) {
	client := newFileUploadsTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/file_uploads?page_size=10&status=uploaded", req.URL.String())

		return fileUploadListResponse()
	})

	_, err := client.FileUploads.ListWithStatus(context.Background(), ListFileUploadsRequest{
		Status:   FileUploadStatusUploaded,
		PageSize: 10,
	})
	require.NoError(t, err)
}

func TestFileUploadsCreateBuildsRequest(t *testing.T) {
	client := newFileUploadsTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "https://notion.test/v1/file_uploads", req.URL.String())
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, "multi_part", body["mode"])
		require.Equal(t, "document.pdf", body["filename"])
		require.Equal(t, "application/pdf", body["content_type"])
		require.Equal(t, float64(3), body["number_of_parts"])

		return fileUploadResponse()
	})

	upload, err := client.FileUploads.Create(context.Background(), CreateFileUploadRequest{
		Mode:          "multi_part",
		Filename:      "document.pdf",
		ContentType:   "application/pdf",
		NumberOfParts: 3,
	})
	require.NoError(t, err)
	require.Equal(t, "upload_a", upload.ID)
	require.Equal(t, "https://notion.test/upload", upload.UploadURL)
}

func TestFileUploadsCreateSupportsExternalURL(t *testing.T) {
	client := newFileUploadsTestClient(t, func(req *http.Request) string {
		var body map[string]any
		require.NoError(t, json.NewDecoder(req.Body).Decode(&body))
		require.Equal(t, "external_url", body["mode"])
		require.Equal(t, "https://example.com/report.pdf", body["external_url"])

		return fileUploadResponse()
	})

	_, err := client.FileUploads.Create(context.Background(), CreateFileUploadRequest{
		Mode:        "external_url",
		ExternalURL: "https://example.com/report.pdf",
	})
	require.NoError(t, err)
}

func TestFileUploadsRetrieveBuildsPath(t *testing.T) {
	client := newFileUploadsTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "https://notion.test/v1/file_uploads/upload_a", req.URL.String())

		return fileUploadResponse()
	})

	upload, err := client.FileUploads.Retrieve(context.Background(), "upload_a")
	require.NoError(t, err)
	require.Equal(t, "upload_a", upload.ID)
}

func TestFileUploadsSendBuildsMultipartBody(t *testing.T) {
	client := newFileUploadsTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "https://notion.test/v1/file_uploads/upload_a/send", req.URL.String())
		require.True(t, strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/form-data; boundary="))

		reader, err := req.MultipartReader()
		require.NoError(t, err)
		fields := map[string]string{}
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
			data, err := io.ReadAll(part)
			require.NoError(t, err)
			fields[part.FormName()] = string(data)
			if part.FormName() == "file" {
				require.Equal(t, "document.pdf", part.FileName())
				require.Equal(t, "application/pdf", part.Header.Get("Content-Type"))
			}
		}
		require.Equal(t, "2", fields["part_number"])
		require.Equal(t, "hello notion", fields["file"])

		return fileUploadResponse()
	})

	upload, err := client.FileUploads.Send(context.Background(), "upload_a", SendFileUploadRequest{
		Filename:    "document.pdf",
		ContentType: "application/pdf",
		Reader:      strings.NewReader("hello notion"),
		PartNumber:  2,
	})
	require.NoError(t, err)
	require.Equal(t, "upload_a", upload.ID)
}

func TestFileUploadsSendRejectsNilReader(t *testing.T) {
	client := NewClient("secret_test", WithBaseURL("https://notion.test/v1"))

	upload, err := client.FileUploads.Send(context.Background(), "upload_a", SendFileUploadRequest{})
	require.Nil(t, upload)
	require.EqualError(t, err, "notion: file upload send request has nil reader")
}

func TestFileUploadsCompleteBuildsPath(t *testing.T) {
	client := newFileUploadsTestClient(t, func(req *http.Request) string {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "https://notion.test/v1/file_uploads/upload_a/complete", req.URL.String())
		require.Nil(t, req.Body)
		require.Empty(t, req.Header.Get("Content-Type"))

		return fileUploadResponse()
	})

	upload, err := client.FileUploads.Complete(context.Background(), "upload_a")
	require.NoError(t, err)
	require.Equal(t, "upload_a", upload.ID)
}

func TestFileUploadsReturnsAPIError(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusConflict,
			Header:     http.Header{"X-Request-Id": []string{"req_file_upload"}},
			Body:       io.NopCloser(strings.NewReader(`{"object":"error","status":409,"code":"conflict_error","message":"upload conflict"}`)),
		}, nil
	})}
	client := NewClient("secret_test", WithHTTPClient(httpClient))

	upload, err := client.FileUploads.Retrieve(context.Background(), "upload_a")
	require.Nil(t, upload)
	require.Error(t, err)

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusConflict, apiErr.StatusCode)
	require.Equal(t, "conflict_error", apiErr.Code)
	require.Equal(t, "upload conflict", apiErr.Message)
	require.Equal(t, "req_file_upload", apiErr.RequestID)
}

func newFileUploadsTestClient(t *testing.T, assert func(*http.Request) string) *Client {
	t.Helper()

	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		body := assert(req)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	})}
	return NewClient(
		"secret_test",
		WithBaseURL("https://notion.test/v1"),
		WithHTTPClient(httpClient),
	)
}

func fileUploadListResponse() string {
	return `{
		"object":"list",
		"type":"file_upload",
		"results":[` + fileUploadResponse() + `],
		"has_more":false,
		"next_cursor":null
	}`
}

func fileUploadResponse() string {
	return `{
		"object":"file_upload",
		"id":"upload_a",
		"created_time":"2026-06-06T10:00:00Z",
		"created_by":{"object":"user","id":"bot_a","name":null,"avatar_url":null},
		"last_edited_time":"2026-06-06T10:01:00Z",
		"in_trash":false,
		"expiry_time":"2026-06-06T11:00:00Z",
		"status":"pending",
		"filename":"document.pdf",
		"content_type":"application/pdf",
		"content_length":12,
		"upload_url":"https://notion.test/upload",
		"complete_url":"https://notion.test/complete",
		"number_of_parts":{"total":3,"sent":1}
	}`
}
