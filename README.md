## Local Development

Build CLI:
```sh
make build-cli
```

Clean playground output:
```sh
make playground-clean
```

Build playground project:
```sh
make playground-build
```

Serve playground:
```sh
make playground-serve
```
### Variable Naming Conventions

- **`inputDir`** – The raw path to the blog directory as provided by the user via the CLI.  
  It should only be used to resolve the normalized path (`pathToBlog`).

- **`pathToBlog`** – A path to the blog directory, resolved relative to the current working directory (`cwd`).  
  This should be used consistently throughout the codebase for all operations on the blog.