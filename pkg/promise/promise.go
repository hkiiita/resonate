package promise

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/resonatehq/resonate/pkg/idempotency"
)

type Promise struct {
	Id                        string            `json:"id"`
	State                     State             `json:"state"`
	Param                     Value             `json:"param,omitempty"`
	Value                     Value             `json:"value,omitempty"`
	Timeout                   int64             `json:"timeout"`
	IdempotencyKeyForCreate   *idempotency.Key  `json:"idempotencyKeyForCreate,omitempty"`
	IdempotencyKeyForComplete *idempotency.Key  `json:"idempotencyKeyForComplete,omitempty"`
	CreatedOn                 *int64            `json:"createdOn,omitempty"`
	CompletedOn               *int64            `json:"completedOn,omitempty"`
	Tags                      map[string]string `json:"tags,omitempty"`
	SortId                    int64             `json:"-"` // unexported
}

func (p *Promise) String() string {
	return fmt.Sprintf(
		"Promise(id=%s, state=%s, param=%s, value=%s, timeout=%d, idempotencyKeyForCreate=%s, idempotencyKeyForUpdate=%s, tags=%s)",
		p.Id,
		p.State,
		p.Param,
		p.Value,
		p.Timeout,
		p.IdempotencyKeyForCreate,
		p.IdempotencyKeyForComplete,
		p.Tags,
	)
}

type State int

const (
	Pending State = 1 << iota
	Resolved
	Rejected
	Canceled
	Timedout
)

func (s State) String() string {
	switch s {
	case Pending:
		return "PENDING"
	case Resolved:
		return "RESOLVED"
	case Rejected:
		return "REJECTED"
	case Canceled:
		return "REJECTED_CANCELED"
	case Timedout:
		return "REJECTED_TIMEDOUT"
	default:
		panic("invalid state")
	}
}

func (s *State) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *State) UnmarshalJSON(data []byte) error {
	var state string
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	switch strings.ToUpper(state) {
	case "PENDING":
		*s = Pending
	case "RESOLVED":
		*s = Resolved
	case "REJECTED":
		*s = Rejected
	case "REJECTED_CANCELED":
		*s = Canceled
	case "REJECTED_TIMEDOUT":
		*s = Timedout
	default:
		return fmt.Errorf("invalid state '%s'", state)
	}

	return nil
}

func (s State) In(mask State) bool {
	return s&mask != 0
}

type Value struct {
	Headers map[string]string `json:"headers,omitempty"`
	Data    []byte            `json:"data,omitempty"`
}

func (v Value) String() string {
	return fmt.Sprintf(
		"Value(headers=%s, data=%s)",
		v.Headers,
		string(v.Data),
	)
}
