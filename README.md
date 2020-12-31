## Overview

Archivist generates data structure definitions from JSON files for any kind of programming language. Besides that, it also provides a library for golang to cut the job of managing runtime configs down to several lines of code.

## Contents

- [Features](#features)
- [Getting Started](#install)
- [Basic Usage Example](#example)
- [Meta File](#.meta)
- Extended Types: [reference](#reference), [datetime](#datetime), [duration](#duration), [i18n](#i18n)
- [Config Group and Subgroup](#conf-group)
- [Overwriting Mechanism](#overwrite)
- [Fast Loading](#easyjson)
- [Hot Reload](#reload)
- [Hot Patching](#patch)
- [Config Extension](#extension)
- [Backward Compatibility Check](#compatibility)
- [Whitelist vs Blacklist](#whitelist-blacklist)
- [Use Javascript as Data File](#.js)
- [Subcommand: paths](#paths)
- [Subcommand: orphan](#orphan)
- [Subcommand: tpls](#tpls)
- [Generate Data Structure Definitions for Another Programming Language](#code-templates)
- [FAQ](#faq)

<a name="features"></a>
## Features

#### Code Generator

* Guess out the most suitable data type for arbitrary JSON field
* Customize data type via [.meta](#.meta) file. E.g. replace struct with map[string]...
* Support extended data types via [.meta](#.meta) file, including [reference](#reference), [datetime](#datetime), [duration](#duration) and [i18n](#i18n)
* Process [.js](#.js) file in addition to .json file, which is much friendly to edit for human-kind
* Show all node paths and data types with the '[paths](#paths)' subcommand
* Detect and drop orphan generated files with the '[orphan](#orphan)' subcommand
* Customize [code templates](#code-templates), by which any kind of programming language can benefit from the tool

#### Library

* Load all .json files with only [several lines of code](#load)
* Select config [group and subgroup](#conf-group) at runtime
* Support all kinds of config [overwriting](#overwrite) manners: file-level, content-level, in-memory, etc.
* Support extended data types, including [reference](#reference), [datetime](#datetime), [duration](#duration) and [i18n](#i18n)
* Support [hot reload](#reload) (atomically switch to a whole new config collection)
* Support [hot patching](#patch)
* Support [backward compatibility check](#compatibility)
* Support [whitelist and blacklist](#whitelist-blacklist)
* Support [config extension](#extension), which stores organized configs
* Thread-safe

<a name="install"></a>
## Getting Started

Note: To use [.js](#.js) as config file, you must have <a href="https://nodejs.org" target="_blank">node.js</a> installed.

#### Install from Binary Release

Download a binary release suitable for your system, unzip it, and then update your PATH environment variable in order that archivist can run from anywhere.

#### Install from Source

``` bash
go get -u github.com/kingsgroupos/archivist/cli/archivist
go get -u github.com/edwingeng/easyjson-alt/easyjson

# To use .js as config file, clone this repository, and then copy cli/script/* to somewhere.
# These scripts are used to convert .js files to .json files.
```

<a name="example"></a>
## Basic Usage Example

#### Initial State of the `foo` Directory

```
foo
├── conf
└── json
    ├── a.json
    └── b.json
```

Content of `a.json`:

``` json
{
    "A1": {
        "A_Struct": {
            "A_Int": 100,
            "A_Float": 1.23,
            "A_Bool": false,
            "A_String": "hello world",
            "A_Array": [
                1,
                2,
                3.4
            ],
            "A_Map": {
                "1": "hello",
                "2": "world"
            }
        }
    }
}
```

Content of `b.json`:

``` json
[
    {
        "Name": "Edwin",
        "Age": 100
    },
    {
        "Name": "Alisa",
        "Age": 20
    }
]
```

#### Generate Data Structure Definitions:

``` bash
archivist generate --outputDir=foo/conf --pkg=conf 'foo/json/*.json'
```

Results:

```
foo
├── conf
│   ├── aConf.go
│   ├── bConf.go
│   ├── collection.go
│   └── collectionExtension.go
└── json
    ├── a.json
    └── b.json
```

Code snippet of `aConf.go`:

``` go
package conf

type AConf struct {
    A1 *AConf_1910679952 `json:"A1" bson:"A1"`
}

// AConf_1910679952 represents /A1
type AConf_1910679952 struct {
    AStruct *AConf_448491524 `json:"A_Struct" bson:"A_Struct"`
}

// AConf_448491524 represents /A1/A_Struct
type AConf_448491524 struct {
    AArray  []float64        `json:"A_Array" bson:"A_Array"`
    ABool   bool             `json:"A_Bool" bson:"A_Bool"`
    AFloat  float64          `json:"A_Float" bson:"A_Float"`
    AInt    int64            `json:"A_Int" bson:"A_Int"`
    AMap    map[int64]string `json:"A_Map" bson:"A_Map"`
    AString string           `json:"A_String" bson:"A_String"`
}
```

Code snippet of `bConf.go`:

``` go
package conf

type BConf []*BConf_1771430469

// BConf_1771430469 represents /[]
type BConf_1771430469 struct {
    Age  int64  `json:"Age" bson:"Age"`
    Name string `json:"Name" bson:"Name"`
}
```

Besides `aConf.go` and `bConf.go`, archivist generates 2 more files, `collection.go` and `collectionExtension.go`, to simplify runtime config management.

Code snippet of `collection.go`:

``` go
package conf

type Collection struct {
    Extension     `json:"-"`
    filename2Conf map[string]interface{}

    AConf AConf
    BConf BConf
}
```

Code snippet of `collectionExtension.go`:

``` go
package conf

func (this *Collection) CompatibleVersions() []string {
    // return nil if you do not need backward compatibility check
    panic("not implemented yet")
}

type Extension struct {
}
```

<a name="load"></a>
#### Load config files

Note: It is **mandatory** to specify a [config group](#conf-group) with the `WATCHER_CONF_GROUP` environment variable.

Create a `local` subfolder under `foo/json` and put an empty file in it:

```
foo
├── conf
│   ├── aConf.go
│   ├── collection.go
│   └── collectionExtension.go
└── json
    ├── a.json
    └── local
        └── placeholder
```

Code:

``` go
// Do NOT forget to set WATCHER_CONF_GROUP=local before running the following code.

// Make sure WATCHER_CONF_GROUP is set
archivist.MustHaveConfGroup()

// Create an instance of archivist
ar := archivist.NewArchivist(archivist.WithRoot("foo/json"))

// Load config files
c1, err := ar.LoadCollection(conf.NewCollection)
if err != nil {
    panic(err)
}

// Set the current config collection
ar.SetCurrentCollection(c1)

// ...

// Get the current config collection
c2, err := ar.FindCollection("")
if err != nil {
    panic(err)
}

// Convert c2 to its actual type
c3 := c2.(*conf.Collection)

// Print out configs
fmt.Println(c3.AConf.A1.AStruct.AInt)
fmt.Println(c3.AConf.A1.AStruct.AMap)
fmt.Println(c3.BConf[0])

// At last, modify Collection.CompatibleVersions() to let it return nil
// ...
```

Output:

```
100
map[1:hello 2:world]
&{100 Edwin}
```

<a name=".meta"></a>
## Meta File

In most cases, archivist automatically determines the data type of each JSON field. But sometimes, the data type of guess is not desired, or you prefer to use an [extended data type](#extended), you may give a .meta file to customize that.

Suppose you have a file `abc.json`, and you want to change the data type of a JSON field. Just create a file `abc.meta` under the same directory where `abc.json` lies in. Archivist will parse it along with `abc.json` and apply the instructions in it.

#### Instruction Syntax of the Meta File

```
path type comment
```

- **path**: path of the field from root. It is not easy to spell out the path of a field, therefore archivist provides a subcommand, [paths](#paths), to make it easy, with which all you need to do is copy/paste.
- **type**: data type. Here is the legal types:
```
int, int8, int16, int32, int64, uint, uint8, uint16, uint32
float32, float64
string
bool
ref@...
datetime
duration
i18n
{}
[]
map[int], map[int8], map[int16], map[int32], map[int64]
map[uint], map[uint8], map[uint16], map[uint32]
map[string]
```
  - `ref@...` means [reference](#reference)
  - `{}` means struct
  - `[]` means array
  - `map[...]` means map, and its key type is `...`
- **comment**: comment for the field (optional)

#### Meta File Example
```
/[]/Age    int32    the number of years someone has lived
/[]/Name   string
```

#### The Two Towers: .meta and .suggested.meta

Sometimes, we need an extra .meta file controlled by someone else other than programmer. For this reason, .suggested.meta file is introduced in. It works in between type guessing and .meta file. In other words, .suggested.meta overrides type guessing and .meta overrides .suggested.meta.

#### More examples

<a href="https://github.com/kingsgroupos/archivist/tree/main/cli/archivist/example/json" target="_blank">https://github.com/kingsgroupos/archivist/tree/main/cli/archivist/example/json</a>

<a name="extended"></a>
## Extended Data Types

Archivist supports extended data types via [.meta](#.meta) file.

<a name="reference"></a>
#### reference

It works like the foreign key of relational database. In the following example, the type of `MyItem` is `ref@d`.

Content of `c.json`:

``` json
{
    "MyItem": 100
}
```

Content of `c.meta`:
```
/MyItem    ref@d
```

Content of `d.json`:
``` json
{
    "100": {
        "Value": 1000
    },
    "200": {
        "Value": 2000
    }
}
```

During code generation, archivist adds an extra field, `MyItemRef`, to `CConf`. During config loading, `MyItemRef` is bound automatically to its corresponding config item.

Code snippet of `cConf.go`:

``` go
type CConf struct {
    MyItem    int64      `json:"MyItem" bson:"MyItem"`
    MyItemRef *DConfItem `json:"-" bson:"-"`
}
```

Code snippet of `dConf.go`:

``` go
type DConf map[int64]*DConfItem

// DConfItem represents /map[]
type DConfItem struct {
    Value int64 `json:"Value" bson:"Value"`
}
```

The referenced config file must satisfy the following criterion: its corresponding type must be map, the key type must be int64
and the value type must be struct.

<a name="datetime"></a>
#### datetime

It helps you decode a datetime string that conforms to <a href="https://en.wikipedia.org/wiki/ISO_8601" target="_blank">ISO 8601</a> to a language specific data type, e.g. `time.Time` in golang. The following example demonstrates how to use the datetime type.

Content of `e.json`:

``` json
{
    "when": "2020-01-01T12:00:00Z"
}
```

Content of `e.meta`:

```
/when     datetime
```

Code snippet of `eConf.go`:

``` go
type EConf struct {
    When time.Time `json:"when" bson:"when"`
}
```

<a name="duration"></a>
#### duration

It helps you decode a duration string to a language specific data type, e.g. `wtime.Duration` in golang. A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h", "d". The following example demonstrates how to use the duration type.

Content of `f.json`:

``` json
{
    "lifetime": "10d12h15m5s"
}
```

Content of `f.meta`:

```
/lifetime   duration
```

Code snippet of `fConf.go`:

``` go
type FConf struct {
    Lifetime wtime.Duration `json:"lifetime" bson:"lifetime"`
}
```

<a name="i18n"></a>
#### i18n

It is an internationalization helper. Comparing with other extended types, i18n is a little complicated. The following example demonstrates how to use the i18n type.

Content of `g.json`:

``` json
{
    "Greetings": "KEY_GREETINGS",
    "GreetingsWithName": "XXX_GreetingsWithName",
    "Number": "KEY_NUMBER"
}
```

Content of `g.meta`:

```
/Greetings           i18n
/GreetingsWithName   i18n
/Number              i18n
```

Code snippet of `gConf.go`:

```
type GConf struct {
    Greetings         archivist.I18n `json:"Greetings" bson:"Greetings"`
    GreetingsWithName archivist.I18n `json:"GreetingsWithName" bson:"GreetingsWithName"`
    Number            archivist.I18n `json:"Number" bson:"Number"`
}
```

Code snippet of some `.go` file:

``` go
import "github.com/kingsgroupos/archivist/lib/go/archivist"

// ...

trans := map[string]map[string]string{
    "en-US": {
        "KEY_GREETINGS":         "Greetings!",
        "XXX_GreetingsWithName": "Greetings {0}!",
        "KEY_NUMBER":            "Number: {0:0.00}",
    },
    "fr-FR": {
        "KEY_GREETINGS":         "Bonjour!",
        "XXX_GreetingsWithName": "Bonjour {0}!",
    },
}

archivist.ResolveI18n = func(lang, key string) string {
    if m1 := trans[lang]; m1 != nil {
        if str, ok := m1[key]; ok {
            return str
        }
    }
    return key
}

fmt.Println(c3.GConf.Greetings.I18n("en-US"))
greetingsWithName, _ := c3.GConf.GreetingsWithName.Sprintf("fr-FR", "Edwin")
fmt.Println(greetingsWithName)
num, _ := c3.GConf.Number.Sprintf("en-US", 12345.6)
fmt.Println(num)
```

Output:

```
Greetings!
Bonjour Edwin!
Number: 12,345.60
```

<a name="conf-group"></a>
## Config Group and Subgroup

Different operating environments may need more or less different configs. The config files under the root directory are base configs. Config group and subgroup are designed to overwrite base configs at file-level and content-level. By convention, each subfolder under the root directory is a config group and each subfolder under a config group directory is a config subgroup. In the following example, _local_, _develop_ and _release_ are config groups, _subgroup1_ and _subgroup2_ are config subgroups.

```
bar
├── develop
├── local
└── release
    ├── subgroup1
    └── subgroup2
```

You can select a config group with the environment variable `WATCHER_CONF_GROUP` and a config subgroup with the environment variable `WATCHER_CONF_SUBGROUP`. Selecting a config group is always required while selecting a config subgroup is optional.

Read the following section to learn about config overwriting.

<a name="overwrite"></a>
## Overwriting Mechanism

- To overwrite a config file at file-level, put a file of the same name under your config group directory. Please note that config subgroup does not support file-level overwriting.
- To overwrite a config file at content-level, put a file of the same basename + [`.tweak.json`|`.local.json`] under your config group/subgroup directory. `.local.json` is designed for local changes, which should be ignored by your version control system.

In the example below, overwriting happens in the following sequence:

1. `release/a.json` overwrites `a.json` at file-level.
2. `release/a.tweak.json` overwrites the previous result at content-level.
3. `release/a.local.json` overwrites the previous result at content-level.

```
qux
├── a.json
└── release
    ├── a.json
    ├── a.local.json
    └── a.tweak.json
```

#### In-memory overwriting

Sometimes, config data comes from another source other than the file system. In such case, you can organize this kind of data with `archivist.Overwrite` and pass them to `LoadCollection`. `archivist.Overwrite` is always applied at the latest stage.

For more information, please refer to [Hot Patching](#patch).

#### Support of File-level & Content-level Overwriting

```
+-----------+------------+---------------+
|           | File-level | Content-level |
+-----------+------------+---------------+
| Group     |     Yes    |      Yes      |
+-----------+------------+---------------+
| Subgroup  |     No     |      Yes      |
+-----------+------------+---------------+
| In-memory |     Yes    |      Yes      |
+-----------+------------+---------------+
```

#### Overwriting Sequences

- Without in-memory file-level data:

```
. 1) Base:     -               - a.json
^ 2) Group     - File-level    - a.json
^ 3) Group     - Content-level - a.tweak.json
^ 4) Group     - Content-level - a.local.json
^ 5) Subgroup  - Content-level - a.tweak.json
^ 6) Subgroup  - Content-level - a.local.json
^ 7) In-memory - Content-level
```

- With in-memory file-level data:

```
. 1) Base      -               - a.json
^ 2) In-memory - File-level
^ 3) In-memory - Content-level
```

<a name="easyjson"></a>
## Fast Loading

<a href="https://github.com/edwingeng/easyjson-alt" target="_blank">easyjson</a> can dramatically increase the speed of config loading. Turn on `--x-easyjson` to earn the benefit.

```
archivist generate --outputDir=foo/conf --pkg=conf --x-easyjson 'foo/json/*.json'
```

<a name="reload"></a>
## Hot Reload

``` go
// Load config files
newColl, err := ar.LoadCollection(conf.NewCollection)
if err != nil {
    panic(err)
}

// Replace the current config collection
ar.SetCurrentCollection(newColl)
```

<a name="patch"></a>
## Hot Patching

With `PatchCollection`, you can apply a patch to an existing config collection. This function creates a new collection and reuses existing configs as much as possible. It does not call the reload callback, and it does not update the config [extension](#extension).

``` go
var overwrites []archivist.Overwrite
overwrites = append(overwrites,
    archivist.Overwrite{
        FileLevel: false,
        Target:    "a.json",
        Data: []byte(`{
            "A1": {
                "A_Struct": {
                    "A_Int": 200
                }
            }
        }`),
    })

newColl, err := ar.PatchCollection(existingColl, overwrites...)
if err != nil {
    panic(err)
}

ar.SetCurrentCollection(newColl)
```

<a name="extension"></a>
## Config Extension

Sometimes, raw configs does not fit your needs well. You have to write some code to organize raw configs and store the result for later use. Config extension is designed for that.

Perhaps you have found that archivist generated a new file named `collectionExtension.go` and added an `Extention` field for the `Collection` struct. After the first code generation, archivist will never touch `collectionExtension.go` again, so you can udpate the file freely without worrying about losing your code.

In short, you can organize raw configs with the reload callback, which is invoked each time after `LoadCollection` succeeds, and store the result into the `Extention` struct.

<a name="compatibility"></a>
## Backward Compatibility Check

Archivist provides a simple mechanism to check whether the app client needs to upgrade its configs. You can update the `CompatibleVersions` function in `collectionExtension.go` to return a list of versions that the collection is compatible with, and when you invoke `FindCollection`, pass the app client config version to the call. If the version is not compatible with the current collection, an error, `archivist.ErrUpgradeNeeded`, is returned.

<a name="whitelist-blacklist"></a>
## Whitelist vs Blacklist

- `WithWhitelist` helps you load only the specified config files.
- `WithBlacklist` helps you ignore the specified config files.

<a name=".js"></a>
## Use Javascript as Data File

Note: To use .js as config file, you must have <a href="https://nodejs.org" target="_blank">node.js</a> installed.

The .js config file must conform to the following format:

``` js
module.exports = {
    // ...
}
```

The archivist command supports .js file while the archivist library does not. Before loading configs, you must convert .js files to .json files and put them into proper directories. Fortunately, two scripts are made to do this. They lies in the `cli/script` directory.

It is highly recommended to have a look at <a href="https://github.com/kingsgroupos/archivist/tree/main/lib/go/overwrite" target="_blank">https://github.com/kingsgroupos/archivist/tree/main/lib/go/overwrite</a> to get more ideas.

<a name="paths"></a>
## Subcommand: paths

This subcommand shows all the node paths and data types of a .json/.js file.

Usage example:

```
### archivist paths <file> [flags]
>>> archivist paths foo/json/a.json

Path                       Type
/                          {}
/A1                        {}
/A1/A_Struct               {}
/A1/A_Struct/A_Array       []
/A1/A_Struct/A_Array/[]    float64
/A1/A_Struct/A_Bool        bool
/A1/A_Struct/A_Float       float64
/A1/A_Struct/A_Int         int64
/A1/A_Struct/A_Map         map[int64]
/A1/A_Struct/A_Map/map[]   string
/A1/A_Struct/A_String      string
```

<a name="orphan"></a>
## Subcommand: orphan

This subcommand finds out (and deletes) orphan files in your code.

```
### archivist orphan <dataDir> <codeDir> <codeFileExt> [ignore1,ignore2,...] [flags]
>>> archivist orphan foo/json foo/conf '.go' --delete

orphan deleted: foo/conf/xxx.go
```

<a name="tpls"></a>
## Subcommand: tpls

This subcommand outputs the built-in code templates.

```
### archivist tpls [flags]
>>> archivist tpls --outputDir tpls

tpls/struct.tpl
tpls/collection.tpl
tpls/collectionExtension.tpl
```

To use your own code templates, you must specify the --tplDir argument when generating code.

``` bash
archivist generate --outputDir=foo/conf --pkg=conf --tplDir=code/template/directory 'foo/json/*.json'
```

<a name="code-templates"></a>
## Generate Data Structure Definitions for Another Programming Language

1. Save the built-in code templates to files with the [tpls](#tpls) subcommand.
2. Modify code templates according to your needs. Syntax reference: <a href="https://godoc.org/text/template" target="_blank">https://godoc.org/text/template</a>.
3. Generate code with your own templates.

``` bash
# Set the file name extension of your code with --codeFileExt
# Set --x-collection=false if you do not need the collection file
# Set --x-collectionExtension=false if you do not need the collection extension file

archivist generate --outputDir=foo/conf --pkg=conf --tplDir=code/template/directory \
    --codeFileExt '.cs' 'foo/json/*.json'
```

Here is the structures and functions exposed to code templates by archivist:

``` go
const (
    ValueKind_Primitive ValueKind = 1
    ValueKind_Struct    ValueKind = 2
    ValueKind_Map       ValueKind = 3
    ValueKind_Array     ValueKind = 4
    ValueKind_Ref       ValueKind = 5
)

// Node contains the meta info of a JSON node.
type Node struct {
    Name   string
    Parent *Node

    ValueKind ValueKind
    Value     struct {
        Primitive    string
        StructFields []*Node
        MapKey       string
        MapValue     *Node
        ArrayValue   *Node
        RawRef       string
        Ref          string
    }

    Notes string
}

// structArgs is passed to struct.tpl when executing the code template.
type structArgs struct {
    Pkg   string
    Nodes []*Node
}

// collectionArgs is passed to collection.tpl when executing the code template.
type collectionArgs struct {
    Pkg                 string
    JSONFiles           []string
    Structs             []string
    RevRefGraph         map[string][]string
    CollectionExtension bool
}

// collectionExtensionArgs is passed to collectionExtension.tpl when executing the code template.
type collectionExtensionArgs struct {
    Pkg string
}

// bsonTag returns whether to generate BSON tag for struct field
func bsonTag() bool

// deepToRef scans the specified node and all of its descendents, if the type of any node is
// reference, it returns true. Otherwise, it returns false.
func deepToRef(node *Node) bool

// deepToStruct scans the specified node and all of its descendents, if the type of any node
// is struct, it returns true. Otherwise, it returns false.
func deepToStruct(node *Node) bool

// depth is an 'depth' manager.
// If action is "+", it increases 'depth' by 1 and returns an empty string.
// If action is "-", it decreases 'depth' by 1 and returns an empty string.
// If action is "0", it resets 'depth' to 0 and returns an empty string.
// If action is "v", it returns the current value of 'depth'.
func depth(action string) interface{}

// graveAccent always returns "`".
func graveAccent() string

// hasPrefix tests whether the string s begins with prefix.
func hasPrefix(s string, prefix string) bool

// hasSuffix tests whether the string s ends with suffix.
func hasSuffix(s string, suffix string) bool

// jsonFile returns the name of the JSON file being processed, e.g. "a.json".
func jsonFile() string

// lcfirst returns s with the first letter mapped to its lower case.
func lcfirst(s string) string

// lookupStructName returns the struct name of the specified node if its type is struct.
func lookupStructName(node *Node) string

// shortenRefName shortens refName by removing its suffix.
func shortenRefName(refName string) string

// stackPop removes a string from the top of the global stack and returns the string.
func stackPop() string

// stackPush pushes string s to the top of the global stack.
func stackPush(s string) string

// toCamel converts string s to the camel case.
func toCamel(s string) string

// toLower returns s with all Unicode letters mapped to their lower case.
func toLower(s string) string

// toPascal converts string s to the pascal case.
func toPascal(s string) string

// toUpper returns s with all Unicode letters mapped to their upper case.
func toUpper(s string) string

// trimPrefix returns s without the provided leading prefix string. If s
// doesn't start with prefix, s is returned unchanged.
func trimPrefix(s string, prefix string) string

// trimSuffix returns s without the provided trailing suffix string. If s
// doesn't end with suffix, s is returned unchanged.
func trimSuffix(s string, suffix string) string

// ucfirst returns s with the first letter mapped to its upper case.
func ucfirst(s string) string
```

<a name="faq"></a>
## FAQ

- **Why JSON? Why not YAML, TOML, XML, CSV, ...?**

JSON is simple yet good enough. It is widely supported by all kinds of programming languages. It supports hierarchical data structure. It provides no choice between field and property when you add a piece of data into a struct. JSON files are not hard to compare and merge.

JSON is not perfect. It has too many double quotation marks. It does not support comment. The last item in a {}/[] must have no trailing comma. That is why archivist also supports Javascript as an auxiliary format. The data definition syntax of Javascript is very close to JSON.

I'd like to support more file formats, but I have no time to do that. Contribution is welcome :)

- **The code templates are too sophisticated. How can I get started?**

The built-in code templates for golang are fully functional yet sophisticated. You can build your own templates in a much simpler way if you do not need so many features. We have other code templates written for `C#`, `Lua` and `EmmyLua`, but they are tightly coupled with our products and are not proper to open source.

Here is a hello-world-like code template of `struct.tpl`:

``` go
// Code generated by archivist. DO NOT EDIT.

package {{.Pkg}}
{{""}}
{{- range $i, $v := .Nodes}}
type {{toPascal (lookupStructName .)}} struct {}
{{- end}}
```

For more information, please refer to [Generate Data Structure Definitions for Another Programming Language](#code-templates).

- **Does archivist support multi-config-version?**

No. I abandoned the idea. But you can still make the feature based on archivist. It's not that difficult.

- **What does WATCHER mean?**

Watcher is our game development solution and archivist is a tiny part of watcher. Watcher is not open source.


&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
&nbsp;<br/>
