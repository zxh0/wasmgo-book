(module
  (import "env" "assert_eq_i32" (func $assert_eq_i32 (param i32 i32)))
  (start $main)
  (func $main (export "main")
    (call $assert_eq_i32
      (i32.const 123)
      (select (i32.const 123) (i32.const 456) (i32.const 1))
    )
    (call $assert_eq_i32
      (i32.const 456)
      (select (i32.const 123) (i32.const 456) (i32.const 0))
    )
    (call $assert_eq_i32
      (i32.const 123)
      (drop (i32.const 123) (i32.const 456))
    )
  )
)
