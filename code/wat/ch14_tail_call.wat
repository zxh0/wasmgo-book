;; --enable-tail-call
(module
  (func $fac (param $x i64) (result i64)
    (return_call $fac-aux (get_local $x) (i64.const 1))
  )
  
  (func $fac-aux (param $x i64) (param $r i64) (result i64)
    (if (i64.eqz (get_local $x))
      (then (return (get_local $r)))
      (else
        (return_call $fac-aux
          (i64.sub (get_local $x) (i64.const 1))
          (i64.mul (get_local $x) (get_local $r))
        )
      )
    )
    (i64.const 0) ;; compiler bug?
  )
)