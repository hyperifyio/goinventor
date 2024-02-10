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

Copyright (c) Heusala Group Ltd. All rights reserved.

`goinventor` initially operates under the HG Evaluation and Non-Commercial License for two years, post which 
it transitions to the MIT license for broader usage, including commercial purposes. Refer to 
[LICENSE.md](LICENSE.md) for details.

**Commercial licenses can be obtained under separate agreements.**

## Usage and Configuration

For usage and configuration options:

```bash
./goinventor --help
```

Output:

```
Usage of ./goinventor:
  -env-prefix string
        Prefix for ENV variables (default "INVENTORY_")
  -host string
        Get host specific values
  -list
        List all hosts
  -nats string
        The NATS server URL (default "nats://127.0.0.1:4222")
  -source string
        The key value store source (env/nats) (default "env")
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

## Inventory Naming Convention

The `goinventor` tool employs a specific naming convention for environment 
variables to efficiently manage and categorize inventory data for Ansible. This 
convention facilitates the dynamic allocation of variables to the appropriate 
groups and hosts within the inventory.

### Format

The environment variables follow these patterns:

1. **Host-Specific Variables**:  
   Format: `PREFIX_group_hostname_key=value`  
   Description: Sets a variable (`key=value`) for a specific host (`hostname`) 
   within a named group (`group`).  
   Example: `INVENTORY_webservers_web1_ansible_host=192.168.1.10`

2. **Host Variables Without Group**:  
   Format: `PREFIX__hostname_key=value`  
   Equivalent to: `PREFIX_ungrouped_hostname_key=value`  
   Description: Assigns a variable to a host not explicitly assigned to any
   other group.  
   Example: `INVENTORY__web1_ansible_host=192.168.1.10`

3. **Group-Specific Variables**:  
   Format: `PREFIX_group__key=value`  
   Description: Sets a variable for all hosts within a specified group.  
   Example: `INVENTORY_webservers__ansible_user=admin`

4. **Global Variables**:
   Format: `PREFIX___key=value`  
   Equivalent to: `PREFIX_all__key=value`  
   Description: Defines global variables applicable to all groups and hosts.  
   Example: `INVENTORY___ansible_connection=ssh`

### Usage in Ansible

These environment variables are parsed by `goinventor` to construct a dynamic
inventory for Ansible. The tool categorizes and structures the data according to
the established naming conventions, ensuring that each variable is accurately
assigned in the Ansible inventory.

By adhering to these conventions, users can dynamically manage complex 
inventories with ease, providing a flexible and powerful approach to configuring
Ansible playbooks and roles.
