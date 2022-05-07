<div align="center">

# preppy

A requirements.txt auto-generator<br>

</div>

## About
`preppy` is a `requirements.txt` auto-generator. It is a simple go binary that recursively scans your project directory and parses all the python files to generate a list of dependencies.

`preppy` is smart enough to not include packages from the standard library or local packages under development and pins the appropriate version if the package is already installed. However `preppy` is not perfect and might miss dev dependencies not referenced inside source files.

**Why not just `pip freeze > requirements.txt`?**

Well manually editing out `requirements.txt` generated so by cross-referencing your source is a pain if you have a lot of dependencies or if you have a shared/global environment.


## Getting Started

preppy provides pre-compiled binary releases for Linux. 

Simply download the latest release from github and put it on $PATH.

On other platforms, building from source is recommended.

**Pre-requisites for building from source**
- [`go`](https://golang.org/) For building the preppy binary
- [`python`](https://www.python.org/) For generating stdlib.json
- [`upx`](https://upx.github.io/) (Optional) For compressing the binary
- [`make`](https://www.gnu.org/software/make/) (Optional) For convenience

## Usage

Simply run `preppy` from the root of your project.
```
preppy
```

Optional Flags:
- `-d` : Dry run. Prints the generated requirements.txt to stdout.
- `-r` : Specify root directory. Default is current directory.

## Roadmap

`preppy` follows the linux philosophy of small, simple tools. It is not a full-featured dependency manager. It is a tool to generate a `requirements.txt` file from your project.

However, it aims to have following capabilities in the near future:

- [x] Add support to parse pre-existing `requirements.txt` files and selectively update dependencies.
- [x] Parsing precommit hooks for dev dependencies.
- [ ] Support for loading packages from Conda instead of pip.
- [ ] Add support for specifying virtual environment using `-e` flag.


## Acknowledgements

[Tree-sitter](https://github.com/tree-sitter/tree-sitter) is a fantastic tool for parsing source code even in presence of syntax errors. Thanks to Tree-sitter team for making such a great tool available. Also thanks to Maxim Sukharev ([smacker](https://github.com/smacker)) for providing go bindings for it.

Thanks to [isort](https://github.com/PyCQA/isort) authors for providing a way to list packages from standard library across multiple python versions.
