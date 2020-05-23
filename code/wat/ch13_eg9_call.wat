(module
  (func $test (param $a i32) (result i32)
    (if (i32.eqz (local.get $a))
      (then (return (i32.const 0)))
    )
    (call $test (i32.sub (local.get $a) (i32.const 1)))
  )
)