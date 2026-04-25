# portmap-cli

A lightweight CLI tool to manage and visualize active port mappings and process bindings on a local machine.

---

## Installation

```bash
go install github.com/yourusername/portmap-cli@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portmap-cli.git
cd portmap-cli
go build -o portmap-cli .
```

---

## Usage

List all active port mappings:

```bash
portmap-cli list
```

Filter by a specific port:

```bash
portmap-cli list --port 8080
```

Show detailed process bindings:

```bash
portmap-cli show --verbose
```

Kill the process bound to a port:

```bash
portmap-cli kill --port 3000
```

Example output:

```
PORT     PROTOCOL  PID     PROCESS
8080     TCP       12345   node
5432     TCP       6789    postgres
3000     TCP       11223   go-app
```

---

## Requirements

- Go 1.21+
- Linux, macOS, or Windows

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

This project is licensed under the [MIT License](LICENSE).