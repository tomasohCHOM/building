# Container

Building a [container using less than 100 lines of Go](https://www.infoq.com/articles/build-a-container-golang/), covering the basics of containerization, creating containers, and building a simple one from scratch using [Go](https://go.dev/).

## What are containers

Often defined as software units that package code and dependencies so that they can run seamlessly across environments. Three major components (among other things) are used to create these containers: namespaces, cgroups, and layered filesystems.

### Namespaces

They provide the isolation needed to run multiple containers on one machine, all with seemingly their own environments. Types of namespaces:

- **PID:**  assigns a set of process IDs independent from PIDs in other namespaces.
- **MNT:** processes in the namespace get their own mount table, allowing them to mount/unmount file systems without affecting the host filesystem.
- **NET:** gives their processes their own network stack (private routing table, set of IP addresses, socket listing, connection tracking table, firewall, etc.).
- **UTS:** allows processes to have different host and domain names from others.
- **IPC:** own IPC resources, like POSIX message groups.
- **USER:** own set of user and group IDs, giving processes `root` privilege within its user namespace and not in others.

### Cgroups

A cgroup is a Linux kernel feature that collects a set of process or task IDs and applies limits on resource usage (CPU, memory, disk I/O, network, etc.) to them, enforcing some level of resource sharing between processes.

### Layered Filesystems

Used to move whole machine images around efficiently. They provide optimizations over creating copies of the root filesystem in each container.

## Quick Notes

- The code lets a process run in an isolated namespace with its corresponding filesystem.
- The article skips setting up cgroups and root filesystem management (efficient downloading and caching).
- Skips container setup as well.
