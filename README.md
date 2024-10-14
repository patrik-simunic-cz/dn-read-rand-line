
# DNAnexus / Task 1 - read line

## Build locally

```bash
git clone https://github.com/patrik-simunic-cz/dn-read-rand-line
cd dn-read-rand-line

go build .
```

## Example usage

### Read line (basic)

Command: `readline <TXT_FILE_PATH> <LINE_INDEX>`

> [!CAUTION]
> **The `<LINE_INDEX>` argument uses 0-based indexing (ie. line index 3 = line 4).**<br>
> _This is unlike the advanced usage command ([`readline read`](#read-line-advanced)), which expectes the actual line number._

#### Example

```bash
./readline ./data/large_file.txt 9876
```

### Generate file

Command: `readline generate`

#### Flags

|Flag|Required|Description|
|-|-|-|
|`--path`|**Yes**|Path to the generated file|
|`--lines`|No|Number of generated lines|
|`--wordsPerLine`|No|Max words per generated line|

#### Example

```bash
./readline generate --path ./data/large_file.txt --lines 10000 --wordsPerLine 20
```

### Read line (advanced)

Command: `readline rand`

#### Flags

|Flag|Required|Description|
|-|-|-|
|`--verbose`|No|Print execution details and statistics|
|`--path`|**Yes**|Path to the file to read|
|`--indexPath`|No|Path to the index file|
|`--line`|**Yes**|Line to print|

#### Example

```bash
./readline read --verbose --path ./data/large_file.txt --indexPath ./index.idx --line 9876
```
