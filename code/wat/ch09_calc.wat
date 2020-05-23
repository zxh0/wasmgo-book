(module
  (import "env" "assert_eq_i64" (func $assert_eq_i64 (param i64 i64)))
  (type $ft0 (func (param i64 i64) (result i64)))
  (table funcref (elem $add $sub $mul))

  (start $main)
  (func $main (export "main")
    (call $assert_eq_i64 (i64.const 5) (call $calc (i64.const 3) (i64.const 2) (i32.const 0)))
    (call $assert_eq_i64 (i64.const 1) (call $calc (i64.const 3) (i64.const 2) (i32.const 1)))
    (call $assert_eq_i64 (i64.const 6) (call $calc (i64.const 3) (i64.const 2) (i32.const 2)))
  )

  (func $calc (param $a i64) (param $b i64) (param $op i32) (result i64)
    (local.get $a)
    (local.get $b)
    (local.get $op)
    (call_indirect (type $ft0))
  )

  (func $add (type $ft0) (i64.add (local.get 0) (local.get 1)))
  (func $sub (type $ft0) (i64.sub (local.get 0) (local.get 1)))
  (func $mul (type $ft0) (i64.mul (local.get 0) (local.get 1)))
)
