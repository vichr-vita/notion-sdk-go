package notion

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
)

// FileUploadsService handles file upload endpoints.
type FileUploadsService struct {
	client *Client
}

// FileUploadStatus filters file upload list results.
type FileUploadStatus string

const (
	FileUploadStatusPending  FileUploadStatus = "pending"
	FileUploadStatusUploaded FileUploadStatus = "uploaded"
	FileUploadStatusExpired  FileUploadStatus = "expired"
	FileUploadStatusFailed   FileUploadStatus = "failed"
)

// ListFileUploadsRequest is the query for listing file uploads.
type ListFileUploadsRequest struct {
	Status      FileUploadStatus `json:"status,omitempty" url:"status,omitempty"`
	StartCursor string           `json:"start_cursor,omitempty" url:"start_cursor,omitempty"`
	PageSize    int              `json:"page_size,omitempty" url:"page_size,omitempty"`
}

// CreateFileUploadRequest is the JSON body for creating a file upload.
type CreateFileUploadRequest struct {
	Mode          string `json:"mode,omitempty"`
	Filename      string `json:"filename,omitempty"`
	ContentType   string `json:"content_type,omitempty"`
	NumberOfParts int    `json:"number_of_parts,omitempty"`
	ExternalURL   string `json:"external_url,omitempty"`
}

// SendFileUploadRequest is the multipart body for sending upload bytes.
type SendFileUploadRequest struct {
	Filename    string
	ContentType string
	Reader      io.Reader
	PartNumber  int
}

// List retrieves one page of file uploads for the current bot connection.
func (s *FileUploadsService) List(ctx context.Context, pagination *Pagination) (*PaginatedResponse[FileUpload], error) {
	return s.ListWithStatus(ctx, ListFileUploadsRequest{
		StartCursor: paginationStartCursor(pagination),
		PageSize:    paginationPageSize(pagination),
	})
}

// ListWithStatus retrieves one page of file uploads, optionally filtered by status.
func (s *FileUploadsService) ListWithStatus(ctx context.Context, request ListFileUploadsRequest) (*PaginatedResponse[FileUpload], error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "file_uploads", request, nil)
	if err != nil {
		return nil, err
	}

	var resp PaginatedResponse[FileUpload]
	if err := s.client.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Create initiates a file upload.
func (s *FileUploadsService) Create(ctx context.Context, request CreateFileUploadRequest) (*FileUpload, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "file_uploads", nil, request)
	if err != nil {
		return nil, err
	}

	var upload FileUpload
	if err := s.client.do(req, &upload); err != nil {
		return nil, err
	}
	return &upload, nil
}

// Retrieve gets a file upload by ID.
func (s *FileUploadsService) Retrieve(ctx context.Context, fileUploadID string) (*FileUpload, error) {
	req, err := s.client.newRequest(ctx, http.MethodGet, "file_uploads/"+fileUploadID, nil, nil)
	if err != nil {
		return nil, err
	}

	var upload FileUpload
	if err := s.client.do(req, &upload); err != nil {
		return nil, err
	}
	return &upload, nil
}

// Send transmits file bytes for a file upload.
func (s *FileUploadsService) Send(ctx context.Context, fileUploadID string, body SendFileUploadRequest) (*FileUpload, error) {
	if body.Reader == nil {
		return nil, fmt.Errorf("notion: file upload send request has nil reader")
	}

	req, err := s.client.newMultipartRequest(ctx, http.MethodPost, "file_uploads/"+fileUploadID+"/send", body)
	if err != nil {
		return nil, err
	}

	var upload FileUpload
	if err := s.client.do(req, &upload); err != nil {
		return nil, err
	}
	return &upload, nil
}

// Complete finalizes a multi-part file upload.
func (s *FileUploadsService) Complete(ctx context.Context, fileUploadID string) (*FileUpload, error) {
	req, err := s.client.newRequest(ctx, http.MethodPost, "file_uploads/"+fileUploadID+"/complete", nil, nil)
	if err != nil {
		return nil, err
	}

	var upload FileUpload
	if err := s.client.do(req, &upload); err != nil {
		return nil, err
	}
	return &upload, nil
}

func (c *Client) newMultipartRequest(ctx context.Context, method, path string, body SendFileUploadRequest) (*http.Request, error) {
	u, err := c.endpoint(path)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	if body.PartNumber > 0 {
		if err := writer.WriteField("part_number", strconv.Itoa(body.PartNumber)); err != nil {
			return nil, err
		}
	}

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeQuotes(body.Filename)))
	if body.ContentType != "" {
		header.Set("Content-Type", body.ContentType)
	}
	part, err := writer.CreatePart(header)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, body.Reader); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, err
	}
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(buf.Bytes())), nil
	}
	req.ContentLength = int64(buf.Len())
	req.Header.Set("Authorization", "Bearer "+c.config.token)
	req.Header.Set("Notion-Version", c.config.version)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

func paginationStartCursor(pagination *Pagination) string {
	if pagination == nil {
		return ""
	}
	return pagination.StartCursor
}

func paginationPageSize(pagination *Pagination) int {
	if pagination == nil {
		return 0
	}
	return pagination.PageSize
}

func escapeQuotes(value string) string {
	replacer := strings.NewReplacer("\\", "\\\\", `"`, `\"`)
	return replacer.Replace(value)
}
