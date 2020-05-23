;; --enable-bulk-memory
(module
  (memory 1 8)
  (data $d1 (offset (i32.const 123)) "Hello, World!") ;; active
  (data $d2 "Goodbye!") ;; passive

  (func $init (param $dst i32) (param $src i32) (param $size i32)
    (memory.init $d2 (local.get $dst) (local.get $src) (local.get $size))
  )
  (func $copy (param $dst i32) (param $src i32) (param $size i32)
    (memory.copy (local.get $dst) (local.get $src) (local.get $size))
  )
  (func $fill (param $dst i32) (param $val i32) (param $size i32)
    (memory.fill (local.get $dst) (local.get $val) (local.get $size))
  )
  (func $drop
    (data.drop $d2)
  )
)