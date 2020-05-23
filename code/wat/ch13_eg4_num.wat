(module
  (func $num (result f32)
    (i32.const 123) (i32.const 456) (i32.const 789)
    (i32.add) (i32.div_s) (f32.convert_i32_s)
  )
)