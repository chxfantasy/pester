package pester

import (
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"
)


//ErrUnexpectedMethod occurs when an http.Client method is unable to be mapped from a calling method in the pester client
var ErrUnexpectedMethod = errors.New("unexpected client method, must be one of Do, Get, Head, Post, or PostFrom")

// ErrReadingBody happens when we cannot read the body bytes
var ErrReadingBody = errors.New("error reading body")

// ErrReadingRequestBody happens when we cannot read the request body bytes
var ErrReadingRequestBody = errors.New("error reading request body")

//no timeout param when calling pester
var ErrNoTimeOutSet = errors.New("error no timeout set for params")

//panic happened during request
var ErrPanicDuringRequest = errors.New("error panic happened during request")

//asyc回调
type RequestCallBackFun func( resp *http.Response, err error, logs string, transmissionParams ...interface{} )

// LogHook is used to log attempts as they happen. This function is never called,
// however, if KeepLog is set to true.
type LogHook func(e ErrEntry)

// BackoffStrategy is used to determine how long a retry request should wait until attempted
type BackoffStrategy func(retry int) time.Duration

// DefaultClient provides sensible defaults
var DefaultClient = &Client{Concurrency: 1, MaxRetries: 3, Backoff: DefaultBackoff, ErrLog: []ErrEntry{}}

// DefaultBackoff always returns 1 second
func DefaultBackoff(_ int) time.Duration {
	return 1 * time.Second
}

// ExponentialBackoff returns ever increasing backoffs by a power of 2
func ExponentialBackoff(i int) time.Duration {
	return time.Duration(1<<uint(i)) * time.Second
}

// ExponentialJitterBackoff returns ever increasing backoffs by a power of 2
// with +/- 0-33% to prevent sychronized reuqests.
func ExponentialJitterBackoff(i int) time.Duration {
	return jitter(int(1 << uint(i)))
}

// LinearBackoff returns increasing durations, each a second longer than the last
func LinearBackoff(i int) time.Duration {
	return time.Duration(i) * time.Second
}

// LinearJitterBackoff returns increasing durations, each a second longer than the last
// with +/- 0-33% to prevent sychronized reuqests.
func LinearJitterBackoff(i int) time.Duration {
	return jitter(i)
}

// New constructs a new DefaultClient with sensible default values
func New() *Client {
	return &Client{
		Concurrency:    DefaultClient.Concurrency,
		MaxRetries:     DefaultClient.MaxRetries,
		Backoff:        DefaultClient.Backoff,
		ErrLog:         DefaultClient.ErrLog,
		wg:             &sync.WaitGroup{},
		RetryOnHTTP429: false,
	}
}

// NewExtendedClient allows you to pass in an http.Client that is previously set up
// and extends it to have Pester's features of concurrency and retries.
func NewExtendedClient(hc *http.Client) *Client {
	c := New()
	c.hc = hc
	return c
}


////////////////////////////////////////
// Provide self-constructing variants //
////////////////////////////////////////

// Do provides the same functionality as http.Client.Do and creates its own constructor
func Do(req *http.Request, timeout time.Duration) (resp *http.Response, err error) {
	c := New()
	return c.Do(req, timeout)
}

// Get provides the same functionality as http.Client.Get and creates its own constructor
func Get(url string, timeout time.Duration) (resp *http.Response, err error) {
	c := New()
	return c.Get(url, timeout)
}

// Head provides the same functionality as http.Client.Head and creates its own constructor
func Head(url string, timeout time.Duration) (resp *http.Response, err error) {
	c := New()
	return c.Head(url, timeout)
}

// Post provides the same functionality as http.Client.Post and creates its own constructor
func Post(url string, body []byte, timeout time.Duration) (resp *http.Response, err error) {
	c := New()
	return c.Post(url, body, timeout)
}

// PostForm provides the same functionality as http.Client.PostForm and creates its own constructor
func PostForm(url string, data url.Values, timeout time.Duration) (resp *http.Response, err error) {
	c := New()
	return c.PostForm(url, data, timeout)
}


