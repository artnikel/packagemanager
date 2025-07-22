# Package Manager (PM)

A simple package manager on Go that allows you to package files into archives, upload them to a remote server via SSH, and install packages.

---

## Installation

### Build
```bash
git clone https://github.com/artnikel/packagemanager.git
cd packagemanager
go mod tidy
go build -o pm cmd/pm/main.go
```

---

## Setup

Create file `config.yaml` in the project dir:

```yaml
ssh:
  host: "your-server.com"
  port: "22"
  user: "username"
  password: "password"       
  key_path: ""               # instead of a password

path:
  remote_package_dir: "/tmp/packages" 
  local_package_dir: "./packages"    
```

---

## Usage

### Packet creation

1. Create a package description file (for example, `packet.json`):

```json
{
  "name": "packet-1",
  "ver": "1.10",
  "targets": [
    "./archive_this1/*.txt",
    {"path": "./archive_this2/*", "exclude": "*.tmp"}
  ],
  "packets": [
    {"name": "packet-3", "ver": "<=2.0"}
  ]
}
```

2. Execute the create command:

```bash
./pm create ./packet.json
```

This will create an archive `packet-1--1.0.0.tar.gz` and upload it to the server.

### Installing packages

1. Create a file with a list of packages (for example, `packages.json`):

```json
{
  "packages": [
    {"name": "packet-1", "ver": ">=1.10"},
    {"name": "packet-2"},
    {"name": "packet-3", "ver": "<=1.10"}
  ]
}
```

2. Execute the install command:

```bash
./pm update ./packages.json
```

The packages will be downloaded from the server and unpacked into the `./packages/`.

---

## File format

### packet.json (package description)

**Fields:**
- `name` - package name
- `ver` - package version
- `targets` - file path array for archiving

### packages.json (installation list)

**Fields:**
- `name` - package name for installation
- `ver` - version condition (optional)
  - `>=1.10` - version greater than or equal to 1.10
  - `<=1.10` - version less than or equal to 1.10
  - `1.0.0` - exact version


---

## Commands

```bash
# Show help
./pm

# Create a packet
./pm create ./packet.json

# Install packages
./pm update ./packages.json
```

