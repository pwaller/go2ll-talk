all:
	go run . > main.ll
	clang -Wno-override-module main.ll
	./a.out || true # Note: exit status isn't set properly.