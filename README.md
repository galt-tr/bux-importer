# Import CLI Tool

This CLI tool will derive a set of addresses for a given xpub and import any
found utxos for that address. 

## Prerequisites
By default, this importer uses the WhatsOnChain API to grab transactions and
utxos. To increase performance, it's recommended to get an API key to increase
the default rate limit. Then, set the following env var:
```
$ export WOC_API_KEY=<api_key>
```

## Build CLI
```
$ go build
```

## Usage
```
$ ./importer -depth=20 -gap-limit=10 -debug=true <raw_xpub>
```

## Flags
```
| flag      | type | description                | default |
| --------- | ---- | -------------------------- | ------- |
| depth     | int  | depth of keys to derive    | 20      |
| gap-limit | int  | gap limit to stop deriving | 10      |
| debug     | bool | enable debug logging       | false   |
```
