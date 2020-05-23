(module
  (func $param (result i32)
    (i32.const 100) (i32.const 200) (drop) (drop)
    (select (i32.const 123) (i32.const 456) (i32.const 1))
  )
)