(module
  (func $test
    (local $i i32)
    (loop
      (local.set $i (i32.add (local.get $i) (i32.const 1)))
      (br_if 0 (i32.lt_s (local.get $i) (i32.const 100)))
    )
  )
)