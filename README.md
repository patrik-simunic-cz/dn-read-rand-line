
# DNAnexus / Task 1 - read random line

## Build locally

```bash
git clone https://github.com/patrik-simunic-cz/dn-read-rand-line
cd dn-read-rand-line

go build .
```

## Example usage

### Generate file

Command: `readline generate`

|Flag|Required|Description|
|-|-|-|
|`--path`|**Yes**|Path to the generated file|
|`--lines`|No|Number of generated lines|
|`--wordsPerLine`|No|Max words per generated line|

```bash
./readline generate --path ./data/large_file.txt --lines 10000 --wordsPerLine 20
```

### Read random line

Command: `readline rand`

|Flag|Required|Description|
|-|-|-|
|`--verbose`|No|Print execution details and statistics|
|`--path`|**Yes**|Path to the file to read|
|`--indexPath`|No|Path to the index file|

```bash
./readline rand --verbose --path ./data/large_file.txt --indexPath ./index.idx
```
