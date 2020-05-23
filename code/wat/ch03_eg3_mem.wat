(module
  (memory 1 8)
  (data (offset (i32.const 100)) "hello")

  (func
    (i32.const 1) (i32.const 2)
    (i32.load offset=100)
    (i32.store offset=100)
    (memory.size) (drop)
    (i32.const 4) (memory.grow) (drop)
  )
)