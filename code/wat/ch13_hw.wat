(module
  (import "env" "print_char" (func $print_char (param i32)))
  (memory (data "Hello, World!\n"))
  (global $hw_addr i32 (i32.const 0))
  (global $hw_len i32 (i32.const 14))
  (func $main (export "main")
    (call $print_str (global.get $hw_addr) (global.get $hw_len))
  )
  (func $print_str (export "print_str")
    (param $addr i32) (param $len i32)
    (local $i i32)
    (loop
      (i32.add (local.get $addr) (local.get $i))            ;; tmp = $addr + $i
      (call $print_char (i32.load8_u))                      ;; $print_char(mem[tmp])
      (local.set $i (i32.add (local.get $i) (i32.const 1))) ;; $i++
      (br_if 0 (i32.lt_u (local.get $i) (local.get $len)))  ;; continue if $i < $len
    )
  )
)
