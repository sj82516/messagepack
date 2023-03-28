package msgpack

import "fmt"

func example() {
    var b []byte
    json := map[string]interface{}{"foo": "bar"}
    b, err := Encode(json)
    if err != nil {
        panic(err)
    }
    
    // byte array of messagepack format
    // [0x82 0xa3 0x66 0x6f 0x6f 0xa3 0x62 0x61 0x72]
    fmt.Println(b)
}
