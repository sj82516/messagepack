package msgpack

import (
    "bytes"
    "testing"
)

type testCase struct {
    dest string
    src  interface{}
    exp  []byte
}

var intTestCases = []testCase{
    {
        dest: "int",
        src:  0,
        exp:  []byte{0x00},
    },
    {
        dest: "int",
        src:  1,
        exp:  []byte{0x01},
    },
    {
        dest: "int",
        src:  127,
        exp:  []byte{0x7f},
    },
    {
        dest: "int",
        src:  128,
        exp:  []byte{0xcc, 0x80},
    },
    {
        dest: "int",
        src:  255,
        exp:  []byte{0xcc, 0xff},
    },
    {
        dest: "int",
        src:  256,
        exp:  []byte{0xcd, 0x01, 0x00},
    },
    {
        dest: "int",
        src:  65535,
        exp:  []byte{0xcd, 0xff, 0xff},
    },
    {
        dest: "int",
        src:  65536,
        exp:  []byte{0xce, 0x00, 0x01, 0x00, 0x00},
    },
    {
        dest: "int",
        src:  4294967295,
        exp:  []byte{0xce, 0xff, 0xff, 0xff, 0xff},
    },
    {
        dest: "int",
        src:  4294967296,
        exp:  []byte{0xcf, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
    },
}

func TestEncode(t *testing.T) {
    test := []testCase{
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
    test = append(test, intTestCases...)
    
    for _, v := range test {
        if act, err := Encode(v.src); err != nil {
            t.Errorf("Error: %s", err)
        } else if !bytes.Equal(act, v.exp) {
            t.Errorf("Expected %v, got %v", v.exp, act)
        }
    }
}
