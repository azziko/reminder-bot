package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"remindbot/lib/e"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Update(offset int, limit int) ([]Update, error) {
	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, query)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err != json.Unmarshal(data, &res) {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string, preview string) error {
	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chatID))
	query.Add("text", text)
	query.Add("disable_web_page_preview", preview)

	_, err := c.doRequest(sendMessageMethod, query)
	if err != nil {
		e.Wrap("failed to send message", err)
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	const errMsg = "failed to do request"

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	return body, nil
}
