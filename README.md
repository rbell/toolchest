# Tool Chest
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
- Additional tools to come!

## Contribution
Contributing to this project is welcome!  There are some guidelines you should consider described in the `CONTRIBUTION.md` file.

## Code Of Conduct
Code of conduct is described in the `CODEOFCONDUCT.md` file.