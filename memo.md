Works in next stage:

- [ ] Support rollback scripts. When error occurs in a single step of a target, the deployment process should be stopped and rollback to the state before this process start.
- [x] Support pull file or directory from nodes. Rethink the design of file functions.
- [ ] Clean ugly code.
- [x] Support customized environment variables from nodes' configuration.
- [ ] Support deploy applications parallelly if their script running orders do not matter.
- [ ] Support network auto-connect scripts. Provided by user, run by Octopoda if needed.
