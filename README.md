Cayenne
-------

A library to support [Cayenne binary encoding](https://developers.mydevices.com/cayenne/docs/lora/#lora-cayenne-low-power-payload)

It borrows all the decoding code from https://github.com/TheThingsNetwork/go-cayenne-lib, adding some helpers functions, goish style and go mod.

```
e := cayenne.NewEncoder()
e.AddGPS(1, 48.8, 2.2, 100.0)

b := e.Bytes()
hexPayload := hex.EncodeToString(b)
fmt.Println("Data", hexPayload)
```