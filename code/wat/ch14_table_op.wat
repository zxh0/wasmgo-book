;; --enable-reference-types
(module
  (table $t1 funcref (elem $f))

  (func $f
    (table.size $t1) (drop)
    (table.grow $t1 (ref.func $f) (i32.const 2)) (drop)
    (table.set $t1 (i32.const 4)
      (table.get $t1 (i32.const 3))
    )
  )
)