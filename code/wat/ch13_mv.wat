(module
  (import "env" "swap0" (func $swap0 (param i32 i32) (result i32 i32)))
  (func $swap1 (export "swap1") (param i32 i32) (result i32 i32)
    (local.get 0) (local.get 1) (call $swap0)
  )
  (func $add
    (i32.const 1) (i32.const 2)
    (block (param i32 i32) (result i32)
      (i32.add)
    )
    (drop)
  )
)