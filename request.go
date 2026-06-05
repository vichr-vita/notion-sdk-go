package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func (c *Client) newRequest(ctx context.Context, method, path string, query any, body any) (*http.Request, error) {
	u, err := c.endpoint(path)
	if err != nil {
		return nil, err
	}

	values, err := encodeQuery(query)
	if err != nil {
		return nil, err
	}
	u.RawQuery = values.Encode()

	var reader io.Reader
	if body != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
		reader = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.config.token)
	req.Header.Set("Notion-Version", c.config.version)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) do(req *http.Request, out any) error {
	if c.config.requestHook != nil {
		if err := c.config.requestHook(req); err != nil {
			return err
		}
	}

	resp, err := c.config.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.config.responseHook != nil {
		if err := c.config.responseHook(resp); err != nil {
			return err
		}
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return newAPIError(resp, data)
	}

	if len(bytes.TrimSpace(data)) == 0 || out == nil {
		return nil
	}

	if raw, ok := out.(*json.RawMessage); ok {
		*raw = append((*raw)[:0], data...)
		return nil
	}

	return json.Unmarshal(data, out)
}

func (c *Client) endpoint(path string) (*url.URL, error) {
	base, err := url.Parse(c.config.baseURL)
	if err != nil {
		return nil, err
	}

	base.Path = strings.TrimRight(base.Path, "/") + "/" + strings.TrimLeft(path, "/")
	return base, nil
}

func encodeQuery(query any) (url.Values, error) {
	values := url.Values{}
	if query == nil {
		return values, nil
	}

	switch q := query.(type) {
	case url.Values:
		return q, nil
	case map[string]string:
		for key, value := range q {
			if value != "" {
				values.Set(key, value)
			}
		}
		return values, nil
	case map[string]any:
		for key, value := range q {
			if err := addQueryValue(values, key, reflect.ValueOf(value), false); err != nil {
				return nil, err
			}
		}
		return values, nil
	}

	rv := reflect.ValueOf(query)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return values, nil
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("notion: unsupported query type %T", query)
	}

	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue
		}

		name, omitempty := queryFieldName(field)
		if name == "" || name == "-" {
			continue
		}

		if err := addQueryValue(values, name, rv.Field(i), omitempty); err != nil {
			return nil, err
		}
	}

	return values, nil
}

func queryFieldName(field reflect.StructField) (string, bool) {
	for _, tagName := range []string{"url", "json"} {
		tag := field.Tag.Get(tagName)
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		return parts[0], contains(parts[1:], "omitempty")
	}

	return field.Name, false
}

func addQueryValue(values url.Values, name string, value reflect.Value, omitempty bool) error {
	if !value.IsValid() {
		return nil
	}

	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}

	if omitempty && value.IsZero() {
		return nil
	}

	switch value.Kind() {
	case reflect.String:
		values.Set(name, value.String())
	case reflect.Bool:
		values.Set(name, strconv.FormatBool(value.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		values.Set(name, strconv.FormatInt(value.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		values.Set(name, strconv.FormatUint(value.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		values.Set(name, strconv.FormatFloat(value.Float(), 'f', -1, value.Type().Bits()))
	default:
		return fmt.Errorf("notion: unsupported query field %q type %s", name, value.Type())
	}

	return nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}
