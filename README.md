# hyperifyio/goinventor

`goinventor` is a versatile command line tool that bridges various key-value
stores with Ansible's dynamic inventory system. It's crafted to simplify and 
enhance the management of Ansible inventories, leveraging the power of different
data sources.

![main branch status](https://github.com/hyperifyio/goinventor/actions/workflows/build.yml/badge.svg?branch=main)
![dev branch status](https://github.com/hyperifyio/goinventor/actions/workflows/build.yml/badge.svg?branch=dev)

## Quick Start Guide

To get started with `goinventor`, download the latest release for your operating
system and execute it with the necessary flags:

```bash
wget https://github.com/hyperifyio/goinventor/releases/download/v0.0.1/goinventor-v0.0.1-linux-amd64.zip
unzip ./goinventor-v0.0.1-linux-amd64.zip
cd goinventor-v0.0.1-linux-amd64
./goinventor
```

For global installation:

```
sudo cp ./goinventor /usr/local/bin/goinventor
```

## Setting Up for Development

To build `goinventor` from source:

```bash
git clone git@github.com:hyperifyio/goinventor.git
cd goinventor
make
./goinventor
```

## License

`goinventor` initially operates under the HG Evaluation and Non-Commercial License for two years, post which it transitions to the MIT license for broader usage, including commercial purposes. Refer to [LICENSE.md](LICENSE.md) for details.

Commercial licenses can be obtained under separate agreements.

## Usage and Configuration

For usage and configuration options:

```bash
./goinventor --help
```

## Integrating with Ansible

Test with Ansible and env based inventory:

```
INVENTORY_hosts_host1_hostname=host1 ansible-inventory -i goinventor --list
```

JSON Output:

```json
{
  "_meta": {
    "hostvars": {
      "host1": {
        "hostname": "host1"
      }
    }
  },
  "all": {
    "children": [
      "ungrouped",
      "hosts"
    ]
  },
  "hosts": {
    "hosts": [
      "host1"
    ]
  }
}
```

YAML Output:

```
INVENTORY_hosts_host1_hostname=host1 ansible-inventory -i goinventor --list --yaml
```

```yaml
all:
  children:
    hosts:
      hosts:
        host1:
          hostname: host1
```

### Retrieving Single Host Variables

Command:

```
INVENTORY_hosts_host1_hostname=host1 ansible-inventory -i goinventor --host host1
```

Result:

```json
{
  "hostname": "host1"
}
```
