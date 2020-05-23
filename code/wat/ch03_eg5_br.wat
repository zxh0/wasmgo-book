(module
  (func
    (block
      (block
        (block
          (br 1)
          (br_if 2 (i32.const 100))
          (br_table 0 1 2 3)
          (return)
        )
      )
    )
  )
)