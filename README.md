# Mock Resource

Used for testing resource (and resource type) actions without depending on the
outside world.

All source fields and parameters are optional.

<a href="https://ci.concourse-ci.org/teams/main/pipelines/resource/jobs/build?vars.type=%22mock%22">
  <img src="https://ci.concourse-ci.org/api/v1/teams/main/pipelines/resource/jobs/build/badge?vars.type=%22mock%22" alt="Build Status">
</a>

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

#### Parameters

* `mirror_self_via_params`: Same as configuring `mirror_self` in source when set to true. Default false.

* `create_files_via_params`: Similar to `create_files` in source; merged in so that additional (or replaced) files can be specified.

#### Files created by the resource

* `./version`: A file containing the version.

* `./privileged`: A file containing a boolean for if the resource was run with elevated privileges.

* `./rootfs/` and `./metadata.json`: Only populated if mirroring self, contains the rootfs and metadata to run this resource as an image.

* Any files specified by the create files source or param.

### `out`:

#### Parameters

* `version`: Version to create.

* `print_env`: Print all environment variables to stdout when set to true. Default false.

* `file`: Contents will be read from the file and emitted as the version
