(module
  (func $test (param $a i32)
    (block
      (block
        (block
          (br_table 0 1 2 3 (local.get $a))
        )
      )
    )
  )
)