## Linux Distro Engine (Go edition)
---
The Linux Distro Engine (Go ed.) is a "Distro Engine", that attempts to identify the specific Linux Distro being used, and return in in a minimal, string format.

This Go module is written without any external, third-party dependencies, helping mimimize potential attack surface and dependency issues outside of the Go lang itself.

This is primarily a learning project, but it can also be used as a dependency to enhance development of Linux applications.
Even though different Linux distributions are usually not critically different from each other, certain distributions have different software installed by default,
and being able to conveniently identify the distro being run on can be part of a strategy to manage such defaults.

---
[![Go](https://github.com/mportizlunyov/linuxdistroengine-go/actions/workflows/go.yml/badge.svg)](https://github.com/mportizlunyov/linuxdistroengine-go/actions/workflows/go.yml)
