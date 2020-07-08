;; --enable-exceptions
(module
  (import "env" "print_i32" (func $print_i32 (param i32)))
  (import "env" "print_i64" (func $print_i64 (param i64)))

  (import "env" "x0" (event $x0 (param f32)))
  (event $x1 (param i32))
  (event $x2 (param i64 i64))
  (export "x1" (event $x1))

  (func $f1
    (throw $x1 (i32.const 123))
    (throw $x2 (i64.const 123) (i64.const 456))
  )
  (func $f2
    (block (result i64 i64)
      (block (result i32)
        (try
          (do (call $f1))
          (catch
            (br_on_exn 1 $x1) ;; --+
            (br_on_exn 2 $x2) ;; --|--+
            (rethrow)         ;;   |  |
          )                   ;;   |  |
        )                     ;;   |  |
        (i32.const 0)         ;;   |  |
      )                       ;;   |  |
      (call $print_i32) ;; <-------+  |
      (i64.const 1)     ;;            |
      (i64.const 2)     ;;            |
    )                   ;;            |
    (call $print_i64)   ;; <----------+
    (call $print_i64)   ;;
  )
)