# Mock Resource

Used for testing resource (and resource type) actions without depending on the
outside world.

All source fields and parameters are optional.

## Source Configuration

* `mirror_self`: For use on `task.image_resource`, returns itself as an image that tasks can run on.

* `initial_version`: Initial version that the resource should emit. Defaults to `mock`.

* `no_initial_version`: Disable initial version, useful for testing resource triggers.

* `force_version`: Always emit a version regardless of any param during `check`.

* `create_files`: Example:
  ```yaml
  create_files:
    file1.yml: |
      foo: bar
    file2.sh |
      #!/bin/sh
      echo "hello world"
  ```

* `check_delay`: Amount of time to sleep before returning from the `check`. Uses Go duration format.

* `check_failure`: Force every check to return this error message.

* `metadata`: List of name-value pairs to return as metadata on every get and put. Example:
  ```yaml
  metadata:
  - name: foo
    value: bar
  - name: baz
    value: qux
  ```

* `log`: Print a message on every action.

## Behavior

### `check`:

If `force_version` is set, then only that version will ever be emitted, otherwise any version passed into the `check` is returned. This way you can generate a new version `foo` by running `fly check-resource --resource pipeline/resource --from version:foo`.

### `in`:

Since the mock resource is included as a base resource type, it can be used as the image for a task without requiring internet connection (unlike registry-image resource). It can also be used to create files without requiring external state.

* `mirror_self_via_params`: Same as configuring `mirror_self` in source when set to true. Default false.

* `create_files_via_params`: Similar to `create_files` in source; merged in so that additional (or replaced) files can be specified.

### `out`:

* `version`: Version to create.

* `print_env`: Print all environment variables to stdout when set to true. Default false.
