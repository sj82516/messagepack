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
