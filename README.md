# Extension UI Demo

A Pencode extension demonstrating actions, dialogs, menus, popups, panels, and UI lifecycle events.

Uses the [Pencode](https://github.com/pencode-labs/pencode) peb v14 SDK. Requires Go 1.27.0 or newer.

## Build and install

```sh
git clone git@github.com:pencode-labs/ui-demo.git
cd ui-demo
go build -o bin/ui-demo .
```

In Pencode's extension manager, choose **Install from local directory** and select this checkout. Extensions run on the workspace server, so build and install on that host for remote workspaces.

On Windows, build `bin/ui-demo.exe` and set `entry` in `resources/extension.json` to `bin/ui-demo.exe` before installing.

## Settings

Open **设置 → 扩展 → 扩展 UI 示例** to change the default name displayed in the result panel. Changes are saved on the workspace host and delivered to the running extension.

## Development

```sh
go test ./...
```

## License

MIT. See [LICENSE](LICENSE).

## Releases

Push a `v*` tag (for example, `v0.1.0`) to this repository, or publish a GitHub Release,
to test and build Linux, macOS and Windows packages for amd64 and arm64. The workflow
uploads ZIP archives and `SHA256SUMS` to the matching release, creating it if needed.
Prerelease tags such as `v0.1.0-rc.1` create prereleases.

Each ZIP contains `bin/`, `resources/`, the README and license. Its manifest uses the
tag version and the correct executable suffix for Windows. Download the package
for the workspace server's OS and architecture, extract it, then install the
extracted directory through Pencode's extension manager.
