# gbs Build Environment

This folder builds the Docker image that gbs-launched containers use to build code.

# Usage

Builds are run by the simple `build.sh` shell script, and are parameterized only by the following environment variables:

- `SITE`/`ORG`/`REPO` - the build will run `git clone https://$SITE/$ORG/$REPO.git`. It only supports `git`
- `OUT_DIR` - this is where the build will write the final binary (which will be called `$REPO`). In most cases you'll want to write this binary to a mounted volume. See the example below.

Here's a full example of how to use the build environment to build the [Minio client](https://github.com/minio/mc):

```console
docker run --rm -v $PWD:/pwd -e SITE=github.com -e ORG=minio -e REPO=mc -e OUT_DIR=/pwd quay.io/arschles/gbs-env:0.0.1
```

# Building

Use the `Makefile` in this directory to build and push the Docker image for the env. `make docker-build` builds the environment, and `make docker-push` pushes it to [quay.io](https://quay.io/arschles/gbs-env)
