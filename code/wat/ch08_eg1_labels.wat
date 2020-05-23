(module
  (func
    (block
      (i32.const 100) (br 0) (drop)
    )
    (loop ;; infinite loop!
      (i32.const 200) (br 0) (drop)
    )
    (if (i32.eqz (i32.const 300))
      (then (i32.const 400) (br 0) (drop) )
      (else (i32.const 500) (br 0) (drop) )
    )
  )
)