package msgpack

import (
    "errors"
    "reflect"
)

type encoder func(interface{}) ([]byte, error)

// Encode encodes the given value into messagepack format byte arrays.
// The value must be a struct, map, slice, or array.
func Encode(v interface{}) ([]byte, error) {
    encoderCh := []encoder{
        encodeNil,
        encodeBool,
        encodeInt,
        encodeUint,
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

// nil
func encodeNil(v interface{}) ([]byte, error) {
    if v != nil {
        return nil, nil
    }
    
    return []byte{0xc0}, nil
}

// bool
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
    return nil, nil
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
    if uInt < 65536 {
        return []byte{0xcd, byte(uInt >> 8), byte(uInt)}, nil
    }
    if uInt < 4294967296 {
        return []byte{0xce, byte(uInt >> 24), byte(uInt >> 16), byte(uInt >> 8), byte(uInt)}, nil
    }
    if uInt < 18446744073709551615 {
        return []byte{0xcf, byte(uInt >> 56), byte(uInt >> 48), byte(uInt >> 40), byte(uInt >> 32), byte(uInt >> 24), byte(uInt >> 16), byte(uInt >> 8), byte(uInt)}, nil
    }
    
    return nil, errors.New("out of range")
}
