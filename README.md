# Extension UI Demo

A Pencode extension demonstrating actions, dialogs, menus, popups, panels, and UI lifecycle events.

Uses the [Pencode](https://github.com/pencode-labs/pencode) peb v6 SDK. Requires Go 1.27.0 or newer.

## Build and install

```sh
git clone https://github.com/pencode-labs/ui-demo.git
cd ui-demo
go build -o bin/ui-demo .
```

In Pencode's extension manager, choose **Install from local directory** and select this checkout. Extensions run on the workspace server, so build and install on that host for remote workspaces.

On Windows, build `bin/ui-demo.exe` and set `entry` in `resources/extension.json` to `bin/ui-demo.exe` before installing.

## Development

```sh
go test ./...
```

## License

MIT. See [LICENSE](LICENSE).
