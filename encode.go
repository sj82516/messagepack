package msgpack

import "reflect"

// Encode encodes the given value into messagepack format byte arrays.
// The value must be a struct, map, slice, or array.
func Encode(v interface{}) ([]byte, error) {
    if v == nil {
        return encodeNil()
    }
    
    if reflect.TypeOf(v) == reflect.TypeOf(true) {
        return encodeBool(v.(bool))
    }
    
    if reflect.TypeOf(v) == reflect.TypeOf(int(0)) {
        return encodeInt(v.(int))
    }
    
    return nil, nil
}

// nil
func encodeNil() ([]byte, error) {
    return []byte{0xc0}, nil
}

// bool
func encodeBool(v bool) ([]byte, error) {
    if v {
        return []byte{0xc3}, nil
    }
    
    return []byte{0xc2}, nil
}

func encodeInt(v int) ([]byte, error) {
    if v >= 0 {
        return encodeUint(uint(v))
    }
    return nil, nil
    //return encodeSignInt(v)
}

// https://github.com/msgpack/msgpack/blob/master/spec.md#int-format-family
func encodeUint(v uint) ([]byte, error) {
    if v <= 127 {
        return []byte{byte(v)}, nil
    }
    if v <= 255 {
        return []byte{0xcc, byte(v)}, nil
    }
    if v <= 65535 {
        return []byte{0xcd, byte(v >> 8), byte(v)}, nil
    }
    if v <= 4294967295 {
        return []byte{0xce, byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
    }
    return []byte{0xcf, byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
}
