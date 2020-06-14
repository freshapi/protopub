gazelle:
	bazel run :gazelle -- update
	bazel run :gazelle -- update-repos -prune -from_file=go.mod -to_macro=repositories.bzl%go_repositories

build:
	bazel build //cmd/protopub:protopub

cross:
	bazel build --config=linux //cmd/protopub:protopub
	bazel build --config=windows //cmd/protopub:protopub
	bazel build --config=darwin //cmd/protopub:protopub

release: cross
	rm -rf .release
	rm -rf .bin
	mkdir .release
	mkdir .bin
	cp bazel-bin/cmd/protopub/linux_amd64_pure_stripped/protopub .release/protopub_linux_amd64
	cp bazel-bin/cmd/protopub/darwin_amd64_pure_stripped/protopub .release/protopub_darwin_amd64
	cp bazel-bin/cmd/protopub/windows_amd64_pure_stripped/protopub.exe .release/protopub_windows_amd64.exe

.PHONY: gazelle build cross release
