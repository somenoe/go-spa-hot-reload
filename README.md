
# Proof of concept: go-app (SPA) client hot reload

This project is a proof of concept demonstrating a Single Page Application (SPA) built with Go, showcasing client hot reloading capabilities.

## How it works

1. The server exposes an endpoint `/updated` that returns the current timestamp.
2. The client-side periodically fetches this timestamp.
3. If the timestamp changes, it triggers a page reload.

## Tools Used

- [air](https://github.com/air-verse/air) is used for rebuilding the server upon saving a source file.
- [go-app](https://github.com/maxence-charriere/go-app) is used for creating SPAs with pure Go.
- [mage](https://github.com/magefile/mage) is used for build commands.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for more details.
