@$const_str = global [16 x i8] c"Hello world %d\0A\00"

declare void @printf(...)

define void @main() {
0:
	%0 = add i64 30, 12
	call void (...) @printf([16 x i8]* @$const_str, i64 %0)
	ret void
}
