package sendcloud

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Event struct {
	name   string
	time   time.Time
	rcpt   string
	msgid  string
	reason string
}

func (e *Event) Name() string    { return e.name }
func (e *Event) Time() time.Time { return e.time }
func (e *Event) Rcpt() string    { return e.rcpt }
func (e *Event) MsgId() string   { return e.msgid }
func (e *Event) Reason() string  { return e.reason }

var (
	ErrMethodNotAllowed = fmt.Errorf("method not allowed")
	ErrBadSignature     = fmt.Errorf("bad signature")
	ErrInvalidTimestamp = fmt.Errorf("invalid timestamp")
	ErrInvalidForm      = fmt.Errorf("invalid form data")
)

type Webhook struct {
	key string
}

func NewWebhook(key string) *Webhook {
	return &Webhook{key}
}

func (wh *Webhook) Handle(w http.ResponseWriter, req *http.Request) (evt *Event, err error) {
	if req.Method != "POST" {
		err = ErrMethodNotAllowed
		w.Header().Set("Allow", "POST")
		http.Error(w, "only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err = req.ParseForm(); err != nil {
		err = ErrInvalidForm
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	ts := req.Form.Get("timestamp")
	token := req.Form.Get("token")
	signature := req.Form.Get("signature")
	calcSig := wh.Signature(ts, token)
	if calcSig != signature {
		err = ErrBadSignature
		log.Printf("ERROR signature mismatch ts:%s, token:%s, sign:%s, calcSig:%s",
			ts, token, signature, calcSig)
		http.Error(w, "bad signature", http.StatusForbidden)
		return
	}

	unix, err := strconv.ParseInt(ts, 10, 64) // millisecond since Unix epoch
	if err != nil {
		err = ErrInvalidTimestamp
		http.Error(w, "invalid timestamp", http.StatusBadRequest)
		return
	}
	evt = &Event{
		time:   time.Unix(0, unix*1e6), // 1 ms = 1e6 ns
		name:   req.Form.Get("event"),
		rcpt:   req.Form.Get("recipient"),
		msgid:  req.Form.Get("emailId"),
		reason: req.Form.Get("message") + ": " + req.Form.Get("reason"),
	}
	return
}

func (wh *Webhook) Signature(timestamp, token string) (calcSig string) {
	h := hmac.New(sha256.New, []byte(wh.key))
	io.WriteString(h, timestamp)
	io.WriteString(h, token)
	calcSig = hex.EncodeToString(h.Sum(nil))
	return
}
