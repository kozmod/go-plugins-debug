TEG ?=test:v0.0.1

clean:
	rm -rf out

plugins:
	go build -o out/plug1.so -gcflags="all=-N -l" -buildmode=plugin ./plugin/one/plugin1.go
	go build -o out/plug2.so -gcflags="all=-N -l" -buildmode=plugin ./plugin/two/plugin2.go

main:
	go build -o "out/main" -gcflags="all=-N -l"

# docker run -it --name test -p 8080:8080 -p 40000:40000 --security-opt='apparmor=unconfident' --cap-add=SYS_PTRACE  test:v0.0.1
# dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec ./out/main -- 2
docker:
	docker image build -t ${TEG}  .
	docker image prune --filter label=stage=builder <<< y