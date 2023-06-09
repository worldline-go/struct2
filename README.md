# struct2

[![License](https://img.shields.io/github/license/worldline-go/struct2?color=red&style=flat-square)](https://raw.githubusercontent.com/worldline-go/struct2/main/LICENSE)
[![Coverage](https://img.shields.io/sonar/coverage/worldline-go_struct2?logo=sonarcloud&server=https%3A%2F%2Fsonarcloud.io&style=flat-square)](https://sonarcloud.io/summary/overall?id=worldline-go_struct2)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/worldline-go/struct2/test.yml?branch=main&logo=github&style=flat-square&label=ci)](https://github.com/worldline-go/struct2/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/worldline-go/struct2?style=flat-square)](https://goreportcard.com/report/github.com/worldline-go/struct2)
[![Go PKG](https://raw.githubusercontent.com/worldline-go/guide/main/badge/custom/reference.svg)](https://pkg.go.dev/github.com/worldline-go/struct2)

This repository helps to work with struct, convert map and get information about that.

This is a modified version of [common struct to map libraries](#inspired-projects) with cool features.

## Usage

```sh
go get github.com/worldline-go/struct2
```

### Map

Get decoder and run `Map` method, default looking the `struct` tag in struct.

Supported tags: `-`, `omitempty`, `string`, `ptr2`, `omitnested`, `flatten`.

Convertion order is __`-`, `omitempty`, `string`, `ptr2`, custom hook function, hooker interface, `omitnested` + `flatten`__

```go
type ColorGroup struct {
    ID     int      `db:"id"`
    Name   string   `db:"name"`
    Colors []string `db:"colors"`
    // custom type with implemented Hooker interface
    // covertion result to time.Time
    Date types.Time `db:"time"`
	// RGB unknown type but to untouch it add omitnested to keep that struct type
	RGB *rgb `db:"rgb,omitempty,omitnested"`
}

//...

// default tagName is `struct`
decoder := struct2.Decoder{
    TagName: "db",
}

// get map[string]interface{}
result := decoder.Map(group)

// or use one line
// result := (&struct2.Decoder{}).SetTagName("db").Map(group) // default tag name is "struct"
```

Custom decoder can be use in struct which have `struct2.Hooker` interface.  
Or set a slice of custom `struct2.HookFunc` functions in decoder.

Check documentation examples.

#### Tags Information

__omitnested__: very helpful to don't want to touch data.

__ptr2__: convert pointer to the concrete value, if pointer is nil new value generating.
ptr2 to effect custom hook functions and hooker interface also omitnested.

### Decode

Decode is working almostly same as the `mitchellh/mapstructure` repo.

Default tag is `struct` in struct.

```go
decoder := struct2.Decoder{
    WeaklyTypedInput: true,
    WeaklyIgnoreSeperator: true,
    TagName: "struct",
    BackupTagName: "json",
}

//...

// input and output could be any, output should be pointer
if err := d.Decode(input, output); err != nil {
    return err
}
```

---

## Inspired Projects

When starting this project, I want to make it from scratch and I can learn more and make some new features but codes turning to `fatih/structs` repo and they are solved lots of problems so I copied parts in there and add some features as hook functions. After that I want to extend that one to make map to struct operations. In that time, I see I need to check all types due to mixing types together. So I copied parts in `mitchellh/mapstructure`. Thanks for all people to make them.

[fatih/structs](https://github.com/fatih/structs)  
[mitchellh/mapstructure](https://github.com/mitchellh/mapstructure)
