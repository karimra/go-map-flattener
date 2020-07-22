package transform

import (
	"encoding/json"
	"testing"
)

func TestFlatten(t *testing.T) {
	in := map[string]interface{}{
		"1": -1,
		"2": map[string]interface{}{
			"k1": "v1",
			"k2": "v2",
			"k3": 1,
			"k4": []int{1, 2, 3, 4},
			"k5": map[string]interface{}{
				"kk1": "vv1",
			},
		},
		"3": true,
		"4": 0.1,
		"5": map[float64]interface{}{
			0.1: 1,
		},
	}
	f := NewFlattener()
	out, err := f.Flatten(in)
	if err != nil {
		t.Errorf("flatten fail: %v", err)
	}
	t.Logf("%+v", out)
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Errorf("marshal fail: %v", err)
	}
	t.Log(string(b))
}
