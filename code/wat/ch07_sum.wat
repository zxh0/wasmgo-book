(module
  (import "env" "assert_eq_i32" (func $assert_eq_i32 (param i32 i32)))
  (start $main)
  (func $main (export "main")
    (call $assert_eq_i32 (call $sum (i32.const 100)) (i32.const 5050))
  )
  (func $sum (param $a i32) (result i32)
    (local.get $a)

    (br_if 0 (i32.eq (local.get $a) (i32.const 0)))
    (br_if 0 (i32.eq (local.get $a) (i32.const 1)))

    (call $sum (i32.sub (local.get $a) (i32.const 1)))
    (i32.add)
  )
)
