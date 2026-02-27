# Contributing to ModelHawk

## Development Setup

### Building

Generate protocol buffer code:

```sh
make generate
```

This generates:
- Go code in `gen/go/`
- TypeScript code in `gen/ts/`
- Documentation in `gen/docs/`

## Making Changes

### Protocol Buffer Changes

1. Edit `.proto` files in `proto/v0/`
2. Run `make -j generate` to regenerate code
3. Commit both the `.proto` files and generated code
4. Open a pull request

The CI will verify that generated files are up-to-date.

## Release Process

Releases are automated via GitHub Actions. Javascript/Typescript files are published to npm as `@dshearer/modelhawk`.

### Steps to Release

1. **Update the version** in [package-config/ts/package.json](package-config/ts/package.json):

2. **Commit and merge via pull request**:

   ```sh
   git add package-config/ts/package.json
   make -j generate
   git commit -m "chore: bump version to 0.2.0"
   git push origin your-branch
   # Create PR, get approval, merge to main
   ```

3. **Push a version tag**:

   ```sh
   git checkout main
   git pull
   git tag v0.2.0
   git push origin v0.2.0
   ```

4. **Automated workflow runs**:
   - Verifies the tag matches the version in `package-config/ts/package.json`
   - Generates fresh protobuf code
   - Builds and publishes `@dshearer/modelhawk` to npm with provenance
   - Creates a GitHub release with auto-generated notes
