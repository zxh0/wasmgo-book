(module
  (type $ft1 (func))
  (type $ft2 (func))
  (table funcref (elem $f1 $f1 $f1))
  (func $f1
    (call $f1)
    (call_indirect (type $ft2) (i32.const 2))
  )
)