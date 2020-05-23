(module
  (import "env" "assert_eq_i32" (func $assert_eq_i32 (param i32 i32)))
  (start $main)
  (func $main (export "main")
    (call $assert_eq_i32 (call $test) (i32.const 223))
  )
  (func $test (result i32)
    (i32.const 100) (block (result i32)
      (i32.const 200) (block (result i32)  
        (i32.const 300) (block (result i32) 
          (i32.const 123) (br 2) ;; <---
        ) (i32.add)
      ) (i32.add)
    ) (i32.add)
  )
)

;; n=0: 723
;; n=1: 423
;; n=2: 223
;; n=3: 123
