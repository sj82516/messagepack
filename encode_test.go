package msgpack

import (
    "bytes"
    "strings"
    "testing"
)

type testCase struct {
    dest    string
    src     interface{}
    srcFunc func() interface{}
    exp     []byte
    expFunc func() []byte
}

var nilTestCases = []testCase{
    {
        dest: "nil",
        src:  nil,
        exp:  []byte{0xc0},
    },
}
var boolTestCases = []testCase{
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
var uintTestCases = []testCase{
    {
        dest: "uint",
        src:  uint(0),
        exp:  []byte{0x00},
    },
    {
        dest: "uint",
        src:  uint(1),
        exp:  []byte{0x01},
    },
    {
        dest: "uint",
        src:  uint(127),
        exp:  []byte{0x7f},
    },
    {
        dest: "uint",
        src:  uint(128),
        exp:  []byte{0xcc, 0x80},
    },
    {
        dest: "uint",
        src:  uint(255),
        exp:  []byte{0xcc, 0xff},
    },
    {
        dest: "uint",
        src:  uint(256),
        exp:  []byte{0xcd, 0x01, 0x00},
    },
    {
        dest: "uint",
        src:  uint(65535),
        exp:  []byte{0xcd, 0xff, 0xff},
    },
    {
        dest: "uint",
        src:  uint(65536),
        exp:  []byte{0xce, 0x00, 0x01, 0x00, 0x00},
    },
    {
        dest: "uint",
        src:  uint(4294967295),
        exp:  []byte{0xce, 0xff, 0xff, 0xff, 0xff},
    },
    {
        dest: "uint",
        src:  uint(4294967296),
        exp:  []byte{0xcf, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
    },
    {
        dest: "uint",
        src:  uint(18446744073709551615),
        exp:  []byte{0xcf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
    },
}
var negIntTestCases = []testCase{
    {
        dest: "negative int",
        src:  -1,
        exp:  []byte{0xff},
    },
    {
        dest: "negative int",
        src:  -32,
        exp:  []byte{0xe0},
    },
    {
        dest: "negative int",
        src:  -33,
        exp:  []byte{0xd0, 0xdf},
    },
    {
        dest: "negative int",
        src:  -128,
        exp:  []byte{0xd0, 0x80},
    },
    {
        dest: "negative int",
        src:  -129,
        exp:  []byte{0xd1, 0xff, 0x7f},
    },
    {
        dest: "negative int",
        src:  -32768,
        exp:  []byte{0xd1, 0x80, 0x00},
    },
    {
        dest: "negative int",
        src:  -32769,
        exp:  []byte{0xd2, 0xff, 0xff, 0x7f, 0xff},
    },
    {
        dest: "negative int",
        src:  -2147483648,
        exp:  []byte{0xd2, 0x80, 0x00, 0x00, 0x00},
    },
    {
        dest: "negative int",
        src:  -2147483649,
        exp:  []byte{0xd3, 0xff, 0xff, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xff},
    },
    {
        dest: "negative int",
        src:  -9223372036854775808,
        exp:  []byte{0xd3, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
    },
}
var floatTestCases = []testCase{
    {
        dest: "float",
        src:  float32(0),
        exp:  []byte{0xca, 0x00, 0x00, 0x00, 0x00},
    },
    {
        dest: "float",
        src:  float32(0.1),
        exp:  []byte{0xca, 0x3d, 0xcc, 0xcc, 0xcd},
    },
    {
        dest: "float",
        src:  float64(0.1),
        exp:  []byte{0xcb, 0x3f, 0xb9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a},
    },
}
var strTestCases = []testCase{
    {
        dest: "string",
        src:  "",
        exp:  []byte{0xa0},
    },
    {
        dest: "string",
        src:  "a",
        exp:  []byte{0xa1, 0x61},
    },
    {
        dest: "string",
        src:  "abc",
        exp:  []byte{0xa3, 0x61, 0x62, 0x63},
    },
    {
        dest: "string",
        src:  "Hello World!",
        exp:  []byte{0xac, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21},
    },
    {
        dest: "string",
        src:  strings.Repeat("b", 32),
        exp:  append([]byte{0xd9, 0x20}, []byte(strings.Repeat("b", 32))...),
    },
    {
        dest: "string",
        src:  strings.Repeat("c", 256),
        exp:  append([]byte{0xda, 0x01, 0x00}, []byte(strings.Repeat("c", 256))...),
    },
    {
        dest: "string",
        src:  strings.Repeat("d", 65536),
        exp:  append([]byte{0xdb, 0x00, 0x01, 0x00, 0x00}, []byte(strings.Repeat("d", 65536))...),
    },
    {
        dest: "string",
        src:  "你好",
        exp:  []byte{0xa6, 0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd},
    },
}
var arrTestCases = []testCase{
    {
        dest: "array",
        src:  []interface{}{},
        exp:  []byte{0x90},
    },
    {
        dest: "array",
        src:  []interface{}{1, 2, 3},
        exp:  []byte{0x93, 0x01, 0x02, 0x03},
    },
    {
        dest: "array",
        src:  []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
        exp:  []byte{0xdc, 0x00, 0x10, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
    },
    {
        dest: "array",
        src:  []interface{}{1, 2, 3, "Hello World!"},
        exp:  []byte{0x94, 0x01, 0x02, 0x03, 0xac, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21},
    },
    {
        dest: "array",
        src:  []interface{}{1, 2, 3, "Hello World!", []interface{}{1, 2, 3}},
        exp:  []byte{0x95, 0x01, 0x02, 0x03, 0xac, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21, 0x93, 0x01, 0x02, 0x03},
    },
    {
        dest: "array",
        src:  strings.Split(strings.Repeat("a", 100), ""),
        expFunc: func() []byte {
            var strA, _ = encodeStr("a")
            var bytes []byte
            bytes = append(bytes, 0xdc, 0x00, 0x64)
            for i := 0; i < 100; i++ {
                bytes = append(bytes, strA...)
            }
            return bytes
        },
    }, {
        dest: "array",
        src:  strings.Split(strings.Repeat("a", 65536), ""),
        expFunc: func() []byte {
            var strA, _ = encodeStr("a")
            var bytes []byte
            bytes = append(bytes, 0xdd, 0x00, 0x01, 0x00, 0x00)
            for i := 0; i < 65536; i++ {
                bytes = append(bytes, strA...)
            }
            return bytes
        },
    },
}
var mapTestCases = []testCase{
    {
        dest: "map",
        src:  map[string]interface{}{},
        exp:  []byte{0x80},
    },
    {
        dest: "map",
        src:  map[string]interface{}{"a": 1, "b": 2, "c": 3},
        exp:  []byte{0x83, 0xa1, 0x61, 0x01, 0xa1, 0x62, 0x02, 0xa1, 0x63, 0x03},
    },
    {
        dest: "map",
        src:  map[string]interface{}{"a": 1, "b": 2, "c": map[string]interface{}{"d": 4, "e": 5, "f": 6}},
        exp:  []byte{0x83, 0xa1, 0x61, 0x01, 0xa1, 0x62, 0x02, 0xa1, 0x63, 0x83, 0xa1, 0x64, 0x04, 0xa1, 0x65, 0x05, 0xa1, 0x66, 0x06},
    },
}

func TestEncode(t *testing.T) {
    var tests []testCase
    tests = append(tests, nilTestCases...)
    tests = append(tests, boolTestCases...)
    tests = append(tests, intTestCases...)
    tests = append(tests, uintTestCases...)
    tests = append(tests, negIntTestCases...)
    tests = append(tests, floatTestCases...)
    tests = append(tests, strTestCases...)
    tests = append(tests, arrTestCases...)
    tests = append(tests, mapTestCases...)
    
    for _, v := range tests {
        src := v.src
        if v.srcFunc != nil {
            src = v.srcFunc()
        }
        
        act, err := Encode(src)
        exp := v.exp
        if v.expFunc != nil {
            exp = v.expFunc()
        }
        
        if err != nil {
            t.Errorf("Error: %s", err)
        } else if !bytes.Equal(act, exp) {
            t.Errorf("Expected %v", exp)
            t.Errorf("Actual   %v", act)
        }
    }
}
