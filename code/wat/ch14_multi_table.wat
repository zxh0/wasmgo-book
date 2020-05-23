;; --enable-reference-types
(module
  (import "env" "t0" (table $t0 1 8 funcref))
  (table $t1 1 8 funcref)
  (table $t2 1 8 funcref)

  (elem $t1 (offset (i32.const 1)) $f1 $f1 $f1)
  (elem $t2 (offset (i32.const 2)) $f2 $f2 $f2)

  (func $f1)
  (func $f2)
  (func $f3
    (call_indirect $t1 (type 0) (i32.const 1))
  )
)