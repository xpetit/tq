# tq

`tq` filters out empty objets and arrays from JSON data. It also filters empty strings (configurable) and shows one value per array (configurable).

## Installation

```
go install github.com/xpetit/tq@latest
```

## Usage

```
cat file.json | jq
```

> ```json
> {
>   "a": 0,
>   "b": null,
>   "c": "",
>   "d": {},
>   "e": [],
>   "f": [
>     {},
>     {
>       "a": [1, 2, 3]
>     }
>   ],
>   "g": [[null, "", 4], 5]
> }
> ```

```
tq file.json | jq
```

> ```json
> {
>   "a": 0,
>   "f": [
>     {
>       "a": [1]
>     }
>   ],
>   "g": [[4]]
> }
> ```

```
tq -empty -limit=2 file.json | jq
```

> ```json
> {
>   "a": 0,
>   "c": "",
>   "f": [
>     {
>       "a": [1, 2]
>     }
>   ],
>   "g": [["", 4], 5]
> }
> ```
