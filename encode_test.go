package msgpack

import (
    "bytes"
    "testing"
)

func TestEncode(t *testing.T) {
    test := []struct {
        dest string
        src  interface{}
        exp  []byte
    }{
        {
            dest: "nil",
            src:  nil,
            exp:  []byte{0xc0},
        },
        {
            dest: "bool",
            src:  true,
            exp:  []byte{0xc3},
        },
        {
            dest: "bool",
            src:  false,
            exp:  []byte{0xc2},
        },
    }
    
    for _, v := range test {
        if act, err := Encode(v.src); err != nil {
            t.Errorf("Error: %s", err)
        } else if !bytes.Equal(act, v.exp) {
            t.Errorf("Expected %v, got %v", v.exp, act)
        }
    }
}
