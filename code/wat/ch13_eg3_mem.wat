(module
  (memory (data "Hello, World!\n"))
  (func $mem
    (i32.const 2) (i32.const 3)
    (i32.load offset=5)
    (i32.store offset=6)
    (memory.size) (drop)
    (i32.const 4) (memory.grow) (drop)
  )
)