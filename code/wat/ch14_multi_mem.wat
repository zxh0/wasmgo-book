;; --enable-multi-mem
(module
(;;
  (import "env" "mem" (memory $m0 1 8))
  (memory $m1 1 8)
  (memory $m2 1 8)

  (data $m1 (offset (i32.const 1)) "Hello, ")
  (data $m2 (offset (i32.const 2)) "World! ")

  (func
    (memory.size $m1) (drop)
    (memory.grow $m2 (i32.const 2)) (drop)
    (i32.load $m1 offset=12 (i32.const 34))
    (i32.store $m2 offset=56 (i32.const 78))
  )
;;)
)