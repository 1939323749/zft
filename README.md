# ZFT - Command Line File Uploader

ZFT is a powerful command-line tool designed to make it easy to upload files to a CDN (Content Delivery Network). With just a single command, you can quickly and securely transfer your files to the CDN, making them available for fast and reliable access from anywhere in the world.

## Features

- Simple and intuitive command-line interface
- Supports multiple file formats
- Fast and secure file upload
- Automatic file optimization for better performance
- Cross-platform compatibility (Windows, macOS, and Linux)

## Installation

To install ZFT, download the latest release for your platform from the [releases](https://github.com/yourusername/zft/releases) page. Extract the archive and add the executable to your system's `PATH`.

## Usage

Once ZFT is installed, you can start uploading files to your CDN with a single command:

```sh
zft upload --file example.png --key your-cdn-api-key
```

This command will upload the file `example.png` to your CDN using the provided API key. You can also upload multiple files at once:

```sh
zft upload --files file1.png,file2.jpg,file3.svg --key your-cdn-api-key
```

For more advanced usage, you can use the following options:

- `--folder`: Upload all the files in a specific folder
- `--recursive`: Upload files in subfolders as well
- `--ignore`: Ignore certain file types or patterns
- `--optimize`: Optimize files before uploading (e.g., compress images)

## Configuration

You can create a configuration file to store your CDN API key and other settings. By default, ZFT looks for a `.zftconfig` file in your home directory or the current working directory.

Here's an example configuration file:

```json
{
  "apiKey": "your-cdn-api-key",
  "optimize": true,
  "ignore": ["*.tmp", "*.log"]
}
```

With this configuration, you can omit the `--key` option when using the `zft upload` command.

## Contributing

We welcome contributions to the ZFT project! If you'd like to contribute, please follow these steps:

1. Fork the repository
2. Create a new branch for your feature or bugfix
3. Make your changes and commit them to your branch
4. Create a pull request, describing your changes in detail

We'll review your pull request and merge it if it meets our quality standards and aligns with the project's goals.

## License

ZFT is released under the [MIT License](LICENSE).