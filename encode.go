// msgpack - Go implementation of MessagePack

package msgpack

import (
    "errors"
    "math"
    "reflect"
    "sort"
)

const (
    Max4ByteLen  = 16
    Max16ByteLen = 65536
    Max32ByteLen = 4294967296
    Max64ByteLen = uint(18446744073709551615)
)

// ErrOutOfRange is returned when the size of value is out of range.
var ErrOutOfRange = errors.New("out of range")

type encoder func(interface{}) ([]byte, error)

// Encode encodes the given value into messagepack format byte arrays.
// The value must be a struct, map, slice, or array.
func Encode(v interface{}) ([]byte, error) {
    return encode(v)
}

func encode(v interface{}) ([]byte, error) {
    encoderCh := []encoder{
        encodeNil,
        encodeBool,
        encodeInt,
        encodeUint,
        encodeFloat,
        encodeBin,
        encodeStr,
        encodeArr,
        encodeMap,
    }
    
    for _, encoder := range encoderCh {
        if b, err := encoder(v); err != nil {
            return nil, err
        } else if b != nil {
            return b, nil
        }
    }
    
    return nil, nil
}

// https://github.com/msgpack/msgpack/blob/master/spec.md#nil-format
func encodeNil(v interface{}) ([]byte, error) {
    if v != nil {
        return nil, nil
    }
    
    return []byte{0xc0}, nil
}

// https://github.com/msgpack/msgpack/blob/master/spec.md#bool-format-family
func encodeBool(v interface{}) ([]byte, error) {
    if t, ok := v.(bool); !ok {
        return nil, nil
    } else {
        if t {
            return []byte{0xc3}, nil
        }
        
        return []byte{0xc2}, nil
    }
}

// https://github.com/msgpack/msgpack/blob/master/spec.md#int-format-family
func encodeInt(v interface{}) ([]byte, error) {
    if reflect.TypeOf(v).Kind() != reflect.Int &&
        reflect.TypeOf(v).Kind() != reflect.Int8 &&
        reflect.TypeOf(v).Kind() != reflect.Int16 &&
        reflect.TypeOf(v).Kind() != reflect.Int32 &&
        reflect.TypeOf(v).Kind() != reflect.Int64 {
        return nil, nil
    }
    
    i := v.(int)
    if i >= 0 {
        uI := uint(i)
        return encodeUint(uI)
    }
    return encodeNegInt(i)
}

func encodeUint(v interface{}) ([]byte, error) {
    if reflect.TypeOf(v).Kind() != reflect.Uint &&
        reflect.TypeOf(v).Kind() != reflect.Uint8 &&
        reflect.TypeOf(v).Kind() != reflect.Uint16 &&
        reflect.TypeOf(v).Kind() != reflect.Uint32 &&
        reflect.TypeOf(v).Kind() != reflect.Uint64 {
        return nil, nil
    }
    
    uInt, _ := v.(uint)
    if uInt < 128 {
        return []byte{byte(uInt)}, nil
    }
    if uInt < 256 {
        return []byte{0xcc, byte(uInt)}, nil
    }
    if uInt < Max16ByteLen {
        return []byte{0xcd, byte(uInt >> 8), byte(uInt)}, nil
    }
    if uInt < Max32ByteLen {
        return []byte{0xce, byte(uInt >> 24), byte(uInt >> 16), byte(uInt >> 8), byte(uInt)}, nil
    }
    if uInt <= Max64ByteLen {
        return []byte{0xcf, byte(uInt >> 56), byte(uInt >> 48), byte(uInt >> 40), byte(uInt >> 32), byte(uInt >> 24), byte(uInt >> 16), byte(uInt >> 8), byte(uInt)}, nil
    }
    
    return nil, ErrOutOfRange
}

func encodeNegInt(v interface{}) ([]byte, error) {
    nI := v.(int)
    if nI >= -32 {
        return []byte{byte(nI)}, nil
    }
    if nI >= -128 {
        return []byte{0xd0, byte(nI)}, nil
    }
    if nI >= -32768 {
        return []byte{0xd1, byte(nI >> 8), byte(nI)}, nil
    }
    if nI >= -2147483648 {
        return []byte{0xd2, byte(nI >> 24), byte(nI >> 16), byte(nI >> 8), byte(nI)}, nil
    }
    if nI >= -9223372036854775808 {
        return []byte{0xd3, byte(nI >> 56), byte(nI >> 48), byte(nI >> 40), byte(nI >> 32), byte(nI >> 24), byte(nI >> 16), byte(nI >> 8), byte(nI)}, nil
    }
    return nil, nil
}

// https://github.com/msgpack/msgpack/blob/master/spec.md#float-format-family
func encodeFloat(v interface{}) ([]byte, error) {
    if reflect.TypeOf(v).Kind() != reflect.Float32 &&
        reflect.TypeOf(v).Kind() != reflect.Float64 {
        return nil, nil
    }
    
    switch v.(type) {
    case float32:
        x := v.(float32)
        bits := math.Float32bits(x)
        bytes := make([]byte, 4)
        for i := 0; i < 4; i++ {
            bytes[3-i] = byte(bits >> uint(8*i))
        }
        bytes = append([]byte{0xca}, bytes...)
        return bytes, nil
    case float64:
        x := v.(float64)
        bits := math.Float64bits(x)
        bytes := make([]byte, 8)
        for i := 0; i < 8; i++ {
            bytes[7-i] = byte(bits >> uint(8*i))
        }
        bytes = append([]byte{0xcb}, bytes...)
        return bytes, nil
    }
    return nil, nil
}

// https://github.com/msgpack/msgpack/blob/master/spec.md#str-format-family
func encodeStr(v interface{}) ([]byte, error) {
    if reflect.TypeOf(v).Kind() != reflect.String {
        return nil, nil
    }
    
    s := v.(string)
    head := []byte{}
    bytes := []byte(s)
    if len(s) < 32 {
        head = []byte{0xa0 + byte(len(s))}
    } else if len(s) < 256 {
        head = []byte{0xd9, byte(len(s))}
    } else if len(s) < Max16ByteLen {
        head = []byte{0xda, byte(len(s) >> 8), byte(len(s))}
    } else if len(s) < Max32ByteLen {
        head = []byte{0xdb, byte(len(s) >> 24), byte(len(s) >> 16), byte(len(s) >> 8), byte(len(s))}
    } else {
        return nil, ErrOutOfRange
    }
    
    return append(head, bytes...), nil
}

// https://github.com/msgpack/msgpack/blob/master/spec.md#map-format-family
func encodeBin(v interface{}) ([]byte, error) {
    if reflect.TypeOf(v).Kind() != reflect.Slice && reflect.TypeOf(v).Kind() != reflect.Array {
        return nil, nil
    }
    // # make sure is byte array
    if reflect.TypeOf(v).Elem().Kind() != reflect.Uint8 {
        return nil, nil
    }
    
    s := reflect.ValueOf(v)
    if s.Len() == 0 {
        return []byte{0xc4, 0}, nil
    }
    
    head := []byte{}
    if s.Len() < 256 {
        head = []byte{0xc4, byte(s.Len())}
    } else if s.Len() < Max16ByteLen {
        head = []byte{0xc5, byte(s.Len() >> 8), byte(s.Len())}
    } else if s.Len() < Max32ByteLen {
        head = []byte{0xc6, byte(s.Len() >> 24), byte(s.Len() >> 16), byte(s.Len() >> 8), byte(s.Len())}
    } else {
        return nil, ErrOutOfRange
    }
    
    return append(head, s.Bytes()...), nil
}

// https://github.com/msgpack/msgpack/blob/master/spec.md#array-format-family
func encodeArr(v interface{}) ([]byte, error) {
    if reflect.TypeOf(v).Kind() != reflect.Slice && reflect.TypeOf(v).Kind() != reflect.Array {
        return nil, nil
    }
    
    s := reflect.ValueOf(v)
    head := []byte{}
    if s.Len() < Max4ByteLen {
        head = []byte{0x90 + byte(s.Len())}
    } else if s.Len() < Max16ByteLen {
        head = []byte{0xdc, byte(s.Len() >> 8), byte(s.Len())}
    } else if s.Len() < Max32ByteLen {
        head = []byte{0xdd, byte(s.Len() >> 24), byte(s.Len() >> 16), byte(s.Len() >> 8), byte(s.Len())}
    } else {
        return nil, ErrOutOfRange
    }
    
    body := []byte{}
    for i := 0; i < s.Len(); i++ {
        b, err := encode(s.Index(i).Interface())
        if err != nil {
            return nil, err
        }
        body = append(body, b...)
    }
    
    return append(head, body...), nil
}

// https://github.com/msgpack/msgpack/blob/master/spec.md#map-format-family
func encodeMap(v interface{}) ([]byte, error) {
    if reflect.TypeOf(v).Kind() != reflect.Map {
        return nil, nil
    }
    
    m := reflect.ValueOf(v)
    head := []byte{}
    if m.Len() < Max4ByteLen {
        head = []byte{0x80 + byte(m.Len())}
    } else if m.Len() < Max16ByteLen {
        head = []byte{0xde, byte(m.Len() >> 8), byte(m.Len())}
    } else if m.Len() < Max32ByteLen {
        head = []byte{0xdf, byte(m.Len() >> 24), byte(m.Len() >> 16), byte(m.Len() >> 8), byte(m.Len())}
    } else {
        return nil, ErrOutOfRange
    }
    
    body := []byte{}
    keys := m.MapKeys()
    // json key must be string type
    sort.Slice(keys, func(i, j int) bool {
        return keys[i].String() < keys[j].String()
    })
    for _, key := range keys {
        b, err := encode(key.Interface())
        if err != nil {
            return nil, err
        }
        
        body = append(body, b...)
        b, err = encode(m.MapIndex(key).Interface())
        if err != nil {
            return nil, err
        }
        body = append(body, b...)
    }
    
    return append(head, body...), nil
}
