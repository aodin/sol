package sol

import (
	"encoding/json"
	"testing"
)

func TestValues(t *testing.T) {
	values := Values{"c": []byte("bytes")}
	b, err := json.Marshal(values)
	if err != nil {
		t.Fatalf("Unexpected error while marshaling: %s", err)
	}
	if string(b) != `{"c":"bytes"}` {
		t.Errorf(`Unexpected JSON marshal: %s != {"c":"bytes"}`, b)
	}
}
