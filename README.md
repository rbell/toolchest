# Tool Chest

![coverage](https://raw.githubusercontent.com/rbell/toolchest/badges/.badges/master/coverage.svg)
![lint](https://github.com/rbell/toolchest/actions/workflows/lint.yml/badge.svg?branch=master)
![Go Report Card](https://goreportcard.com/badge/github.com/rbell/toolchest?cache=v1)
![Release](https://img.shields.io/github/release/rbell/toolchest.svg?style=flat-square)

Provides a suite of golang tools that can be referenced and used in golang projects.

## License
This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.

## Tools
- Work Queue, located in `toolchest/workqueue`, provides a que which work can be submitted to and worked on by multiple go routines.  Features include:
  - Ability to configure queue length and number of go routines operating on queue
  - Error monitoring
  - Prioritizing work submitted to queue
  - Option to name work submitted to queue for later reference
  - Stopping and breaking the queue
  - Resize queue length
  - Dequeue work in queue
  - View work in queue and it's current state, priority and position in queue.
- ValidationError
  - Facilitates error reflecting validation issues to a user. 
  - Supports warnings and errors
  - Reflects validation errors over nested types (i.e. Customer has an address)
- Storage
  - SafeMap
    - A thread safe map
  - FifoMapCache
    - A thread safe map with a maximum size.  When the cache is full, the oldest entries are evicted.
  - GenericStack
    - A generic stack data structure
- Propositions
  - Provides a set of proposition functions in Go. These functions allow evaluations of various conditions on various types, each function returning either true or false.
- SliceOps
  - Provides a set of utility functions for working with slices in Go. These functions include operations such as cutting, removing, inserting, filtering, pushing, popping, and more.

- Additional tools to come!

## Contribution
Contributing to this project is welcome!  There are some guidelines you should consider described in the `CONTRIBUTION.md` file.

## Code Of Conduct
Code of conduct is described in the `CODEOFCONDUCT.md` file.